package gobot

type Adapter interface {
	Init() error

	Start()
}
