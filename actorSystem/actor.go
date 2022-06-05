package actorSystem

type Actor[T any] struct {
	incoming    chan<- T
	stopChannel chan<- interface{}
}

func (self Actor[T]) Send(data T) {
	self.incoming <- data
}

func (self Actor[T]) Stop() {
	self.stopChannel <- nil
}

type fn func(any)

func NewActor[T any](bufferSize uint, receiver func(T)) Actor[T] {
	incoming := make(chan T, bufferSize)
	stop := make(chan interface{})
	a := Actor[T]{
		incoming:    incoming,
		stopChannel: stop,
	}
	go func() {
		for {
			select {
			case <-stop:
				return
			case message := <-incoming:
				receiver(message)
			}
		}
	}()
	return a
}