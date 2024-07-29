package main

func main() {
	initRedis()
	RunTraceRouteServer(":8080")
}
