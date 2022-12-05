package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"golang.org/x/sys/unix"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Unable to run: %v", err)
	}
}

func run() error {
	const basePath = "/var/lib/containers/storage/overlay"
	backingFsBlockDev := path.Join(basePath, "backingFsBlockDev")

	if _, err := os.Stat(backingFsBlockDev); err == nil {
		log.Printf("The block device %s already exists, skipping recreation", backingFsBlockDev)
		return nil
	}

	log.Printf("Recreating %s", backingFsBlockDev)

	var stat unix.Stat_t
	if err := unix.Stat(basePath, &stat); err != nil {
		return fmt.Errorf("stat base path: %w", err)
	}

	backingFsBlockDevTmp := backingFsBlockDev + ".tmp"

	if err := unix.Mknod(backingFsBlockDevTmp, unix.S_IFBLK|0600, int(stat.Dev)); err != nil {
		return fmt.Errorf("failed to mknod %s: %w", backingFsBlockDevTmp, err)
	}
	if err := unix.Rename(backingFsBlockDevTmp, backingFsBlockDev); err != nil {
		return fmt.Errorf("failed to rename %s to %s: %w", backingFsBlockDevTmp, backingFsBlockDev, err)
	}

	return nil
}
