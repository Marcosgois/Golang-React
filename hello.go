package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type Twitter struct {
	ID      string
	User    string
	Message string
	Data    string
	Time    time.Time
}

type twitterHandlers struct {
	sync.Mutex
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
	h.Lock()
	i := 0
	for _, twitter := range h.store {
		twitters[i] = twitter
		i++
	}
	h.Unlock()

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

	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("Need content-type 'application/json', but got '%s'", ct)))
		return
	}

	var twitter Twitter
	err = json.Unmarshal(bodyBytes, &twitter)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	twitter.Time = time.Now()

	h.Lock()
	h.store[twitter.ID] = twitter
	w.Write([]byte(fmt.Sprintf("Postando ID: '%s'", twitter.ID)))
	defer h.Unlock()
}

func newTwitterHandlers() *twitterHandlers {
	return &twitterHandlers{
		store: map[string]Twitter{},
	}
}

func main() {
	twitterHandlers := newTwitterHandlers()
	http.HandleFunc("/twitters", twitterHandlers.twitters)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
