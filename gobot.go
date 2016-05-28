package gobot

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"net/http"
	_ "net/http/pprof"
	"sync"
)

var (
	workerMu  sync.RWMutex
	adapterMu sync.RWMutex
	workers   = make(map[string]Worker)
	adapter   Adapter
)

func RegisterWorker(name string, worker Worker) {
	workerMu.Lock()
	defer workerMu.Unlock()

	if worker == nil {
		panic("gotbot: Worker cannot be nil.")
	}
	if _, exist := workers[name]; exist {
		panic("gobot: Worker exist : " + name)
	}
	log.Debugf("Add Worker %s", name)
	workers[name] = worker
}

func RegisterAdapter(newadapter Adapter) {
	adapterMu.Lock()
	defer adapterMu.Unlock()
	adapter = newadapter
}

func StartGoBot() error {
	err := initAdapter()
	if err != nil {
		log.Error(err)
		return err
	}

	err = initWorkers()
	if err != nil {
		log.Error(err)
		return err
	}
	return http.ListenAndServe("localhost:6060", nil)
}

func initAdapter() error {
	if adapter == nil {
		return errors.New("gobot's adapter can not be nil.")
	}
	err := adapter.Init()
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func initWorkers() error {
	for name, worker := range workers {
		err := worker.Init()
		return fmt.Errorf("Init worker %s Fail. %s", name, err.Error())
	}
	return nil
}
