package connectors

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
)

type HttpEndpointConnector struct {
	data any
	mu   sync.RWMutex
}

func NewHttpEndpointConnector(r *chi.Mux, route string) *HttpEndpointConnector {
	hec := HttpEndpointConnector{}

	r.Get(route, func(w http.ResponseWriter, r *http.Request) {
		hec.mu.RLock()
		payload, err := json.Marshal(hec.data)
		hec.mu.RUnlock()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(payload)
	})

	return &hec
}

func (hec *HttpEndpointConnector) Send(data any) error {
	hec.mu.Lock()
	defer hec.mu.Unlock()

	hec.data = data
	return nil
}

func (hec *HttpEndpointConnector) GetData() any {
	hec.mu.RLock()
	defer hec.mu.RUnlock()

	return hec.data
}
