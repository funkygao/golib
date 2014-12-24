package stringio

import "testing"

func assertEquals(x interface{}, y interface{}) bool {
    return x == y
}

func assertNotEqual(x interface{}, y interface{}) bool {
    return x != y
}

func simpleCreateAndSeek(b []byte) (sio *stringIO) {
    sio = StringIO()
    _, _ = sio.Write(b)
    sio.Seek(0, 0)
    return
}

func TestFd(t *testing.T) {
    sio := StringIO()
    fd, err := sio.Fd()
    if assertEquals(fd, -1) == false {
        t.Errorf("Invalid fd returns")
    }
    if assertEquals(err, nil) {
        t.Errorf("Invalid fd does not raise Error")
    }
    sio.Close()
}

func TestGoName(t *testing.T) {
    sio := StringIO()
    if assertEquals(sio.Name(), sio.GoString()) == false {
        t.Errorf("StringIO name mismatch. Name: %s,"+
            "GoString: %s",
            sio.Name(),
            sio.GoString())
    }
    sio.Close()
}

func TestEmptyStringIO(t *testing.T) {
    sio := StringIO()
    if assertEquals(sio.String(), "") == false {
        t.Errorf("Empty StringIO string is not \"\"")
    }
    sio.Close()
}

func TestClose(t *testing.T) {
    sio := StringIO()
    sio.Close()
    if assertNotEqual(sio.Name(), "StringIO <closed>") {
        t.Errorf("Closed StringIO name invalid")
    }
}

func TestWriteBytes(t *testing.T) {
    b := []byte{'a', 'b', 'c'}
    sio := StringIO()
    n, err := sio.Write(b)
    sio.Seek(0, 0)
    if assertEquals(n, 0) || assertNotEqual(n, 3) {
        t.Errorf("Error writing bytes")
    }
    if assertNotEqual(err, nil) {
        t.Errorf("Write bytes returned error: %s", err)
    }
    sio.Close()
}

func TestSeek(t *testing.T) {
    sio := simpleCreateAndSeek([]byte{'a', 'b', 'c'})
    if assertNotEqual(sio.String(), "abc") {
        t.Errorf("Error Seek result")
    }
}

func TestGetValueString(t *testing.T) {
    sio := simpleCreateAndSeek([]byte{'a', 'b', 'c'})
    s := sio.GetValueString()
    if assertNotEqual(s, sio.String()) {
        t.Errorf("GetValueString return invalid value")
    }
    sio.Close()
    s = sio.GetValueString()
    if assertNotEqual(s, "<nil>") {
        t.Errorf("GetValueString return not <nil> " +
            "on closed StringIO object")
    }
}

func TestGetValueBytes(t *testing.T) {
    sio := simpleCreateAndSeek([]byte{'a', 'b', 'c'})
    bb := sio.GetValueBytes()
    if assertNotEqual(int(bb[0]), int('a')) ||
        assertNotEqual(int(bb[1]), int('b')) ||
        assertNotEqual(int(bb[2]), int('c')) {
        t.Errorf("GetValueBytes return invalid value")
    }
    sio.Close()
    bb = sio.GetValueBytes()
    if assertNotEqual(len(bb), 0) {
        t.Errorf("GetValueBytes return non-zero byte " +
            "array on closed StringIO")
    }
}

func TestStringDiffGetValueString(t *testing.T) {
    sio := simpleCreateAndSeek([]byte{'a', 'b', 'c', 'd', 'e'})
    sio.Seek(3, 0)
    s := sio.GetValueString()
    ss := sio.String()
    if assertEquals(s, ss) {
        t.Errorf("GetValueString and String return the " +
            "same string value when pos is not 0")
    }
}

func TestReadBytes(t *testing.T) {
    b := make([]byte, 3)
    sio := simpleCreateAndSeek([]byte{'a', 'b', 'c'})
    sio.Read(b) // read 3 bytes, pos is now == last
    if assertNotEqual(int(b[0]), int('a')) ||
        assertNotEqual(int(b[1]), int('b')) ||
        assertNotEqual(int(b[2]), int('c')) {
        t.Errorf("Read return invalid value")
    }
    sio.Close()
}

