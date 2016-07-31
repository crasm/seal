package main

import (
	log "github.com/go-ozzo/ozzo-log"
)

var logger *log.Logger

func init() {
	logger = log.NewLogger()
	logger.Targets = []log.Target{log.NewConsoleTarget()}
	logger.Open()
}
