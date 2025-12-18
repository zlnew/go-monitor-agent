// Package subscribers
package subscribers

type Subscriber interface {
	Handle(event any)
}

type EventBus interface {
	Subscribe(eventName string, handler func(event any))
}
