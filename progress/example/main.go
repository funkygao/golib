package main

import (
    bar "github.com/funkygao/golib/progress"
    "time"
)

func main() {
    const N int = 10
    p := bar.New(N)
    for i:=0; i<N; i++ {
        p.ShowProgress(i)
        time.Sleep(time.Second)
    }
}
