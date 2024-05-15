package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/NodirBobiev/judge/internals/assert"
	"github.com/NodirBobiev/judge/internals/errorutil"
	"golang.org/x/sync/errgroup"
)

func matchByWords(actualReader, expectedReader io.Reader) error {
	actualScanner := bufio.NewScanner(actualReader)
	actualScanner.Split(bufio.ScanWords)

	expectedScanner := bufio.NewScanner(expectedReader)
	expectedScanner.Split(bufio.ScanWords)
	wordCounter := 0
	for {
		wordCounter++
		actualScan := actualScanner.Scan()
		expectedScan := expectedScanner.Scan()
		if !actualScan && !expectedScan {
			return nil
		}

		actualWord := actualScanner.Text()
		expectedWord := expectedScanner.Text()

		if err := assert.Equalf(expectedWord, actualWord, "not matching word %d", wordCounter); err != nil {
			return err
		}

	}
}

func checkEmpty(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	scanner.Scan()
	return assert.Equalf("", scanner.Text(), "not empty")
}

func main() {
	fs_bundle := os.Args[1]

	errorutil.Must(syscall.Sethostname([]byte("python-pc")), "set hostname")
	errorutil.Must(syscall.Chroot(fs_bundle), "change root")
	errorutil.Must(syscall.Chdir("/"), "change directory")
	// errorutil.Must(syscall.Mount("/proc", "proc", "proc", 0, ""), "mount proc")
	defer syscall.Unmount("/proc", 0)

	cmd := exec.Command("python", "code.py")

	inputFile, err := os.Open("test.in")
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	outputFile, err := os.Open("test.out")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	cmd.Stdin = inputFile

	stdoutPipe, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()

	errorutil.Must(cmd.Start(), "start command")

	var g errgroup.Group
	g.Go(func() error {
		return assert.Assert(matchByWords(stdoutPipe, outputFile), "WA")
	})
	g.Go(func() error {
		return assert.Assert(checkEmpty(stderrPipe), "RE")
	})

	if err := assert.Assert(cmd.Wait(), "RE"); err != nil {
		fmt.Println(err)
		return
	}

	if err := g.Wait(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("AC")
}
