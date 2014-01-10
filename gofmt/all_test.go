package gofmt

import (
    "fmt"
    "testing"
)

func TestByteSize(t *testing.T) {
    type testCase struct {
        b        ByteSize
        expected string
    }

    cases := []testCase{
        {638048, "623.09KB"},
        {12121212212, "11.29GB"}}
    for _, b := range cases {
        s := fmt.Sprintf("%s", b.b)
        if s != b.expected {
            t.Error("exptected:", b.expected, " real:", s)
        }
    }
}

func TestByteSizeConsts(t *testing.T) {
    if MB/KB != 1024 {
        t.Error("MB/KB != 1024")
    }

    if KB != 1024 {
        t.Error("KB")
    }
    if MB != 1024*KB {
        t.Error("MB", int64(MB))
    }
    if GB != 1024*MB {
        t.Error("GB", int64(GB))
    }
    if TB != 1024*GB {
        t.Error("TB", int64(TB))
    }
    if PB != 1024*TB {
        t.Error("PB", int64(PB))
    }
    if EB != 1024*PB {
        t.Error("EB", int64(EB))
    }
}

func TestComma(t *testing.T) {
    expected := "33,434,433"
    got := Comma(33434433)
    if got != expected {
        // echo got result
        t.Error("Expected:", expected, "Got:", got)
    }
}
