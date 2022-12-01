package core

import (
	"auto-swap/config"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger/v2"
	"github.com/dgraph-io/badger/v2/options"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
	"sync"
	"time"
)

type LvDb struct {
	exit chan struct{}
	sync.WaitGroup
	*badger.DB
	*time.Ticker
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}

func (lvdb *LvDb) Init(name string) error {

	lvdb.exit = make(chan struct{})
	lvdb.Ticker = time.NewTicker(5 * time.Minute)

	dir := path.Join(config.GetDataDir(), name)
	if !isExist(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	opts := badger.DefaultOptions(dir).
		WithValueLogLoadingMode(options.FileIO).
		WithLevelOneSize(64 << 20). //64m
		WithMaxTableSize(4 << 20).  //4m
		WithNumMemtables(2).
		WithCompactL0OnClose(false).
		WithLoadBloomsOnOpen(false).
		WithNumVersionsToKeep(1)
	if runtime.GOOS != "windows" {
		opts = opts.WithTruncate(true)
	}
	db, err := badger.Open(opts)
	if err != nil {
		return err
	}
	lvdb.DB = db
	lvdb.startGcLog()
	return nil
}

func (lvdb *LvDb) startGcLog() {
	go func() {
		lvdb.Add(1)
		defer lvdb.Done()
		defer lvdb.Stop()
		for {
			select {
			case <-lvdb.exit:
				return
			case <-lvdb.C:
				lvdb.gcLog()
			}
		}
	}()
}

func (lvdb *LvDb) gcLog() {
	AGAIN: err := lvdb.RunValueLogGC(0.5)
	if err == nil {
		goto AGAIN
	}
}

func (lvdb *LvDb) Close() error {
	close(lvdb.exit)
	lvdb.Wait()
	lvdb.Ticker.Stop()
	return lvdb.DB.Close()
}

const NodePre ="node_"

func (lvdb *LvDb) CacheNode(node *Node) {
	_ = lvdb.Update(func(txn *badger.Txn) error {
		b,err := json.Marshal(node)
		if err != nil {
			logrus.Errorf("lvdb CacheNode, json.Marshal err: %v", err)
		}
		err =  txn.Set([]byte(NodePre+node.Path),b)
		if err != nil {
			logrus.Errorf("lvdb CacheNode, txn.Set err: %v", err)
		}
		return nil
	})
}

func (lvdb *LvDb) IterNode() chan interface{} {
	chNode :=  make(chan interface{},100)
	go func() {
		defer close(chNode)
		_ = lvdb.View(func(txn *badger.Txn) error {
			opt := badger.DefaultIteratorOptions
			opt.Prefix = []byte(NodePre)
			iter := txn.NewIterator(opt)
			defer iter.Close()
			for iter.Rewind(); iter.Valid(); iter.Next() {
				err := iter.Item().Value(func(val []byte) error {
					node := &Node{}
					err := json.Unmarshal(val, node)
					if err != nil {
						return err
					}else{
						chNode <- node
					}
					return nil
				})
				if err != nil {
					chNode <- err
					break
				}
			}
			return nil
		})
	}()
	return chNode
}

const LkSlippage = "slippage"
func (lvdb *LvDb) CacheSlippage() error {
	return lvdb.Update(func(txn *badger.Txn) error {
		b := make([]byte,8,8)
		binary.BigEndian.PutUint64(b,uint64(config.Cfg.Node.Slippage))
		return txn.Set([]byte(LkSlippage),b)
	})
}
func (lvdb *LvDb) GetSlippage() (slippage int64, err error) {
	return slippage, lvdb.View(func(txn *badger.Txn) error {
		b := make([]byte,8,8)
		binary.BigEndian.PutUint64(b,uint64(config.Cfg.Node.Slippage))
		item,err :=  txn.Get([]byte(LkSlippage))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			slippage = int64(binary.BigEndian.Uint64(val))
			return nil
		})
	})
}

const LkMaxPath2 = "MaxPath2"
func (lvdb *LvDb) CacheMaxPath2(max int) error {
	return lvdb.Update(func(txn *badger.Txn) error {
		b := make([]byte,8,8)
		binary.BigEndian.PutUint64(b,uint64(max))
		return txn.Set([]byte(LkMaxPath2),b)
	})
}
func (lvdb *LvDb) GetMaxPath2() (maxPath2 int, err error) {
	return maxPath2, lvdb.View(func(txn *badger.Txn) error {
		item,err :=  txn.Get([]byte(LkMaxPath2))
		if err == badger.ErrKeyNotFound {
			return nil
		}
		if err != nil  {
			return err
		}
		return item.Value(func(val []byte) error {
			maxPath2 = int(binary.BigEndian.Uint64(val))
			return nil
		})
	})
}

const ReqId = "ReqId"
func (lvdb *LvDb) CacheReqId(reqId uint64) error {
	return lvdb.Update(func(txn *badger.Txn) error {
		b := make([]byte,8,8)
		binary.BigEndian.PutUint64(b,reqId)
		return txn.Set([]byte(ReqId),b)
	})
}
func (lvdb *LvDb) GetReqId() (reqId uint64, err error) {
	return reqId, lvdb.View(func(txn *badger.Txn) error {
		item,err :=  txn.Get([]byte(ReqId))
		if err == badger.ErrKeyNotFound {
			return nil
		}
		if err != nil  {
			return err
		}
		return item.Value(func(val []byte) error {
			reqId = binary.BigEndian.Uint64(val)
			return nil
		})
	})
}