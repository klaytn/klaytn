package main

import (
	"fmt"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	var allCnt, nilCnt int64
	var db *leveldb.DB
	var err error

	db, err = leveldb.OpenFile("/var/kend/data/klay/chaindata/snapshot", nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		if len(key) == 0 || len(value) == 0 {
			nilCnt++
		} else {
			if allCnt%100000 == 0 {
				fmt.Printf("allCnt : %d\tnilCnt : %d\n", allCnt, nilCnt)
			}
			allCnt++
		}
	}
	iter.Release()
	err = iter.Error()
	if err != nil {
		panic(err)
	}

	fmt.Printf("finally allCnt : %d\tnilCnt : %d\n", allCnt, nilCnt)
	os.Exit(0)
}
