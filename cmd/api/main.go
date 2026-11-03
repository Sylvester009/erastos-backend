package main

import "erastos-backend/internal/api"

func main() {
	router := api.BookRouter()
	router.SetTrustedProxies(nil)
	router.Run("localhost:8080")
}
