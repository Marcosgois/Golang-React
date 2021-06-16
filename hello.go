package main

import (
	"encoding/json"
	"net/http"
)

type Twitter struct {
	User    string
	Message string
	Data    string
	ID      string
}

type twitterHandlers struct {
	store map[string]Twitter
}

func (h *twitterHandlers) twitters(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
		return
	case "POST":
		h.post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Não Permitido"))
	}
}

func (h *twitterHandlers) get(w http.ResponseWriter, r *http.Request) {
	twitters := make([]Twitter, len(h.store))

	i := 0
	for _, twitter := range h.store {
		twitters[i] = twitter
		i++
	}

	jsonBytes, err := json.Marshal(twitters)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
func (h *twitterHandlers) post(w http.ResponseWriter, r *http.Request) {

}
func newTwitterHandlers() *twitterHandlers {
	return &twitterHandlers{
		store: map[string]Twitter{
			"id1": {
				User:    "Marcos Gois",
				Message: "Hello World",
				Data:    "16/06/2021",
				ID:      "id1",
			},
		},
	}
}

func main() {
	twitterHandlers := newTwitterHandlers()
	http.HandleFunc("/twitter", twitterHandlers.get)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
