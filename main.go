package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/btcsuite/websocket"
)

var ArgAddr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", defaultHandler)
	log.Println("Listenning on", *ArgAddr)
	log.Fatal(http.ListenAndServe(*ArgAddr, nil))
}

// Allows all origins here
func checkOrigin(r *http.Request) bool {
	return true
}

var upgrader = websocket.Upgrader{
	CheckOrigin: checkOrigin,
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	if err != nil {
		log.Println("Error sending message: ", err.Error())
	}

	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		if string(msg) != "" {
			log.Printf("recv: %s", msg)
			log.Printf(fmt.Sprintf("X-Forwarded-For: %s", r.Header.Get("X-Forwarded-For")))

			err = c.WriteMessage(mt, []byte(fmt.Sprintf("X-Forwarded-For: %s", r.Header.Get("X-Forwarded-For"))))
			err = c.WriteMessage(mt, []byte(fmt.Sprintf("pong: %s", msg)))
			if err != nil {
				log.Println("Error sending message: ", err.Error())
			}
		}
	}
}
