package main

func main() {
	go StartHealthMonitor("node1:9000")
	StartHttpServer()
}
