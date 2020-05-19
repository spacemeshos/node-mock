package utils

// Bus -
type Bus struct {
	channels map[int]chan interface{}
	cookie   int
}

// Init -
func (b *Bus) Init() {
	b.channels = make(map[int]chan interface{})
}

// Register -
func (b *Bus) Register() (chan interface{}, int) {
	channel := make(chan interface{})

	b.channels[b.cookie] = channel

	b.cookie++

	return channel, b.cookie - 1
}

// Delete -
func (b *Bus) Delete(cookie int) {
	delete(b.channels, cookie)
}

// Send -
func (b *Bus) Send(msg interface{}) {
	for _, c := range b.channels {
		c <- msg
	}
}
