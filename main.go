package main

import (
	"fridge/src/controller"
)

func main() {
	var restServer = controller.RestServer{}

	restServer.RunRESTServer()
}
