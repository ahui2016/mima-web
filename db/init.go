package db

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"

	"github.com/ahui2016/mima-web/mima"
	"github.com/ahui2016/mima-web/util"
	"golang.org/x/crypto/nacl/secretbox"
)

const (
	KeySize     = 32
	NonceSize   = 24
	dbName      = "mimaweb.db"
	fragmentExt = ".frag"
)

var (
	errTitleEmpty = errors.New("mima.Title is empty")
)

type (
	Nonce     = [NonceSize]byte
	SecretKey = [KeySize]byte
	Mima      = mima.Mima
	Operation = mima.Operation
)

// NewFirstMima create the very first mima, to save a random serect key.
func newFirstMima() *Mima {
	m := new(Mima)
	key := randomKey()
	m.Username = RandomString64()
	m.Password = util.Base64Encode(key[:])
	m.CreatedAt = time.Now().UnixNano()
	return m
}

func newNonce() (nonce Nonce, err error) {
	_, err = rand.Read(nonce[:])
	return
}

func seal64(data []byte, key *SecretKey) (string, error) {
	nonce, err := newNonce()
	if err != nil {
		return "", err
	}
	sealed := secretbox.Seal(nonce[:], data, &nonce, key)
	return util.Base64Encode(sealed), nil
}

func decrypt64(sealed64 string, key *SecretKey) (*Mima, error) {
	sealed, err := util.Base64Decode(sealed64)
	if err != nil {
		return nil, err
	}
	if len(sealed) < NonceSize {
		return nil, errors.New("it's not a secretbox")
	}
	var nonce Nonce
	copy(nonce[:], sealed[:NonceSize])
	mimaJSON, ok := secretbox.Open(nil, sealed[NonceSize:], &nonce, key)
	if !ok {
		return nil, errors.New("db.decrypt: secretbox open fail")
	}
	m := new(Mima)
	if err := json.Unmarshal(mimaJSON, m); err != nil {
		return nil, err
	}
	return m, nil
}

func base64DecodeToSecretKey(s string) (*SecretKey, error) {
	b, err := util.Base64Decode(s)
	if err != nil {
		return nil, err
	}
	key := bytesToKey(b)
	return &key, nil
}

func bytesToKey(b []byte) (key SecretKey) {
	copy(key[:], b)
	return
}

func writeFile(fullPath, sealed64 string) error {
	return ioutil.WriteFile(fullPath, []byte(sealed64), 0644)
}

func RandomString64() string {
	size := 64
	someBytes := make([]byte, size)
	if _, err := rand.Read(someBytes); err != nil {
		panic(err)
	}
	return util.Base64Encode(someBytes)
}

func randomKey() SecretKey {
	return sha256.Sum256([]byte(RandomString64()))
}
