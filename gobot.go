package gobot

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/eternnoir/gobot/payload"
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
	Name       string
	workers    map[string]Worker
	adapters   map[string]Adapter
	ConfigPath string
}

func NewDefaultGobot(botname string) *Gobot {
	ret := &Gobot{}
	ret.Name = botname
	ret.workers = workers
	ret.adapters = adapters
	ret.ConfigPath = "./"
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

func (bot *Gobot) Receive(message *payload.Message) {
	log.Infof("Receive new message. %#v", message)
	if message.SourceAdapter == "" {
		panic("Message's SourceAdapter Id must be seted.")
	}
	for name, worker := range bot.workers {
		// Call workers process
		log.Debugf("Call worker %s process message %#v", name, message)
		err := worker.Process(bot, message)
		if err != nil {
			log.Error(err)
		}
	}
}

func (bot *Gobot) Send(text string) {
	for an, adapter := range bot.adapters {
		log.Debugf("Use adapter %s, Send message %s", an, text)
		go adapter.Send(text)
	}
}

func (bot *Gobot) SendToChat(text, chatroom string) {
	for an, adapter := range bot.adapters {
		log.Debugf("Use adapter %s, Send message %s to ChatRoom %s", an, text, chatroom)
		go adapter.SendToChat(text, chatroom)
	}
}

func (bot *Gobot) Reply(orimessage *payload.Message, text string) error {
	adapter := bot.adapters[orimessage.SourceAdapter]
	return adapter.Reply(orimessage, text)
}

func (bot *Gobot) startAdaperts() {
	for name, adapter := range bot.adapters {
		log.Infof("Start Adapter %s", name)
		go adapter.Start()
	}
}

func (bot *Gobot) initAdapter() error {
	for name, adapter := range bot.adapters {
		err := adapter.Init(bot)
		if err != nil {
			return fmt.Errorf("Init Adapter %s Fail. %s", name, err.Error())
		}
	}
	return nil
}

func (bot *Gobot) initWorkers() error {
	for name, worker := range bot.workers {
		err := worker.Init(bot)
		if err != nil {
			return fmt.Errorf("Init worker %s Fail. %s", name, err.Error())
		}
	}
	return nil
}
