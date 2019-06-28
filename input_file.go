package tg

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
)

// InputFile represents the file that should be uploaded to the telegram.
type InputFile struct {
	// Filename
	Name string

	// Body of file
	Body io.Reader
}

// NewInputFile creates the InputFile from provided name and body reader.
func NewInputFile(name string, body io.Reader) InputFile {
	return InputFile{
		Name: name,
		Body: body,
	}
}

// NewInputFileBytes creates input file from provided name.
//
// Example:
//   file := NewInputFileBytes("test.txt", []byte("test, test, test..."))
func NewInputFileBytes(name string, body []byte) InputFile {
	return NewInputFile(
		name,
		bytes.NewReader(body),
	)
}

// NewInputFileLocal creates the InputFile from provided local file.
// This method just open file by provided path.
// So, you should close it AFTER send.
//
// Example:
//
//   file, err := NewInputFileLocal("test.png")
//   if err != nil {
//       return err
//   }
//   defer file.Close()
//
func NewInputFileLocal(path string) (InputFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return InputFile{}, err
	}

	return NewInputFile(
		filepath.Base(file.Name()),
		file,
	), nil
}

// NewInputFileLocalBuffer creates the InputFile from provided local file path.
// This function copy the file to the buffer in memory, so you do not need to close it.
func NewInputFileLocalBuffer(path string) (InputFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return InputFile{}, err
	}

	name := filepath.Base(file.Name())

	body := &bytes.Buffer{}

	if _, err := io.Copy(body, file); err != nil {
		return InputFile{}, err
	}

	return NewInputFile(name, body), nil
}

// Close call close method on body if it implements io.ReadCloser,
// else - do nothing.
func (file InputFile) Close() error {
	if body, ok := file.Body.(io.ReadCloser); ok {
		return body.Close()
	}

	return nil
}

// AddFileToRequest implementation of Media interface.
func (file InputFile) AddFileToRequest(k string, r *Request) {
	r.AddFile(k, file)
}
