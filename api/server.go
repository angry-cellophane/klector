package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io.klector/klector/storage"
	"log"
	"net/http"
)

type Api interface {
	Stop()
}

type server struct {
	router  *httprouter.Router
	storage *storage.Storage
}

func (s *server) Stop() {
	// noop
}

func (s *server) store(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var events storage.Events
	if err := json.NewDecoder(r.Body).Decode(&events); err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	if events.Events == nil || len(events.Events) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("no events"))
		return
	}

	if err := (*s.storage).Write(&events); err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(204)
}

func (s *server) query(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var query storage.Query

	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
	log.Printf("received query %v", query)

	resultSet, err := (*s.storage).Query(&query)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	if err := json.NewEncoder(w).Encode(resultSet); err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
}

func (s *server) start() error {
	s.router.POST("/api/v1/event", s.store)
	s.router.POST("/api/v1/query", s.query)

	return http.ListenAndServe(":4479", s.router)
}

func Create(storage *storage.Storage) error {
	server := server{
		router:  httprouter.New(),
		storage: storage,
	}
	return server.start()
}
