package db

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ahui2016/mima-web/mima"
	"github.com/ahui2016/mima-web/util"
)

type DB struct {
	allItems  []*Mima
	userKey   *SecretKey
	key       *SecretKey
	Directory string
	FullPath  string
}

func NewDB(directory string) *DB {
	return &DB{
		Directory: directory,
		FullPath:  filepath.Join(directory, dbName),
	}
}

func (db *DB) IsEmpty() bool {
	if len(db.allItems) > 0 {
		return false
	}
	return true
}

func (db *DB) Create(password string) {
	if db.fullPathExist() {
		log.Fatalf("%s already exists", db.FullPath)
	}
	m := newFirstMima()
	mimaJSON, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	userKey := sha256.Sum256([]byte(password))
	sealed64, err := seal64(mimaJSON, &userKey)
	if err != nil {
		panic(err)
	}
	if err := writeFile(db.FullPath, sealed64); err != nil {
		panic(err)
	}
}

func (db *DB) Init(password string) error {
	if db.IsReady() {
		return errors.New("DB.allItems is not empty")
	}
	if db.fullPathNotExist() {
		return fmt.Errorf("%s do not exists", db.FullPath)
	}
	if err := db.setKeys(password); err != nil {
		return err
	}
	return nil
}

func (db *DB) Reset() {
	db.allItems = nil
	db.userKey = nil
	db.key = nil
}

func (db *DB) CheckPassword(password string) bool {
	userKey := sha256.Sum256([]byte(password))
	return *db.userKey == userKey
}

func (db *DB) IsReady() bool {
	if len(db.allItems) > 0 {
		return true
	}
	return false
}

func (db *DB) setKeys(password string) error {
	userKey := sha256.Sum256([]byte(password))
	db.userKey = &userKey
	if err := db.setKeyAndFirst(); err != nil {
		return err
	}
	return nil
}

// setKeyAndFirst sets db.key and append ths first item into db.allItems.
func (db *DB) setKeyAndFirst() error {
	info, err := os.Lstat(db.FullPath)
	if err != nil {
		return err
	}
	if info.Size() < NonceSize {
		return fmt.Errorf("%s is empty", db.FullPath)
	}
	m, err := db.getFirstItem()
	if err != nil {
		return fmt.Errorf("DB.setKeyAndFirst: %w", err)
	}
	key, err := base64DecodeToSecretKey(m.Password)
	if err != nil {
		return err
	}
	db.key = key
	db.allItems = append(db.allItems, m)
	return nil
}

func (db *DB) getFirstItem() (m *Mima, err error) {
	scanner, file, err := util.NewFileScanner(db.FullPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	for scanner.Scan() {
		sealed64 := scanner.Text()
		m, err = db.decryptFirst(sealed64)
		break
	}
	return
}

// retrieveItems retrieves all items (except the first one) from db.FullPath.
func (db *DB) retrieveItems() error {
	scanner, file, err := util.NewFileScanner(db.FullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	n := 0
	for scanner.Scan() {
		n++
		sealed64 := scanner.Text()
		if n == 1 {
			continue
		}
		m, err := db.decrypt(sealed64)
		if err != nil {
			return fmt.Errorf("decrypting line %d : %w", n, err)
		}
		db.allItems = append(db.allItems, m)
	}
	return scanner.Err()
}

func (db *DB) decryptFirst(sealed64 string) (*Mima, error) {
	return decrypt64(sealed64, db.userKey)
}

func (db *DB) decrypt(sealed64 string) (*Mima, error) {
	return decrypt64(sealed64, db.key)
}

func (db *DB) fullPathNotExist() bool {
	_, err := os.Stat(db.FullPath)
	if os.IsNotExist(err) {
		return true
	}
	if err != nil {
		panic(err)
	}
	return false
}

func (db *DB) fullPathExist() bool {
	return !db.fullPathNotExist()
}

// Insert inserts a new Mima into the end of db.allItems.
// The db.allItems is sorted by Mima.UpdatedAt.
func (db *DB) Insert(m *Mima) error {
	if m.Title == "" {
		return errTitleEmpty
	}
	db.allItems = append(db.allItems, m)
	return db.encryptWriteFragment(m, mima.Insert)
}

func (db *DB) encryptWriteFragment(m *Mima, op Operation) error {
	m.Operation = op
	sealed64, err := db.encrypt(m)
	if err != nil {
		return nil
	}
	return db.writeFragFile(sealed64)
}

func (db *DB) encrypt(m *Mima) (string, error) {
	mimaJSON, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return seal64(mimaJSON, db.key)
}

func (db *DB) writeFragFile(sealed64 string) error {
	fragPath := filepath.Join(db.Directory, util.TimestampFilename(fragmentExt))
	log.Print(fragPath)
	return writeFile(fragPath, sealed64)
}

// AllItems return all items except the very first item and soft-deleted items,
// and hides password and notes. Sorted decending by Mima.UpdatedAt.
// TODO: cache
func (db *DB) AllItems() (all []*Mima) {
	if db.Len()-1 == 0 {
		return
	}
	for i := db.Len() - 1; i > 0; i-- {
		m := db.allItems[i].HideSecrets()
		if m.IsDeleted() {
			continue
		}
		all = append(all, m)
	}
	return
}

func (db *DB) Len() int {
	return len(db.allItems)
}

func (db *DB) GetByID(id string) (i int, m *Mima, err error) {
	for i = 1; i < db.Len(); i++ {
		m = db.allItems[i]
		if m.ID == id {
			return
		}
	}
	err = fmt.Errorf("id: %s not found", id)
	return
}
