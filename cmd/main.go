package main

import (
	"github.com/Joe5451/go-oauth2-server/internal"
	"github.com/Joe5451/go-oauth2-server/internal/config"

	"fmt"
)

func main() {
	if err := config.InitializeAppConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	router, err := internal.InitializeApp()
	if err != nil {
		panic(err)
	}

	router.Run("localhost:8080")
}
