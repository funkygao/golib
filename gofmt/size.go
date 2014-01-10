package gofmt

import "fmt"

type ByteSize float64

func (b ByteSize) String() string {
    switch {
    case b >= EB:
        return fmt.Sprintf("%.2fEB", b/EB)
    case b >= PB:
        return fmt.Sprintf("%.2fPB", b/PB)
    case b >= TB:
        return fmt.Sprintf("%.2fTB", b/TB)
    case b >= GB:
        return fmt.Sprintf("%.2fGB", b/GB)
    case b >= MB:
        return fmt.Sprintf("%.2fMB", b/MB)
    case b >= KB:
        return fmt.Sprintf("%.2fKB", b/KB)
    }

    return fmt.Sprintf("%.2fB", b)
}

const (
    _            = iota
    KB  ByteSize = 1 << (10 * iota)
    MB
    GB
    TB
    PB
    EB
)
