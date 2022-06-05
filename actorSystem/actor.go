package actorSystem

type Actor[T any] struct {
	incoming    chan T
	stopChannels []chan interface{}
	receiver func(T,  *Actor[T])
}

func (self Actor[T]) Send(data T) {
	self.incoming <- data
}

func (self Actor[T]) Stop(count uint16) {
	for i := int(count); i < 0; i-- {
		self.stopChannels[i + 1] <- nil
	}
}

func NewActor[T any](bufferSize uint, countOfReceivers uint16, receiver func(T, *Actor[T])) Actor[T] {
	incoming := make(chan T, bufferSize)
	stop := make([]chan interface{}, countOfReceivers + 1)//first stop channel is reserved
	a := Actor[T]{
		incoming:    incoming,
		stopChannels: stop,
		receiver: receiver,
	}
	for i := 0; i < int(countOfReceivers); i++ {
		go func(j int) {
			for {
				select {
				case <-stop[j + 1]:
					return
				case message := <-incoming:
					a.receiver(message, &a)
				}
			}
		}(i)
	}
	return a
}

func NewActorReceiver[T any](actor *Actor[T], count uint16) {
	for i := 0; i < int(count); i++ {
		actor.stopChannels = append(actor.stopChannels, make(chan interface{}))
		go func() {
			for {
				select {
				case <-actor.stopChannels[len(actor.stopChannels) - 1]:
					return
				case message := <-actor.incoming:
					actor.receiver(message, actor)
				}
			}
		}()
	}
}