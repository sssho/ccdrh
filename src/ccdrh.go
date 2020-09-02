package ccdrh

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	exitOK    = 0
	exitError = 1
)

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

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

func validateInputFile(path string) error {
	ifile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer ifile.Close()

	buf := make([]byte, 2)

	_, err = ifile.Read(buf)
	if err != nil {
		return err
	}
	if string(buf) != "$'" {
		return &errorString{"Invalid cdr history file format"}
	}
	return nil
}

func DefaultCacheFilePath() string {
	return filepath.Join(os.Getenv("HOME"), ".cache", "shell", "chpwd-recent-dirs")
}

func Run(cacheFile string) int {
	if cacheFile == "" {
		cacheFile = DefaultCacheFilePath()
	} else {
		err := validateInputFile(cacheFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return exitError
		}
	}

	err := updateCachefile(cacheFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitError
	}
	return exitOK
}
