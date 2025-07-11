package echo

type Notifier interface {
	Notify(message string) error
}

type Handler interface {
	Handle(message string) error
}
