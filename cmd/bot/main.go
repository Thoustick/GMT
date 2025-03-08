package main

import ( 
	"os"

	"github.com/Thoustick/GMT/pkg/logger"
	"github.com/Thoustick/GMT/internal/bot"
)

func main() {
	err := bot.Init()
	if err !=nil {
		logger.Fatal(err, "Error while initializing bot")
		os.Exit(1)
	}
}