package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/ant0ine/go-json-rest/rest"
)

func main() {
	// HTMLページを表示
	// http.Handle("/", http.FileServer(http.Dir("static")))
	// if err := http.ListenAndServe(":8080", nil); err != nil {
	// 	log.Fatal("ListenAndServe: ", err)
	// }

	// API起動
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Post("/countries", PostCountry),
		rest.Get("/countries/:code", GetCountry),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}

// Country 国codeと名前
type Country struct {
	Code string
	Name string
}

// ここに一時的にPOSTされたデータを保存しているんだな
var store = map[string]*Country{}

// これ何？
var lock = sync.RWMutex{}

// PostCountry 国を投稿する
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

// GetCountry 国々を取得する
func GetCountry(w rest.ResponseWriter, r *rest.Request) {
	code := r.PathParam("code")

	lock.RLock()
	var country *Country
	if store[code] != nil {
		country = &Country{}
		*country = *store[code]
	}
	lock.RUnlock()

	if country == nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(country)
}