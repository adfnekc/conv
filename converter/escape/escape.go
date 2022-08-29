package escape

import (
	"bytes"
	"conv/converter"
	"errors"
	"fmt"
	"strconv"
)

func init() {
	converter.AddConverter("c_escape", "C Style escape", new(CEscape))
	converter.AddConverter("mulit_line_escape", "C Style escape with mulit line", new(MuliLineEscape))
}

func Escape(src []byte) ([]byte, error) {
	var dst = make([]byte, len(src)/4)
	round := 0
	for i := 0; i < len(src); {
		step := 4
		if string(src[i]) != "\\" {
			return nil, errors.New("malformed escape string")
		}
		if string(src[i+1]) == "0" {
			step = 2
			dst[round] = 0
		} else if string(src[i+1]) == "x" {
			if string(src[i+3]) == "\\" {
				dst[round] = src[i+2]
				step = 3
			} else {
				i16, err := strconv.ParseInt(string(src[i+2:i+4]), 16, 0)
				if err != nil {
					return nil, err
				}
				dst[round] = byte(i16)
			}
		} else {
			return nil, errors.New("malformed escape string")
		}
		i += step
		round++
	}
	return dst, nil
}
func Unescape(src []byte) []byte {
	var dst = make([]byte, len(src)*4)
	prefix := []byte("\\x")
	for i, v := range src {
		dst[i*4] = prefix[0]
		dst[i*4+1] = prefix[1]
		dst[i*4+2] = []byte(fmt.Sprintf("%x", v))[0]
		dst[i*4+3] = []byte(fmt.Sprintf("%x", v))[1]
	}
	return dst
}

type CEscape struct {
}

func (c CEscape) From(src []byte) ([]byte, error) {
	return Escape(src)
}

func (c CEscape) To(src []byte) []byte {
	return Unescape(src)
}

type GOEscape struct {
}

func (c GOEscape) From([]byte) ([]byte, error) {
	return nil, nil
}

func (c GOEscape) To(src []byte) []byte {
	var dst = make([]byte, len(src)*4+1)
	prefix := []byte("\\x")
	suffix := []byte(",")
	dst[0] = []byte("{")[0]
	for i, v := range src {
		dst[i*4+1] = prefix[0]
		dst[i*4+2] = prefix[1]
		dst[i*4+3] = []byte(fmt.Sprintf("%x", v))[0]
		dst[i*4+4] = suffix[3]
	}
	return append(dst, []byte("}")...)
}

type MuliLineEscape struct {
}

func (c MuliLineEscape) From(src []byte) ([]byte, error) {
	src = bytes.TrimSpace(src)
	if !(bytes.HasSuffix(src, []byte(`"`)) && bytes.HasPrefix(src, []byte(`"`))) {
		return nil, errors.New("malformed MuliLineEscape string")
	}
	var dst []byte
	lines := bytes.Split(src, []byte{0x0a})
	for _, line := range lines {
		line = bytes.TrimSpace(line)
		line = bytes.Trim(line, `"`)
		line = bytes.TrimSpace(line)
		line = bytes.TrimRight(line, `\`)
		line = bytes.TrimSpace(line)
		line = bytes.Trim(line, `"`)

		dst = append(dst, line...)
	}
	return Escape(dst)
}

func (c MuliLineEscape) To(src []byte) []byte {
	return nil
}
