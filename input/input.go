package input

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
)

func readStdin() ([]byte, error) {
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}
	if (fileInfo.Mode() & os.ModeNamedPipe) != os.ModeNamedPipe {
		return nil, errors.New("cant read named pipe")
	}
	var b []byte
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		b = append(b, s.Bytes()...)
	}
	return b, nil
}

func readInput() ([]byte, error) {
	scan := bufio.NewScanner(os.Stdin)
	var b []byte
	fmt.Print("Enter characters:", "\n")
	for scan.Scan() {
		line := scan.Bytes()
		if bytes.Equal(line, []byte{0x0a}) {
			break
		}
		b = append(b, line...)
	}

	return b, nil
}

var ReadSeqList = []func() ([]byte, error){
	readStdin, readInput,
}

func ReadSeq() []byte {
	for _, readFunc := range ReadSeqList {
		b, err := readFunc()
		if err != nil {
			continue
		}
		if len(b) == 0 {
			continue
		}
		return b
	}
	return nil
}
