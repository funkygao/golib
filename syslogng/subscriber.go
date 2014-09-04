package syslogng

import (
    "bufio"
    "os"
)

func Subscribe() chan string {
    bio := bufio.NewReader(os.Stdin)
    ch := make(chan string)
    go func() {
        for {
            line, err := readline(bio)
            if err != nil {
            }

            ch <- string(line)
        }
    }()

    return ch
}

func readline(bio *bufio.Reader) ([]byte, error) {
    line, isPrefix, err := bio.ReadLine()
    if !isPrefix {
        return line, err
    }

    // line is too long, read till eol
    buf := append([]byte(nil), line...)
    for isPrefix && err == nil {
        line, isPrefix, err = bio.ReadLine()
        buf = append(buf, line...)
    }
    return buf, err
}
