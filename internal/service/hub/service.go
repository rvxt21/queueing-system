package hub

import (
	"sync"

	"github.com/rvxt21/queueing-system/internal/domain"
)

// subscriberBuffer is how many updates a slow client can fall behind before updates are dropped.
const subscriberBuffer = 16

type Subscriber struct {
	msgs chan domain.TicketStatusChanged
}

type Hub struct {
	mu          sync.Mutex
	subscribers map[*Subscriber]struct{}
}

func NewHub() *Hub {
	return &Hub{
		subscribers: make(map[*Subscriber]struct{}),
	}
}

func (h *Hub) Subscribe() *Subscriber {
	sub := &Subscriber{msgs: make(chan domain.TicketStatusChanged, subscriberBuffer)}

	h.mu.Lock()
	defer h.mu.Unlock()

	h.subscribers[sub] = struct{}{}

	return sub
}

func (h *Hub) Unsubscribe(sub *Subscriber) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.subscribers, sub)
}

func (h *Hub) Publish(msg domain.TicketStatusChanged) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for sub := range h.subscribers {
		select {
		case sub.msgs <- msg:
		default: // client too slow, drop the update
		}
	}
}

func (s *Subscriber) Messages() <-chan domain.TicketStatusChanged {
	return s.msgs
}
