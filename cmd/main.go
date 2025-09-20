// @cyptocurrency deliver API
// @version 1.0
// @description REST API for delivering cryptocurrencies
// @host localhost:8080
// @BasePath /
package main

import (
	"log"
	"os"

	"github.com/100bench/cryptocurrency_provider.git/app"
)

func main() {
	if err := app.RunApp(); err != nil {
		log.Printf("Application exited with error: %v", err)
		os.Exit(1)
	}
}
