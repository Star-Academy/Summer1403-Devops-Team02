package main

import "log"

func main() {
	err := initRedis()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected to Redis")

	RunTraceRouteServer(":8080")
}
