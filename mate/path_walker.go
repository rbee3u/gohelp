package mate

import (
	"bytes"
	"strconv"
)

func NewPathWalker() *PathWalker {
	return new(PathWalker)
}

type PathWalker []path

func (pw *PathWalker) Enter(p path) {
	*pw = append(*pw, p)
}

func (pw *PathWalker) Exit() {
	*pw = (*pw)[:len(*pw)-1]
}

func (pw PathWalker) String() string {
	return pw.Prefix(len(pw))
}

func (pw PathWalker) Prefix(n int) string {
	buf := new(bytes.Buffer)

	for i, p := range pw[:n] {
		if i > 0 {
			buf.WriteByte('_')
		}

		p.writeToBuffer(buf)
	}

	return buf.String()
}

type path interface {
	writeToBuffer(buf *bytes.Buffer)
}

type IntegerPath int

func (p IntegerPath) writeToBuffer(buf *bytes.Buffer) {
	buf.WriteString(strconv.Itoa(int(p)))
}

type StringPath string

func (p StringPath) writeToBuffer(buf *bytes.Buffer) {
	buf.WriteString(string(p))
}
