package ccdrh

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	exitOK    = 0
	exitError = 1
)

func existingDirectory(r io.Reader) io.Reader {
	var buffer bytes.Buffer
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		dir := strings.Trim(scanner.Text()[1:], "'")
		fineinfo, err := os.Stat(dir)
		if err != nil {
			continue
		}

		if fineinfo.IsDir() {
			buffer.WriteString(scanner.Text() + "\n")
		}
	}

	return &buffer
}

func backupCachefile(cacheFile string, backupFile string) error {
	src, err := os.Open(cacheFile)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(backupFile)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return nil
}

func updateCachefile(cacheFile string) error {
	backupFile := cacheFile + ".bkup"
	err := backupCachefile(cacheFile, backupFile)
	if err != nil {
		return err
	}

	ifile, err := os.Open(backupFile)
	if err != nil {
		return err
	}
	defer ifile.Close()

	// Overwrite cacheFile
	ofile, err := os.Create(cacheFile)
	if err != nil {
		return err
	}
	defer ofile.Close()

	validDirs := existingDirectory(ifile)
	_, err = io.Copy(ofile, validDirs)
	if err != nil {
		return err
	}

	return nil
}

func Run() int {
	cacheFile := filepath.Join(os.Getenv("HOME"), ".cache", "shell", "chpwd-recent-dirs")

	err := updateCachefile(cacheFile)
	if err != nil {
		return exitError
	}
	return exitOK
}
