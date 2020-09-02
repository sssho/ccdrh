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

func updateCachefile(cacheFile string) error {
	ifile, err := os.Open(cacheFile)
	if err != nil {
		return err
	}
	defer ifile.Close()
	validDirs := existingDirectory(ifile)
	ifile.Close()

	// Overwrite cacheFile
	ofile, err := os.Create(cacheFile)
	if err != nil {
		return err
	}
	defer ofile.Close()

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
