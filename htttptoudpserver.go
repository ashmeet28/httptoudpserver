package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	httpToUdpCh := make(chan string)

	http.HandleFunc("GET /sendpacket", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Methods", "*")
		w.Header().Add("Access-Control-Expose-Headers", "*")

		if r.Method == http.MethodGet && r.URL.Path == "/sendpacket" {
			p := r.Header.Get("httptoudpserver-content")

			if p != "" {
				httpToUdpCh <- p
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

	go func() {
		c, err := net.Dial("udp", "127.0.0.1:4242")

		if err != nil {
			os.Exit(1)
		}

		var nextPacketId uint64 = 1

		for {
			c.Write([]byte(((<-httpToUdpCh) + "-" + strconv.FormatUint(nextPacketId, 10))))
			nextPacketId++
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second * 1)
			httpToUdpCh <- "0-0-0-0-0-0-0-0-0"
		}
	}()

	log.Fatal(http.ListenAndServe(":3000", nil))
}
