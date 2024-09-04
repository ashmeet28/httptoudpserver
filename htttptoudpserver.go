package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("GET /sendpacket", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Methods", "*")
		w.Header().Add("Access-Control-Expose-Headers", "*")

		if r.Method == http.MethodGet && r.URL.Path == "/sendpacket" {

			d := r.Header.Get("httptoudpserver-content")
			if d != "" {
				fmt.Println(d)
			}
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	http.HandleFunc("OPTIONS /sendpacket", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Methods", "*")
		w.Header().Add("Access-Control-Expose-Headers", "*")

		if r.URL.Path == "/sendpacket" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	log.Fatal(http.ListenAndServe(":3000", nil))
}
