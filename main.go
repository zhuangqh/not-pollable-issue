package main

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"syscall"
	"time"
)

func main() {
	fmt.Printf("%v start test\n", time.Now().Format(time.TimeOnly))
	go manyTun()
	for {
		f()
	}
}

func manyTun() {
	for {
		ff, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "open dev %v", err)
			os.Exit(1)
		}

		if err := ff.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "close dev")
			os.Exit(1)
		}
	}
}

// GenerateID generates a random unique id.
func GenerateID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func f() {
	fp := fmt.Sprintf("/tmp/myfifo/%s", GenerateID())
	defer func() {
		os.Remove(fp)
	}()

	if err := syscall.Mkfifo(fp, uint32(0600&os.ModePerm)); err != nil && !os.IsExist(err) {
		fmt.Fprintf(os.Stderr, "create fifo read: %v", err)
		os.Exit(1)
	}

	mode := os.O_RDONLY | syscall.O_NONBLOCK
	fReq, err := os.OpenFile(fp, mode, 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open fifo read: %v", err)
		os.Exit(1)
	}

	wReq, err := os.OpenFile(fp, os.O_WRONLY|syscall.O_NONBLOCK, 0700)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open fifo write: %v", err)
		os.Exit(1)
	}

	_, err = wReq.Write([]byte("some message"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "write: %v", err)
		os.Exit(1)
	}

	buf := make([]byte, 4)
	_, err = fReq.Read(buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v fifo read: %v\n", time.Now().Format(time.TimeOnly), err)
		os.Exit(1)
	}

	//fmt.Printf("got %v\n", string(buf))

	fReq.Close()
	wReq.Close()
}
