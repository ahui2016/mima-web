package db

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/ahui2016/mima-web/mima"
	"github.com/ahui2016/mima-web/tarball"
	"github.com/ahui2016/mima-web/util"
)

type DB struct {
	allItems     []*Mima
	homeCache    []*Mima
	recycleCache []*Mima
	userKey      *SecretKey
	key          *SecretKey
	Directory    string
	FullPath     string
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
	userKey := sha256.Sum256([]byte(password))
	db.userKey = &userKey
	m := newFirstMima()
	sealed64, err := db.encryptFirst(m)
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
	if err := db.refillItems(); err != nil {
		return err
	}
	fragFiles, err := db.getFragments()
	if err != nil {
		return err
	}
	if len(fragFiles) == 0 {
		return nil
	}
	filesToBackup := append(fragFiles, db.FullPath)
	if err := db.backupToTar(filesToBackup); err != nil {
		return err
	}
	if err := db.removeOldTarball(); err != nil {
		return err
	}
	if err := db.updateFromFragments(fragFiles); err != nil {
		return err
	}
	if err := db.rewriteFullPath(); err != nil {
		return err
	}
	return util.DeleteFiles(fragFiles)
}

func (db *DB) removeOldTarball() error {
	tarballs, err := db.getTarPaths()
	if err != nil {
		return err
	}
	if len(tarballs) > maxTarballs {
		if err := os.Remove(tarballs[0]); err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) backupToTar(files []string) error {
	tarFilePath := filepath.Join(db.Directory, util.TimestampFilename(tarballExt))
	return tarball.Create(tarFilePath, files)
}

func (db *DB) rewriteFullPath() error {
	dbFile, err := os.Create(db.FullPath)
	if err != nil {
		return err
	}
	defer dbFile.Close()

	var sealed64 string
	dbWriter := bufio.NewWriter(dbFile)
	for i, m := range db.allItems {
		if i == 0 {
			sealed64, err = db.encryptFirst(m)
		} else {
			sealed64, err = db.encrypt(m)
		}
		if err != nil {
			return err
		}
		if err := util.BufWriteln(dbWriter, sealed64); err != nil {
			return err
		}
	}
	return dbWriter.Flush()
}

func (db *DB) updateFromFragments(fragFiles []string) error {
	for _, f := range fragFiles {
		frag, err := db.readAndDecrypt(f)
		if err != nil {
			return err
		}
		if frag.Operation == mima.Insert {
			db.allItems = append(db.allItems, frag)
			continue
		}
		i, item, err := db.GetByID(frag.ID)
		if err != nil {
			return err
		}
		switch frag.Operation {
		case mima.Insert: // already proccess above
		case mima.Update:
			if item.UpdateFromFrag(frag) {
				db.moveToEnd(i)
			}
		case mima.SoftDelete:
			item.Delete()
		case mima.UnDelete:
			item.UnDelete()
		case mima.DeleteForever:
			db.deleteByIndex(i)
		}
	}
	return nil
}

func (db *DB) moveToEnd(i int) {
	item := db.allItems[i]
	db.allItems = append(db.allItems[:i], db.allItems[i+1:]...)
	db.allItems = append(db.allItems, item)
}

func (db *DB) deleteByIndex(i int) {
	db.allItems = append(db.allItems[:i], db.allItems[i+1:]...)
}

// func (db *DB) deleteByID(id string) (*Mima, error) {
// 	i, m, err := db.GetByID(id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	db.allItems = append(db.allItems[:i], db.allItems[i+1:]...)
// 	return m, nil
// }

func (db *DB) readAndDecrypt(filePath string) (*Mima, error) {
	blob, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return db.decrypt(string(blob))
}

func (db *DB) getFragments() ([]string, error) {
	return util.GetPathsByExt(db.Directory, fragmentExt)
}

func (db *DB) getTarPaths() ([]string, error) {
	return util.GetPathsByExt(db.Directory, tarballExt)
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

// refillItems retrieves all items (except the first one) from db.FullPath.
func (db *DB) refillItems() error {
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

func (db *DB) Update(form *EditForm) error {
	if form.Title == "" {
		return errTitleEmpty
	}
	i, m, err := db.GetByID(form.ID)
	if err != nil {
		return err
	}
	changeIndex, writeFrag := m.UpdateFromForm(form)
	if changeIndex {
		db.moveToEnd(i)
	}
	if writeFrag {
		return db.encryptWriteFragment(m, mima.Update)
	}
	return nil
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

func (db *DB) encryptFirst(m *Mima) (string, error) {
	mimaJSON, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return seal64(mimaJSON, db.userKey)
}

func (db *DB) writeFragFile(sealed64 string) error {
	fragPath := filepath.Join(db.Directory, util.TimestampFilename(fragmentExt))
	return writeFile(fragPath, sealed64)
}

func (db *DB) updateHomeCache() (all []*Mima) {
	for i := db.Len() - 1; i > 0; i-- {
		m := db.allItems[i].HideSecrets()
		if m.IsDeleted() {
			continue
		}
		all = append(all, m)
	}
	return
}

func (db *DB) DeletedItems() (deleted []*Mima) {
	for i := 1; i < db.Len(); i++ {
		m := db.allItems[i].HideSecrets()
		if !m.IsDeleted() {
			continue
		}
		deleted = append(deleted, m)
	}
	sort.Slice(deleted, func(i, j int) bool {
		return deleted[i].DeletedAt > deleted[j].DeletedAt
	})
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

func (db *DB) DeleteHistory(id, datetime string) error {
	_, m, err := db.GetByID(id)
	if err != nil {
		return err
	}
	if err := m.DeleteHistory(datetime); err != nil {
		return err
	}
	return db.encryptWriteFragment(m, mima.Update)
}

func (db *DB) RecycleByID(id string) error {
	_, m, err := db.GetByID(id)
	if err != nil {
		return err
	}
	m.Delete()
	return db.encryptWriteFragment(m, mima.SoftDelete)
}

func (db *DB) RecoverByID(id string) error {
	_, m, err := db.GetByID(id)
	if err != nil {
		return err
	}
	if !m.IsDeleted() {
		return errors.New(id + " not found in recycle bin")
	}
	m.UnDelete()
	return db.encryptWriteFragment(m, mima.UnDelete)
}

func (db *DB) DeleteForever(id string) error {
	i, m, err := db.GetByID(id)
	if err != nil {
		return err
	}
	db.deleteByIndex(i)
	return db.encryptWriteFragment(m, mima.DeleteForever)
}

func (db *DB) GetByAlias(alias string) (mimas []*Mima) {
	if alias == "" {
		return
	}
	for i := 1; i < db.Len(); i++ {
		m := db.allItems[i]
		if mima.IsDeleted() {
			continue
		}
		if mima.Alias == alias {
			mimas = append(mimas, mima)
		}
	}
	return
}
