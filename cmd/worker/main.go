package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"syscall"

	"github.com/NodirBobiev/judge/internals/errorutil"
	"github.com/NodirBobiev/judge/internals/kafka"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		log.Println("Recieved termination signal. Graceful shutdown")
		cancel()
	}()

	kafkaClient, err := kafka.NewKafkaClient([]string{"localhost:29092"})
	if err != nil {
		fmt.Printf("Failed to create Kafka client: %v\n", err)
		return
	}
	defer kafkaClient.Close()

	topic := "file-uploads"

	for {
		select {
		case <-ctx.Done():
			return
		default:
			message, err := kafkaClient.ConsumeMessageWithContext(ctx, topic)
			if err != nil {
				fmt.Printf("Error consuming message: %v\n", err)
				continue
			}
			if message != nil {
				// err := os.WriteFile(message.Filename, message.Content, 0644)
				// if err != nil {
				// 	fmt.Printf("Failed to write file: %v\n", err)
				// 	continue
				// }
				fmt.Printf("File '%s' saved successfully\n", message.Filename)
				fmt.Println(string(message.Content))
			}
		}
	}
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

func prepare_fs_bundle(src_path string) string {
	fs_path := "fs_bundles/python3.12"

	errorutil.Must(copyFile(src_path, path.Join(fs_path, "code.py")), "copy source code")
	errorutil.Must(copyFile("tests/test.in", path.Join(fs_path, "test.in")), "copy test.in")
	errorutil.Must(copyFile("tests/test.out", path.Join(fs_path, "test.out")), "copy test.out")

	return fs_path
}

func run(src_path string) {
	fs_bundle := prepare_fs_bundle(src_path)

	cmd := exec.Command("./worker-py", fs_bundle)
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
