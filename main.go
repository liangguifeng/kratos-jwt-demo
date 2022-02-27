package main

import (
	"github.com/go-kratos/kratos/v2/log"
	"project/startup"
)

func main() {
	app, err := startup.Init()
	if err != nil {
		log.Fatalf("app run err: %v", err)
	}

	if err = app.Run(); err != nil {
		log.Fatal(err)
	}
}
