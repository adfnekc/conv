package input

import (
	"bufio"
	"errors"

	"os"
	"reflect"
	"runtime"
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

var ReadSeqList = []func() ([]byte, error){
	readStdin, readInput,
}

func ReadSeq() []byte {
	for _, readFunc := range ReadSeqList {
		b, err := readFunc()
		if err != nil {
			// log.Printf("readFunc <%v> err :%s", getFunctionName(readFunc), err)
			continue
		}
		if len(b) == 0 {
			continue
		}
		return b
	}
	return nil
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
