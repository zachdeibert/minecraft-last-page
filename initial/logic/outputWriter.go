package logic

import (
	"os"
	"strings"
)

type outputWriter struct {
	file *os.File
}

func createOutputWriter(filename string) (*outputWriter, error) {
	f, e := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR|os.O_SYNC, 0644)
	if e != nil {
		return nil, e
	}
	if _, e = f.Seek(0, 2); e != nil {
		f.Close()
		return nil, e
	}
	return &outputWriter{
		file: f,
	}, nil
}

func (ow outputWriter) Write(strs []string) error {
	if strs != nil && len(strs) > 0 {
		str := strings.Join(strs, "\n") + "\n"
		_, e := ow.file.WriteString(str)
		return e
	}
	return nil
}

func (ow outputWriter) Close() error {
	return ow.file.Close()
}
