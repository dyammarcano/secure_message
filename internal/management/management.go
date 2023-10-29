package management

import (
	"crypto/rand"
	"encoding/json"
	"github.com/dyammarcano/secure_message/internal/config"
	"github.com/dyammarcano/secure_message/internal/database"
	"os"
	"sync"
	"time"
)

var mgmt *Management

type (
	Management struct {
		db    *database.Database
		keys  []Keys
		mutex sync.RWMutex
	}

	Keys struct {
		ID         uint      `gorm:"primarykey" yaml:"-" json:"-"`
		CreatedAt  time.Time `yaml:"-" json:"-"`
		UpdatedAt  time.Time `yaml:"-" json:"-"`
		PrivateKey []byte    `yaml:"private_key" json:"private_key"`
		Hits       int       `yaml:"hits" json:"hits"`
	}
)

func init() {
	m, err := newManagement()
	if err != nil {
		panic(err)
	}
	mgmt = m
}

func newManagement() (*Management, error) {
	db, err := database.NewDatabase(config.C.DatabaseFilePath)
	if err != nil {
		return nil, err
	}

	m := &Management{
		db:    db,
		keys:  make([]Keys, 1024),
		mutex: sync.RWMutex{},
	}

	if err = m.db.Migrate(&Keys{}); err != nil {
		return nil, err
	}

	if err := m.db.FindAll(&m.keys); err != nil {
		return nil, err
	}

	if len(m.keys) == 0 {
		if err := m.generateKeys(); err != nil {
			return nil, err
		}
	}

	return m, err
}

func (m *Management) generateKeys() error {
	m.keys = make([]Keys, 1024)
	for i := range m.keys {
		bytes := make([]byte, 44) // 352 bits key
		if _, err := rand.Read(bytes); err != nil {
			return err
		}
		m.keys[i] = Keys{PrivateKey: bytes}
	}

	if err := m.saveKeys(); err != nil {
		return err
	}

	return nil
}

func (m *Management) saveKeys() error {
	if err := m.db.BulkInsert(m.getKeys()); err != nil {
		return err
	}
	return nil
}

func (m *Management) setKeys(k []Keys) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.keys = k
}

func (m *Management) overrideData() error {
	// delete database file
	m.db.Close()
	if err := os.Remove(config.C.DatabaseFilePath); err != nil {
		return err
	}

	db, err := database.NewDatabase(config.C.DatabaseFilePath)
	if err != nil {
		return err
	}

	m.db = db
	if err = m.db.Migrate(&Keys{}); err != nil {
		return err
	}

	if err := m.db.BulkInsert(m.getKeys()); err != nil {
		return err
	}

	return nil
}

func (m *Management) getKeys() []Keys {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.keys
}

func (m *Management) hitKey(k Keys) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	k.Hits++
	if err := m.db.Update(k); err != nil {
		return err
	}

	m.updateSlice(k)
	return nil
}

func (m *Management) unHitKey(k Keys) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if k.Hits == 0 {
		return nil
	}

	k.Hits--
	if err := m.db.Update(k); err != nil {
		return err
	}

	m.updateSlice(k)
	return nil
}

func (m *Management) updateSlice(k Keys) {
	for i := range m.keys {
		if m.keys[i].ID == k.ID {
			m.keys[i].Hits = k.Hits
		}
	}
}

func (m *Management) exportKeys() ([]byte, error) {
	if err := m.db.FindAll(&m.keys); err != nil {
		return nil, err
	}

	data, err := json.Marshal(m.keys)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (m *Management) importKeys(keys []byte) error {
	var k []Keys
	if err := json.Unmarshal(keys, &k); err != nil {
		return err
	}

	m.keys = k
	if err := m.overrideData(); err != nil {
		return err
	}

	return nil
}

func GetKeys() []Keys {
	return mgmt.getKeys()
}

func HitKey(k Keys) error {
	return mgmt.hitKey(k)
}

func UnHitKey(k Keys) error {
	return mgmt.unHitKey(k)
}

func ExportKeys() ([]byte, error) {
	return mgmt.exportKeys()
}

func ImportKeys(keys []byte) error {
	return mgmt.importKeys(keys)
}
