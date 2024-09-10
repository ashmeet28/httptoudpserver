package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	httpToUdpCh := make(chan string)

	var wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	var reqBalancer []int
	var reqBalancerMu sync.Mutex

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/" {
			c, err := wsUpgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Println(err)
				return
			}

			reqBalancerMu.Lock()

			connIndex := len(reqBalancer)
			reqBalancer = append(reqBalancer, 0)

			reqBalancerMu.Unlock()

			for {
				messageType, p, err := c.ReadMessage()
				if err != nil {
					reqBalancerMu.Lock()
					reqBalancer[connIndex] = 2
					reqBalancerMu.Unlock()

					log.Println(err)
					return
				}

				for {
					canSend := false

					reqBalancerMu.Lock()

					isPending := false

					for _, s := range reqBalancer {
						if s == 0 {
							isPending = true
							break
						}
					}

					if !isPending {
						for i, s := range reqBalancer {
							if s == 1 {
								reqBalancer[i] = 0
							}
						}
					}

					if reqBalancer[connIndex] == 0 {
						canSend = true
					}

					reqBalancerMu.Unlock()

					if canSend {
						if err := c.WriteMessage(messageType, p); err != nil {
							reqBalancerMu.Lock()
							reqBalancer[connIndex] = 2
							reqBalancerMu.Unlock()

							log.Println(err)
							return
						} else {
							reqBalancerMu.Lock()
							reqBalancer[connIndex] = 1
							reqBalancerMu.Unlock()
							break
						}
					}

				}
			}
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
