package main

import (
	"math/rand"
	"time"

	"github.com/antonvlasov/avito/handler"
	"github.com/antonvlasov/avito/httpserver"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}
func main() {

	handler.Run(10)
	httpserver.Run()
}
