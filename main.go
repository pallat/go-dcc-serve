package main

import (
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
	"dccserve/balance"
)

func main() {
	go balance.Start()
	api := rest.NewApi()
	// api.Use(rest.DefaultDevStack...)
	api.Use(&rest.CorsMiddleware{
		RejectNonCorsRequests: false,
		OriginValidator: func(origin string, request *rest.Request) bool {
			return true
		},
		AllowedMethods:                []string{"GET", "OPTIONS"},
		AllowedHeaders:                []string{"Accept", "Content-Type", "X-Requested-With", "X-Custome-Header", "Origin"},
		AccessControlAllowCredentials: true,
		AccessControlMaxAge:           3600,

		// RejectNonCorsRequests: false,
		// OriginValidator: func(origin string, request *rest.Request) bool {
		// 	return true //origin == "http://my.other.host"
		// },
		// AllowedMethods: []string{"OPTIONS", "GET", "POST", "PUT"},
		// AllowedHeaders: []string{
		// 	"Accept", "Content-Type", "X-Custom-Header", "Origin"},
		// AccessControlAllowCredentials: true,
		// AccessControlMaxAge:           3600,
	})
	router, err := rest.MakeRouter(
		&rest.Route{"GET", "/balance/:corp/:subr", balance.Balance},
		// &rest.Route{"GET", "/dtn/:subr", dserve.DTNBalance},
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)

	fmt.Println("Start api 8088")
	log.Fatal(http.ListenAndServe(":8088", api.MakeHandler()))
}
