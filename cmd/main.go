package main

import (
	"fmt"
	"log"

	"github.com/Dhairya3124/e-commerce-coupon-system/internal/db"
)

func main() {
	dbConn := db.NewDB()
	fmt.Println(dbConn, "done")
	log.Fatal(dbConn)

}
