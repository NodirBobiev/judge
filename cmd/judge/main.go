package main

import (
	"io"
	"os"
	"os/exec"
	"path"
	"syscall"

	"github.com/NodirBobiev/judge/internals/errorutil"
)

func main() {
	run()
}

func copyFile(src, dest string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

func prepare_fs_bundle() string {
	fs_path := "/fs_bundles/python3.12"
	src_path := "/code_examples/helloworld.py"
	dest_path := path.Join(fs_path, "code.py")
	errorutil.Must(copyFile(src_path, dest_path), "copy source code")

	return fs_path
}

func run() {
	fs_bundle := prepare_fs_bundle()

	cmd := exec.Command("./judge-py", fs_bundle)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER,
		Unshareflags: syscall.CLONE_NEWNS,
		Credential:   &syscall.Credential{Uid: 0, Gid: 0},
		UidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: os.Getuid(), Size: 1},
		},
		GidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: os.Getgid(), Size: 1},
		},
	}

	errorutil.Must(cmd.Run(), "run command")
}
