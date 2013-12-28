// codable stream or pipe
package pipestream

import (
	"bufio"
	"io"
	"os/exec"
)

// Stream data
type Stream struct {
	name   string   // command name
	arg    []string // command arguments
	cmd    *exec.Cmd
	reader *bufio.Reader
	writer *bufio.Writer
	pw     io.WriteCloser
	pr     io.ReadCloser
}

// Constructor factory
func New(cmd string, arg ...string) *Stream {
	return &Stream{name: cmd, arg: arg}
}

func (this *Stream) Open() error {
	this.cmd = exec.Command(this.name, this.arg...)

	var err error
	this.pr, err = this.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	this.pw, err = this.cmd.StdinPipe()
	if err != nil {
		return err
	}

	// startup
	if err := this.cmd.Start(); err != nil {
		return err
	}

	// prepare the reader/writer
	this.reader = bufio.NewReader(this.pr)
	this.writer = bufio.NewWriter(this.pw)
	return nil
}

// get reader to read from the pipe output
func (this *Stream) Reader() *bufio.Reader {
	return this.reader
}

// get writer to write to the pipe input
func (this *Stream) Writer() *bufio.Writer {
	return this.writer
}

// close the stream
func (this *Stream) Close() error {
	// close my writer stream so that client can get EOF
	this.pw.Close()

	// wait for client to exit
	if err := this.cmd.Wait(); err != nil {
		return err
	}

	return nil
}
