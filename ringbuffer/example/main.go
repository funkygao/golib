package main

import (
    "time"
    "fmt"
    "github.com/textnode/gringo"
)

var size uint64 = 5000000

func noopPayload(param gringo.Payload) {}

var pl gringo.Payload = *gringo.NewPayload(1)

func gringoProducer(outGringo *gringo.Gringo, done chan int) {
    var i uint64
    for ; i < size; i++ {
        outGringo.Write(pl)
    }
    done <- 0
}

func gringoForwarder(inGringo *gringo.Gringo, outGringo *gringo.Gringo, done chan int) {
    var i uint64
    for ; i < size; i++ {
        var rez gringo.Payload = inGringo.Read()
        outGringo.Write(rez)
    }
    done <- 0
}

func gringoConsumer(inGringo *gringo.Gringo, done chan int) {
    var i uint64
    for ; i < size; i++ {
        var rez gringo.Payload = inGringo.Read()
        noopPayload(rez)
    }
    done <- 0
}

func gringoRunner() {
    doneChan := make(chan int)

    var r1 *gringo.Gringo = gringo.NewGringo()
    var r2 *gringo.Gringo = gringo.NewGringo()

    var startTime time.Time = time.Now()

    go gringoProducer(r1, doneChan)
    go gringoForwarder(r1, r2, doneChan)
    go gringoConsumer(r2, doneChan)

    <-doneChan; <-doneChan; <-doneChan

    fmt.Println("gringoRunner seconds passed:", time.Since(startTime))
}


func chanProducer(outChan chan gringo.Payload, done chan int) {
    var i uint64
    for ; i < size; i++ {
        outChan <- pl
    }
    done <- 0
}

func chanForwarder(inChan chan gringo.Payload, outChan chan gringo.Payload, done chan int) {
    var i uint64
    for ; i < size; i++ {
        o := <- inChan
        outChan <- o
    }
    close(inChan)
    done <- 0
}

func chanConsumer(inChan chan gringo.Payload, done chan int) {
    var i uint64
    for ; i < size; i++ {
        o := <- inChan
        noopPayload(o)
    }
    close(inChan)
    done <- 0
}

func chanRunner() {
    doneChan := make(chan int)

    r1 := make(chan gringo.Payload, 4096)
    r2 := make(chan gringo.Payload, 4096)

    var startTime time.Time = time.Now()

    go chanProducer(r1, doneChan)
    go chanForwarder(r1, r2, doneChan)
    go chanConsumer(r2, doneChan)

    <-doneChan; <-doneChan; <-doneChan

    fmt.Println("chanRunner seconds passed:", time.Since(startTime))
}

func main() {
    gringoRunner()
    gringoRunner()
    chanRunner()
    chanRunner()
    gringoRunner()
    chanRunner()
}
