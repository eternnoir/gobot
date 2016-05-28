package gobot

type Adapter interface {
	Init() error

	Start() error

	Response() error
}
