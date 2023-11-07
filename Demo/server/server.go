package main

import (
	"river/network"
)

func main() {
	s := network.NewServer("tesst server", "127.0.0.1", 8888)
	s.Server()
}
