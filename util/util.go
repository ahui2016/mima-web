package util

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
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
	return filepath.Glob(pattern)
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
