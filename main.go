package main


func main() {
	server := NewAPIServer(":5555")
	server.Run()
}
