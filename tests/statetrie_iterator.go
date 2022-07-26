package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

var lock sync.RWMutex

func main() {
	var allCnt, nilCnt [4]int64
	wg := new(sync.WaitGroup)

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(i int) {
			var db *leveldb.DB
			var err error
			fmt.Printf("%d start\n", i)
			db, err = leveldb.OpenFile(fmt.Sprintf("/var/kend/data/klay/chaindata/statetrie/%d", i), nil)
			if err != nil {
				panic(err)
			}
			defer db.Close()

			iter := db.NewIterator(nil, nil)
			for iter.Next() {
				key := iter.Key()
				value := iter.Value()
				if len(key) == 0 || len(value) == 0 {
					lock.Lock()
					nilCnt[i]++
					lock.Unlock()
				}
				lock.Lock()
				allCnt[i]++
				lock.Unlock()
			}
			iter.Release()
			wg.Done()
			err = iter.Error()
			fmt.Printf("%d finished, err=%v\n", i, err)
		}(i)
	}
	go func() {
		for {
			time.Sleep(time.Second * 5)
			fmt.Printf("allCnt : %d\tnilCnt : %d\n", allCnt[0]+allCnt[1]+allCnt[2]+allCnt[3], nilCnt[0]+nilCnt[1]+nilCnt[2]+nilCnt[3])
		}
	}()
	wg.Wait()
	fmt.Printf("finally allCnt : %d\tnilCnt : %d\n", allCnt[0]+allCnt[1]+allCnt[2]+allCnt[3], nilCnt[0]+nilCnt[1]+nilCnt[2]+nilCnt[3])
	os.Exit(0)
}
