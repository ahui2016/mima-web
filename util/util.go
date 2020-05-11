package util

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"math/big"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

const (
	ISO8601            = "2006-01-02T15:04:05.999Z"
	databaseFolderName = "MimaWebDB"
)

func NewID() string {
	var max int64 = 100_000_000
	n, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		panic(err)
	}
	timestamp := time.Now().Unix()
	idInt64 := timestamp*max + n.Int64()
	return strconv.FormatInt(idInt64, 36)
}

func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func Base64Decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func NewFileScanner(fullPath string) (*bufio.Scanner, *os.File, error) {
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, nil, err
	}
	return bufio.NewScanner(file), file, nil
}

func GetPathsByExt(dir, ext string) ([]string, error) {
	pattern := filepath.Join(dir, "*"+ext)
	filePaths, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	sort.Strings(filePaths)
	return filePaths, nil
}

func DatabaseDefaultDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(home, databaseFolderName)
}

func TimestampFilename(ext string) string {
	name := strconv.FormatInt(time.Now().UnixNano(), 10)
	return name + ext
}

func BufWriteln(w *bufio.Writer, box64 string) (err error) {
	_, err = w.WriteString(box64 + "\n")
	return
}

func DeleteFiles(filePaths []string) error {
	for _, f := range filePaths {
		if err := os.Remove(f); err != nil {
			return err
		}
	}
	return nil
}

func TimeNow() string {
	return time.Now().Format(ISO8601)
}
