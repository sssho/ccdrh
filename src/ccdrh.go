package ccdrh

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
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

func exists(cdrtext string) bool {
	dir := strings.Trim(cdrtext[1:], "'")
	fineinfo, err := os.Stat(dir)
	if err != nil {
		return false
	}

	if !fineinfo.IsDir() {
		return false
	}
	return true
}

func updateCachefile(cacheFile string) error {
	ifile, err := os.Open(cacheFile)
	if err != nil {
		return err
	}
	defer ifile.Close()

	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()

		scanner := bufio.NewScanner(ifile)

		for scanner.Scan() {
			if exists(scanner.Text()) {
				io.WriteString(pw, scanner.Text()+"\n")
			}
		}
	}()

	tmpfile, err := ioutil.TempFile("", "ccdrh")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name())

	if _, err := io.Copy(tmpfile, pr); err != nil {
		return err
	}
	ifile.Close()
	if err := tmpfile.Close(); err != nil {
		return err
	}

	// Overwrite cacheFile
	// Use io.Copy to avoid cross-device link error
	if err := renameFile(tmpfile.Name(), cacheFile); err != nil {
		return err
	}
	return nil
}

func renameFile(in string, out string) error {
	ifile, err := os.Open(in)
	if err != nil {
		return err
	}
	defer ifile.Close()

	ofile, err := os.Create(out)
	if err != nil {
		return err
	}
	defer ofile.Close()

	if _, err := io.Copy(ofile, ifile); err != nil {
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