func TestOutBoundRead(t *testing.T) {
    b := make([]byte, 3)
    sio := simpleCreateAndSeek([]byte{'a', 'b', 'c'})
    sio.Seek(1, 2) // set pos to the end of buf
    n, err := sio.Read(b)
    // because pos is now beyond the last write,
    // it returns 0 read and no err.
    if assertNotEqual(n, 0) {
        t.Errorf("Outbound read returns non-zero read count")
    }
    if assertEquals(err, nil) {
        t.Errorf("Outbound read does not return EOF")
    }
    b = make([]byte, 1)
    sio.Seek(2, 0) // Get the pos back inside last
    n, err = sio.Read(b)
    if assertNotEqual(n, 1) {
        t.Errorf("Read return invalid value")
    }
    if assertNotEqual(err, nil) {
        t.Errorf("Read return error: %s", err)
    }
    assertEquals(err, nil)
    if assertNotEqual(int(b[0]), 'c') {
        t.Errorf("Read incorrect value")
    }
    sio.Close()
}

func TestResizeWrite(t *testing.T) {
    b := make([]byte, 5000) // bigger than the buffer
    for i, _ := range b {
        b[i] = 'a'
    }
    sio := simpleCreateAndSeek(b)
    p := make([]byte, 5000)
    n, err := sio.Read(p)
    if assertNotEqual(n, 5000) {
        t.Errorf("Resize write error")
    }
    if assertNotEqual(err, nil) {
        t.Errorf("Resize write error: %s", err)
    }
}

func TestReadWriteString(t *testing.T) {
    s := "This is a test\n"
    sio := StringIO()
    sio.WriteString(s)
    sio.Seek(0, 0)
    p := make([]byte, len(s))
    _, err := sio.Read(p)
    if assertNotEqual(string(p), s) {
        t.Errorf("WriteString error: %s", err)
    }
}

func TestReadWriteUnicodeString(t *testing.T) {
    s := "今天是 2009年 12月 13日 星期日\n"
    sio := StringIO()
    sio.WriteString(s)
    sio.Seek(0, 0)
    p := make([]byte, len(s))
    n, err := sio.Read(p)
    if assertNotEqual(len(s), n) || assertNotEqual(string(p), s) ||
       assertNotEqual(err, nil) {
        t.Errorf("Write Unicode string error: %s", err)
    }
}

func TestErrorOnClosedStringIO(t *testing.T) {
    sio := StringIO()
    sio.Close()
    n, err := sio.Seek(0, 0)
    if assertNotEqual(n, int64(0)) {
        t.Errorf("Seek on closed StringIO error: %s", err)
    }
}

func TestTruncate(t *testing.T) {
    sio := simpleCreateAndSeek([]byte{'a', 'b', 'c', 'd', 'e'})
    sio.Truncate(3)
    s := sio.GetValueString()
    ss := sio.String()
    if assertNotEqual(s, ss) || assertNotEqual(s, "abc") ||
        assertNotEqual(ss, "abc") {
        t.Errorf("Truncate failed")
    }
    sio.Truncate(0)
    s = sio.GetValueString()
    ss = sio.String()
    if assertNotEqual(s, "") || assertNotEqual(ss, "") {
        t.Errorf("Truncate to 0 failed")
    }
}

func TestSeekPosTruncate(t *testing.T) {
    sio := simpleCreateAndSeek([]byte{'a', 'b', 'c', 'd', 'e'})
    sio.Seek(2, 0)
    sio.Truncate(2)
    s := sio.GetValueString()
    ss := sio.String()
    if assertNotEqual(s, "abcd") || assertNotEqual(ss, "cd") {
        t.Errorf("Truncate failed")
    }
}

func TestReadAt(t *testing.T) {
    sio := simpleCreateAndSeek([]byte{'a', 'b', 'c', 'd', 'e'})
    b := make([]byte, 2)
    sio.ReadAt(b, 3)
    if assertNotEqual(string(b), "de") {
        t.Errorf("ReadAt returns invalid value: %s", string(b))
    }
}

func TestWriteAt(t *testing.T) {
    sio := simpleCreateAndSeek([]byte{'a', 'b'})
    b := []byte{'c', 'd'}
    sio.WriteAt(b, 2)
    s := sio.GetValueString()
    if assertNotEqual(s, "abcd") {
        t.Errorf("WriteAt append failed")
    }
    sio.WriteAt(b, 0)
    s = sio.GetValueString()
    if assertNotEqual(s, "cdcd") {
        t.Errorf("WriteAt overwrite failed")
    }
    // not allow to seek outside of buf size,
    // but allow to create byte array hold within
    // the buf size. The output string
    sio.WriteAt(b, 100)
    s = sio.GetValueString()
    bb := sio.GetValueBytes()
    if assertNotEqual(len(bb), 102) {
        t.Errorf("WriteAt byte array hole failed")
    }
    if assertNotEqual(len(s), 102) {
        t.Errorf("WriteAt byte array hole with invalid string length")
    }
}
