package payload

type Message struct {
	Id            string
	FromUser      *User
	Text          string
	Payload       interface{}
	SourceAdapter string
}
