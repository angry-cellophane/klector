package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io.klector/klector/storage"
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
	var event storage.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}

	(*s.storage).Write(&event)
	w.WriteHeader(204)
}

func (s *server) query(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var query storage.Query

	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}

	resultSet := (*s.storage).Query(&query)

	if err := json.NewEncoder(w).Encode(resultSet); err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(200)
}

func (s *server) start() error {
	s.router.POST("/api/v1/event", s.store)
	s.router.GET("/api/v1/query", s.query)

	return http.ListenAndServe(":4479", s.router)
}

func Create(storage *storage.Storage) error {
	server := server{
		router:  httprouter.New(),
		storage: storage,
	}
	return server.start()
}
