package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/muratsat/chat/backend/api/handlers"
)

type App struct {
	port int
}

func CreateApp(port int) App {
	var app App
	app.port = port
	return app
}

func (app *App) Run() {
	handlers.SetupHandlers()
	url := fmt.Sprintf(":%d", app.port)
	err := http.ListenAndServe(url, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
