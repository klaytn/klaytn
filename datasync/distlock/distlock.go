// Copyright 2020 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

package distlock

import (
	"errors"
	"github.com/go-redsync/redsync"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/mna/redisc"
	"time"
)

var lockManager *distLockManager
var distLockNotFoundErr = errors.New("there is no distributed lock for requested lock name")

// distLockManager holds redsync.Redsync object and lockMap
type distLockManager struct {
	rs      *redsync.Redsync
	lockMap map[string]DistLock
}

// DistLock is an interface to wrap redsync.Mutex
type DistLock interface {
	Extend() bool
	Lock() error
	Unlock() bool
}

var lockNames = []string{
	"dbsyncer-worker",
	"dbsyncer-checker",
}

func createPool(addr string, opts ...redigo.DialOption) (*redigo.Pool, error) {
	return &redigo.Pool{
		MaxIdle:     5,
		MaxActive:   10,
		IdleTimeout: time.Minute,
		Dial: func() (redigo.Conn, error) {
			return redigo.Dial("tcp", addr, opts...)
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}, nil
}

func initRedSync(redisEndpoints []string) {
	redigoCluster := &redisc.Cluster{
		StartupNodes: redisEndpoints,
		DialOptions:  []redigo.DialOption{redigo.DialConnectTimeout(5 * time.Second)},
		CreatePool:   createPool,
	}

	lockManager = &distLockManager{
		rs:      redsync.New([]redsync.Pool{redigoCluster}),
		lockMap: make(map[string]DistLock),
	}

	for _, name := range lockNames {
		lockManager.lockMap[name] = lockManager.rs.NewMutex(name,
			redsync.SetTries(1),
			redsync.SetExpiry(30*time.Second))
	}
}

func GetLock(name string) (DistLock, error) {
	lock, exist := lockManager.lockMap[name]
	if !exist {
		return nil, distLockNotFoundErr
	}
	return lock, nil
}
