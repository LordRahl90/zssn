package main

import (
	"log"
	"zssn/servers"

	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	server := servers.New(db)

	if err := server.Start(":3500"); err != nil {
		log.Fatal(err)
	}
}
