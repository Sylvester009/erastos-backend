package main

func main() {
	router := bookRouter()
	router.Run("localhost:8080")
}
