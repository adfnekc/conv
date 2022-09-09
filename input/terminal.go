package input

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

const (
	ioctlReadTermios  = unix.TCGETS
	ioctlWriteTermios = unix.TCSETS
	pollInterval      = 5 * time.Millisecond
	inputInterval     = 300 * time.Millisecond
	pollDelay         = 50
	consoleDevice     = "/dev/tty"
)

var (
	devPrefixes = [...]string{"/dev/pts/", "/dev/"}
	ttyin       = openTtyIn()
	input       = []rune{}
)

func setNonblock(fd int, nonblock bool) error {
	return syscall.SetNonblock(fd, nonblock)
}

func rawread(fd int, b []byte) (int, error) {
	n, err := syscall.Read(fd, b)
	return n, err
}

func openTtyIn() *os.File {
	in, err := os.OpenFile(consoleDevice, syscall.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		log.Panic("Failed to open " + consoleDevice)
	}
	return in
}

func getchar(fd int, nonblock bool) (int, bool) {
	b := make([]byte, 1)
	err := setNonblock(fd, nonblock)
	if err != nil {
		fmt.Println(err)
		return 0, false
	}
	_, err = rawread(fd, b)
	if err != nil {
		// fmt.Println(err)
		return 0, false
	}
	return int(b[0]), true
}

func setNoLineInput(fd int) (*term.State, error) {
	termios, err := unix.IoctlGetTermios(fd, ioctlReadTermios)
	if err != nil {
		return nil, err
	}

	oldState, err := term.GetState(fd)
	if err != nil {
		return nil, err
	}

	termios.Lflag &^= unix.ECHO | unix.ICANON
	termios.Cc[unix.VMIN] = 1
	termios.Cc[unix.VTIME] = 0
	if err := unix.IoctlSetTermios(fd, ioctlWriteTermios, termios); err != nil {
		return nil, err
	}

	return oldState, nil
}

func readch(ch chan byte) {
	fd := int(ttyin.Fd())
	origState, err := setNoLineInput(fd)
	if err != nil {
		panic(err)
	}
	defer term.Restore(fd, origState)

	for {
		c, ok := getchar(fd, false)
		if !ok {
			time.Sleep(pollInterval)
			continue
		}
		ch <- byte(c)

		retries := pollDelay
		for {
			c, ok := getchar(fd, true)
			if !ok {
				if retries > 0 {
					retries--
					time.Sleep(pollInterval)
					continue
				}
				break
			}
			ch <- byte(c)
		}
		time.Sleep(pollInterval)
	}
}

func readInput() ([]byte, error) {
	var b []byte
	chchan := make(chan byte)
	go readch(chchan)

	ch, ok := <-chchan
	if !ok {
		return b, nil
	}
	b = append(b, ch)

	timer := time.NewTimer(inputInterval)
	timer.Stop()
	for {
		select {
		case <-timer.C:
			return b, nil
		case ch, ok := <-chchan:
			if !ok {
				return b, nil
			}
			if ch == 0xa {
				timer.Reset(inputInterval)
			}
			b = append(b, ch)
		}
	}
}
