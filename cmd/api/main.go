package main

import (
	"swift-menu-session/app"
)

func main() {
	a := app.SetupAPP()
	app.SetupHandlers(a)
	select {}
}
