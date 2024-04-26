package main

import (
	"errors"
	"os"
	"os/exec"
	"syscall"

	"github.com/NodirBobiev/judge/internals/errorutil"
)

func main() {
	fs_bundle := os.Args[1]

	errors.Join()

	errorutil.Must(syscall.Sethostname([]byte("python-pc")), "set hostname")
	errorutil.Must(syscall.Chroot(fs_bundle), "change root")
	errorutil.Must(syscall.Chdir("/"), "change directory")
	errorutil.Must(syscall.Mount("/proc", "proc", "proc", 0, ""), "mount proc")
	defer syscall.Unmount("/proc", 0)

	cmd := exec.Command("python", "code.py")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	errorutil.Must(cmd.Run(), "run command")

}
