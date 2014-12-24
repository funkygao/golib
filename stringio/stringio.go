// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stringio

import (
    "errors"
    "fmt"
    "io"
    "syscall"
)

const buf_size = 4096

var OSError error = errors.New("I/O operation on closed file")

// A stringIO object is similar to a File object.
// It mimics all File I/O operations by implementing the
// common interfaces where File is implemented.
// The main difference is StringIO never read/write to filesystem.
// All operations are done in memory by accessing its underly buffer.
// The difference b/w bytes.Buffer is that StringIO supports
// Random access where Buffer does not.
// Most buffer operations are similar to bytes.Buffer.
// StringIO also does not support any non I/O operations such as
// Mkdir, Stats, Symlink, etc which does not have a real semantics for
// buffer manipulation.
// A StringIO object can not be reused once it is closed, just like
// the file object.
type stringIO struct {
    buf       []byte
    isclosed  bool
    pos, last int
    name      string
}

// Factory method served as the constructor.
func StringIO() *stringIO {
    buf := make([]byte, buf_size)
    sio := new(stringIO)
    sio.buf = buf
    sio.isclosed = false
    sio.pos = 0
    sio.last = 0
    sio.name = fmt.Sprintf("StringIO <%p>", sio)
    return sio
}

// Query for stringio object's fd is an error.
func (s *stringIO) Fd() (fd int, err error) {
    return -1, syscall.EINVAL
}

func (s *stringIO) GoString() string { return s.name }

func (s *stringIO) Name() string { return s.name }

// Return the unread buffer.
func (s *stringIO) String() string {
    if s.isClosed() {
        return "<nil>"
    }
    return string(s.buf[s.pos:s.last])
}

// Return stored buffer until the last written position.
func (s *stringIO) GetValueString() string {
    if s.isClosed() {
        return "<nil>"
    }
    return string(s.buf[0:s.last])
}

// Return stored buffer as a byte array.
func (s *stringIO) GetValueBytes() []byte {
    if s.isClosed() {
        return s.buf[0:0]
    }
    return s.buf[0:s.last]
}

// Call Close will release the buffer/memory.
func (s *stringIO) Close() (err error) {
    s.Truncate(0)
    s.isclosed = true
    s.name = "StringIO <closed>"
    return
}

func (s *stringIO) Truncate(n int) {
    if s.isClosed() != true {
        if n == 0 {
            s.pos = 0
            s.last = 0
        }
        s.last = s.pos + n
        s.buf = s.buf[0:s.last]
    }
}

func (s *stringIO) Seek(offset int64, whence int) (ret int64, err error) {
    if s.isClosed() {
        return 0, OSError
    }
    pos, length := int64(s.pos), int64(len(s.buf))
    int64_O := int64(0)
    switch whence {
    case 0:
        ret = offset
    case 1:
        ret = offset + pos
    case 2:
        ret = offset + length
    default:
        return 0, syscall.EINVAL
    }
    if ret < int64_O {
        ret = int64_O
    }
    // stringIO currently does not support Seek beyond the
    // buf end, whereas posix does allow seek outside of
    // the file size, which will end up with a file hole.
    // However, StringIO does allow a byte hold within its
    // buffer size.
    if ret > length {
        ret = length
    }
    // Unfortunately, this will have to be a downcast.
    s.pos = int(ret)
    return
}

func (s *stringIO) Read(b []byte) (n int, err error) {
    if s.isClosed() {
        return 0, OSError
    }
    if s.pos >= len(s.buf) {
        return 0, io.EOF
    }
    return s.readBytes(b)
}

func (s *stringIO) ReadAt(b []byte, offset int64) (n int, err error) {
    if s.isClosed() {
        return 0, OSError
    }
    s.setPos(offset)
    return s.readBytes(b)
}

// stringIO Write will always be success until memory is used up
// or system limit is reached.
func (s *stringIO) Write(b []byte) (n int, err error) {
    if s.isClosed() {
        return 0, OSError
    }
    return s.writeBytes(b)
}

func (s *stringIO) WriteAt(b []byte, offset int64) (n int, err error) {
    if s.isClosed() {
        return 0, OSError
    }
    s.setPos(offset)
    return s.writeBytes(b)
}

func (s *stringIO) WriteString(str string) (ret int, err error) {
    b := syscall.StringByteSlice(str)
    return s.Write(b[0 : len(b)-1])
}


// private methods
func (s *stringIO) readBytes(b []byte) (n int, err error) {
    if s.pos > s.last {
        return 0, nil
    }
    n = len(b)
    // Require more than what we have only get what we have.
    // In other words, empty bytes will not be sent out.
    if s.pos+n > s.last {
        n = s.last - s.pos
    }
    copy(b, s.buf[s.pos:s.pos+n])
    s.pos += n
    return
}

func (s *stringIO) writeBytes(b []byte) (n int, err error) {
    n = len(b)
    if n > s.length() {
        s.resize(n)
    }
    copy(s.buf[s.pos:s.pos+n], b)
    s.pos += n
    if s.pos > s.last {
        s.last = s.pos
    }
    return
}

func (s *stringIO) setPos(offset int64) {
    pos, int64_O, length := int64(s.pos), int64(0), int64(len(s.buf))
    pos = offset
    if offset < int64_O {
        pos = int64_O
    }
    if offset > length {
        pos = length
    }
    s.pos = int(pos)
}

func (s *stringIO) length() int { return len(s.buf) - s.pos }

func (s *stringIO) isClosed() bool { return s.isclosed == true }

// Stolen from bytes.Buffer (Use the same algorithm)
func (s *stringIO) resize(n int) {
    if len(s.buf)+n > cap(s.buf) {
        buf := make([]byte, 2*cap(s.buf)+n)
        copy(buf, s.buf[0:])
        s.buf = buf
    }
}
