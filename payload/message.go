package payload

type Message struct {
	Id       string
	FromUser *User
	Text     string
}
