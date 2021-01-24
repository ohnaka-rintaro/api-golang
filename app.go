package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/sqs/goreturns/returns"
)

func main() {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Post("/countries", PostCountry)
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}

type Country struct {
	Code string
	Name string
}

var store = map[string]*Country{}

// これ何？
var lock = sync.RWMutex{}

func PostCountry(w rest.ResponseWriter, r *rest.Request) {
	country := Country{}
	err := r.DecodeJsonPayload(&country)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if country.Code == "" {
		rest.Error(w, "country code required", 400)
		return
	}
	if country.Name == "" {
		rest.Error(w, "country name required", 400)
		return
	}
	// 意味がよくわからない
	lock.Lock()
	store[country.Code] = &country
	lock.Unlock()
	w.WriteJson(&country)
}