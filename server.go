package main

import (
	"fmt"
	_"net/http"
	"os"
    "os/signal"
	"github.com/mlo77/webobs"
    "bufio"
)

var wobs *webobs.Server

func exit() {
    os.Exit(0)
}

func playerMsg(tag string, data []byte) {
    wobs.WriteCh<-webobs.Message{Tag:"screen", Data: data }
}

func listenInput(allDone chan bool) {
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        fmt.Println(scanner.Text())
        wobs.ShutClientsForTag("player")
    }
    // for {
    //     select {
    //     case scanner.Scan() :
    //         fmt.Println(scanner.Text())
    //     case <-allDone :
    //         return
    //     }
    // }
}

func main() {
	// creates a http server on 8080 
	wobs = webobs.StartServer(":8080")

	// set a channel ready to be connected
	// from the address http://localhost:8080/fooTag
    wobs.SetChannel("screen", nil, "floppybird")
    wobs.SetChannel("player", playerMsg, "player")

    // capture ctrl+c, so we can exit properly
    allDone := make(chan bool)
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, os.Interrupt)
    go listenInput(allDone)

    go func() {
        <-sig // blocks until something is read in the channel
        allDone<-true
    }()

    <-allDone

    fmt.Println("exit")
}
