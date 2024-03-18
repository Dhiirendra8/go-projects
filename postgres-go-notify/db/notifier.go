package db

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/lib/pq"
)

func ListenForNotifications(connString string, channel string) {
	// db, err := Connect(connString)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	listener := CreateListener(connString, channel)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	notification := make(chan *pq.Notification)

	go func() {
		for {
			select {
			case <-sig:
				fmt.Println("Exiting...")
				return
			case notify := <-notification:
				fmt.Println("Received notification:")
				fmt.Println("  Channel:", notify.Channel)
				fmt.Println("  Payload:", notify.Extra)
			}
		}
	}()

	errorCh := make(chan error, 1)
	for {
		select {
		case <-errorCh:
			fmt.Println("Stopping listener....")
			return
		default:
			err := listener.Ping()
			if err != nil {
				log.Println("Error pinging listener:", err)
				errorCh <- err
			}

			select {
			case <-sig:
				fmt.Println("Stopping listener...")
				return
			case n := <-listener.Notify:
				notification <- n
			case <-time.After(time.Second * 300):
				return
			}
		}
	}
}
