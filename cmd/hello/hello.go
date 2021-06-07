package main

import (
    "fmt"
    log "github.com/sirupsen/logrus"
)

func main() {
    fmt.Println("hello world")
    log.Trace("Something very low level.")
    log.Debug("Useful debugging information.")
    log.Info("Something noteworthy happened!")
    log.Warn("You should probably take a look at this.")
    log.Error("Something failed but I'm not quitting.")
}