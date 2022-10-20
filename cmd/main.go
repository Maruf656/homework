package main

import (
	"github.com/Abdulhalim92/server/internal/router"
	"log"
	"os"
)

func main() {
	err := router.StartRouter()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

//REST -- что это такое
