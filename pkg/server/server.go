package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/ocramh/challenge2/pkg/content"
	"github.com/ocramh/challenge2/pkg/indexer"
	"github.com/ocramh/challenge2/pkg/provider"
)

const (
	addRoute = "/add"
	getRoune = "/get"
)

type Server struct {
	srv  http.Server
	prv  *provider.Provider
	port int
}

func NewLocalServer(port int, prv *provider.Provider) *Server {
	return &Server{
		port: port,
		prv:  prv,
	}
}

func (s Server) addItemHandler(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("item")

	b, err := s.prv.AddItem([]byte(item))
	if err != nil {
		log.Fatalf(err.Error())
	}

	j, err := json.Marshal(b)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Fprint(w, string(j))
}

func (s Server) getItemHandler(w http.ResponseWriter, r *http.Request) {
	blockCid := r.URL.Query().Get("cid")
	blockData, err := s.prv.GetItem(blockCid)
	if err != nil {
		var status int
		if errors.Is(err, indexer.ErrNoItemFound) {
			status = http.StatusNotFound
		} else if errors.Is(err, content.ErrInvalidCID) {
			status = http.StatusBadRequest
		} else {
			status = http.StatusInternalServerError
		}

		w.WriteHeader(status)
		fmt.Fprint(w, fmt.Sprintf(`{"error": "%s" }`, err.Error()))
		return
	}

	j, err := json.Marshal(blockData)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Fprint(w, string(j))
}

func (s *Server) Start() {
	http.HandleFunc(addRoute, s.addItemHandler)
	http.HandleFunc(getRoune, s.getItemHandler)

	log.Printf("server running at localhost:%d \n", s.port)
	http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}
