package db

import (
	"log"
	"time"

	"github.com/lib/pq"
)

func CreateListener(connString string, channel string) *pq.Listener {
	listener := pq.NewListener(connString, 10*time.Second, time.Minute, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Println("Listener error:", err)
		}
	})

	err := listener.Listen(channel)
	if err != nil {
		log.Fatal(err)
	}

	return listener
}
