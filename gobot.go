package gobot

import (
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
	adapters  = make(map[string]Adapter)
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

type Gobot struct {
	Name     string
	workers  map[string]Worker
	adapters map[string]Adapter
}

func NewDefaultGobot(botname string) *Gobot {
	ret := &Gobot{}
	ret.Name = botname
	ret.workers = workers
	ret.adapters = adapters
	return ret
}

func RegisterAdapter(name string, newadapter Adapter) {
	adapterMu.Lock()
	defer adapterMu.Unlock()
	if newadapter == nil {
		panic("gobot: Adapter cannot be nil.")
	}
	if _, exist := adapters[name]; exist {
		panic("gobot: " + name + " exist.")
	}
	log.Debugf("Add adapter %s", name)
	adapters[name] = newadapter
}

func (bot *Gobot) StartGoBot() error {
	err := bot.initAdapter()
	if err != nil {
		log.Error(err)
		return err
	}

	err = bot.initWorkers()
	if err != nil {
		log.Error(err)
		return err
	}
	go bot.startAdaperts()
	return http.ListenAndServe("localhost:6060", nil)
}

func (bot *Gobot) startAdaperts() {
	for name, adapter := range bot.adapters {
		log.Infof("Start Adapter %s", name)
		go adapter.Start()
	}
}

func (bot *Gobot) initAdapter() error {
	for name, adapter := range bot.adapters {
		err := adapter.Init()
		if err != nil {
			return fmt.Errorf("Init Adapter %s Fail. %s", name, err.Error())
		}
	}
	return nil
}

func (bot *Gobot) initWorkers() error {
	for name, worker := range bot.workers {
		err := worker.Init()
		return fmt.Errorf("Init worker %s Fail. %s", name, err.Error())
	}
	return nil
}