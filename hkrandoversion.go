package main

import (
	"bufio"
	"encoding/binary"
	"crypto/sha1"
	"os"
	"io"
	"regexp"
	"strings"
	"fmt"
)

type latin1Reader struct {
	src *bufio.Reader
}

func (r latin1Reader) ReadRune() (rune, int, error) {
	b, err := r.src.ReadByte()
	return rune(b), 1, err
}

func randoVersion(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	base, err := baseVersion(f)
	if err != nil {
		return "", err
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return base + " (???)", err
	}
	hash, err := hashVersion(f)
	if err != nil {
		return base + " (???)", err
	}
	return fmt.Sprintf("%s (%d)", base, hash), nil
}

func baseVersion(f io.ReadSeeker) (string, error) {
	wrap := func(err error) error {
		return fmt.Errorf("error determining base version: %w", err)
	}
	const unknown = "?.??"

	re := regexp.MustCompile(`\d\x00.\x00\d\x00\d(\x00[A-Z])+`)
	loc := re.FindReaderIndex(latin1Reader{bufio.NewReader(f)})
	
	if _, err := f.Seek(int64(loc[0]), io.SeekStart); err != nil {
		return unknown, wrap(err)
	}
	buf := make([]byte, loc[1]-loc[0])
	if _, err := io.ReadFull(f, buf); err != nil {
		return unknown, wrap(err)
	}
	return strings.ReplaceAll(string(buf), "\x00", ""), nil
}

func hashVersion(f io.Reader) (int, error) {
	wrap := func(err error) error {
		return fmt.Errorf("error determining hash version: %w", err)
	}

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return -1, wrap(err)
	}
	sha := h.Sum(make([]byte, 0, sha1.Size))
	ver := int32(0)
	for i := 0; i < len(sha)-1; i += 4 {
		ver = 17 * ver + 31 * int32(binary.LittleEndian.Uint32(sha[i:i+4]))
	}
	ver %= 997
	if ver < 0 {
		ver = -ver
	}
	return int(ver), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage:", os.Args[0], "<rando dll>")
		os.Exit(2)
	}
	v, err := randoVersion(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(v)
}