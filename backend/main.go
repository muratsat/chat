package main

import (
	"os"
	"strconv"

	"github.com/muratsat/chat/api"
)

func main() {
	port, err := strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		port = 8081
	}

	app := api.CreateApp(port)
	app.Run()
}
