package event

type Event interface {
	Subscribe(topic string, tag string, consumer string, body interface{}, handler func(interface{}) error) error
}