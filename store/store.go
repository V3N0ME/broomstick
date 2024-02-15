package store

import (
	"broomstick/controller"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

type configStore struct {
	ds    map[string]int
	wal   *os.File
	mutex sync.RWMutex
}

func NewConfigStore() *configStore {

	wal, err := os.OpenFile("config_store.json", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	all, err := ioutil.ReadFile("config_store.json")
	if err != nil {
		panic(err)
	}

	ds := make(map[string]int)

	if len(all) > 0 {

		if err := json.Unmarshal(all, &ds); err != nil {
			panic(err)
		}
	}

	return &configStore{
		wal: wal,
		ds:  ds,
	}
}

func generateKey(info controller.NodePoolInfo) string {
	return fmt.Sprintf("%s-%s-%s", info.ProjectID, info.Cluster, info.NodePool)
}

func (cs *configStore) StoreNodePoolInfo(info controller.NodePoolInfo, size int) error {

	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	cs.ds[generateKey(info)] = size

	bytes, err := json.Marshal(cs.ds)
	if err != nil {
		return err
	}

	err = cs.wal.Truncate(0)
	if err != nil {
		return err
	}
	_, err = cs.wal.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = cs.wal.Write(bytes)
	if err != nil {
		return err
	}

	cs.wal.Sync()

	return nil
}

func (cs *configStore) GetNodePoolInfo(pool controller.NodePoolInfo) (int, error) {

	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	size, ok := cs.ds[generateKey(pool)]
	if !ok {
		return 0, errors.New("pool not found in store")
	}

	return size, nil
}
