package api

import (
	"net/http"
)

type Broker struct {
	Clients    map[chan string]bool
	deadClient chan chan string
	newClient  chan chan string
	message    chan string
}

func NewBroker() *Broker {
	return &Broker{
		Clients:    make(map[chan string]bool),
		deadClient: make(chan chan string),
		newClient:  make(chan chan string),
		message:    make(chan string),
	}
}

func (b *Broker) Start() {
	go func() {
		for {
			select {
			case client := <-b.newClient:
				b.Clients[client] = true
			case client := <-b.deadClient:
				delete(b.Clients, client)
				close(client)
			case msg := <-b.message:
				for client := range b.Clients {
					client <- msg
				}
			}
		}
	}()
}

func (b *Broker) Stop() {
	close(b.deadClient)
	close(b.newClient)
	close(b.message)
}

func (b *Broker) Broadcast(msg string) {
	b.message <- msg
}

func (b *Broker) AddClient(client chan string) {
	b.newClient <- client
}

func (b *Broker) RemoveClient(client chan string) {
	b.deadClient <- client
}

func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Configure SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Verify streaming support
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming is not supported", http.StatusInternalServerError)
		return
	}

	clientChan := make(chan string)
	b.AddClient(clientChan)

	// Notify client disconnection
	notifyDesconnection := r.Context().Done()

	go func() {
		defer b.RemoveClient(clientChan)
		<-notifyDesconnection
	}()

	for msg := range clientChan {
		w.Write([]byte("data: " + msg + "\n\n"))
		flusher.Flush()
	}

}
