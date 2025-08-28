package main

import (
	"fmt"
	"log"

	"github.com/Dhairya3124/e-commerce-coupon-system/internal/db"
	"github.com/Dhairya3124/e-commerce-coupon-system/internal/repository"
)

func main() {
	dbConn := db.NewDB()
	repo := repository.NewCouponRepository(dbConn.DB)

	fmt.Println(dbConn, repo, "done")
	log.Fatal(dbConn)

}
