package main

import (
	"GoPostgresNotifier/db"
	"fmt"
)

const (
	connString = "user=postgres password=password dbname=mydb sslmode=disable"
	channel    = "myChannel"
)

func main() {
	fmt.Println("Listening to channel :", channel)
	db.ListenForNotifications(connString, channel)
}
