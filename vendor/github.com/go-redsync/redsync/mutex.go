package redsync

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/gomodule/redigo/redis"
)

// A DelayFunc is used to decide the amount of time to wait between retries.
type DelayFunc func(tries int) time.Duration

// A Mutex is a distributed mutual exclusion lock.
type Mutex struct {
	name   string
	expiry time.Duration

	tries     int
	delayFunc DelayFunc

	factor float64

	quorum int

	genValueFunc func() (string, error)
	value        string
	until        time.Time

	pools []Pool
}

// Lock locks m. In case it returns an error on failure, you may retry to acquire the lock by calling this method again.
func (m *Mutex) Lock() error {
	value, err := m.genValueFunc()
	if err != nil {
		return err
	}

	for i := 0; i < m.tries; i++ {
		if i != 0 {
			time.Sleep(m.delayFunc(i))
		}

		start := time.Now()

		n := m.actOnPoolsAsync(func(pool Pool) bool {
			return m.acquire(pool, value)
		})

		now := time.Now()
		until := now.Add(m.expiry - now.Sub(start) - time.Duration(int64(float64(m.expiry)*m.factor)))
		if n >= m.quorum && now.Before(until) {
			m.value = value
			m.until = until
			return nil
		}
		m.actOnPoolsAsync(func(pool Pool) bool {
			return m.release(pool, value)
		})
	}

	return ErrFailed
}

// Unlock unlocks m and returns the status of unlock.
func (m *Mutex) Unlock() bool {
	n := m.actOnPoolsAsync(func(pool Pool) bool {
		return m.release(pool, m.value)
	})
	return n >= m.quorum
}

// Extend resets the mutex's expiry and returns the status of expiry extension.
func (m *Mutex) Extend() bool {
	n := m.actOnPoolsAsync(func(pool Pool) bool {
		return m.touch(pool, m.value, int(m.expiry/time.Millisecond))
	})
	return n >= m.quorum
}

func genValue() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func (m *Mutex) acquire(pool Pool, value string) bool {
	conn := pool.Get()
	defer conn.Close()
	reply, err := redis.String(conn.Do("SET", m.name, value, "NX", "PX", int(m.expiry/time.Millisecond)))
	return err == nil && reply == "OK"
}

var deleteScript = redis.NewScript(1, `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end
`)

func (m *Mutex) release(pool Pool, value string) bool {
	conn := pool.Get()
	defer conn.Close()
	status, err := redis.Int64(deleteScript.Do(conn, m.name, value))

	return err == nil && status != 0
}

var touchScript = redis.NewScript(1, `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("pexpire", KEYS[1], ARGV[2])
	else
		return 0
	end
`)

func (m *Mutex) touch(pool Pool, value string, expiry int) bool {
	conn := pool.Get()
	defer conn.Close()
	status, err := redis.Int64(touchScript.Do(conn, m.name, value, expiry))

	return err == nil && status != 0
}

func (m *Mutex) actOnPoolsAsync(actFn func(Pool) bool) int {
	ch := make(chan bool)
	for _, pool := range m.pools {
		go func(pool Pool) {
			ch <- actFn(pool)
		}(pool)
	}
	n := 0
	for range m.pools {
		if <-ch {
			n++
		}
	}
	return n
}
