package main

import (
	"example/mongodemo/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("MongoDB API")
	r := router.Router()
	fmt.Println("starting server...")
	log.Fatal(http.ListenAndServe(":4050", r))
	fmt.Println("server started at port 4050")
}
