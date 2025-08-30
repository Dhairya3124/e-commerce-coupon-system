package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Dhairya3124/e-commerce-coupon-system/internal/api"
	"github.com/Dhairya3124/e-commerce-coupon-system/internal/cache"
	"github.com/Dhairya3124/e-commerce-coupon-system/internal/db"
	"github.com/Dhairya3124/e-commerce-coupon-system/internal/repository"
	"github.com/Dhairya3124/e-commerce-coupon-system/internal/service"

	"github.com/go-chi/chi"
)

func main() {
	dbConn := db.NewDB()
	repo := repository.NewCouponRepository(dbConn.DB)

	cache := cache.NewLRU(100)
	serv := service.NewCouponService(repo, cache)
	handler := api.NewCouponHandler(serv)

	fmt.Println(dbConn, repo, "done")
	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {

		r.Route("/coupons", func(r chi.Router) {
			r.Post("/", handler.CreateCouponHandler)
			r.Get("/", handler.GetAllCouponsHandler)
		})

	})
	srv := http.Server{
		Addr:    ":3000",
		Handler: r,
	}
	log.Fatal(srv.ListenAndServe())

}
