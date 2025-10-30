package main

import (
	"bender/internal/bot"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		panic(err)
	}
}

func main() {
	err := bot.Start()

	if err != nil {
		return
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
