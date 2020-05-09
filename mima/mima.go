package mima

import (
	"time"

	"github.com/ahui2016/mima-web/util"
)

type Mima struct {
	ID        string
	Title     string
	Alias     string
	Username  string
	Password  string
	Notes     string
	CreatedAt int64
	UpdatedAt int64
	DeletedAt int64
	Operation Operation
	History   []*History
}

func NewFrom(form *AddForm) *Mima {
	now := time.Now().UnixNano()
	return &Mima{
		ID:        util.NewID(),
		Title:     form.Title,
		Alias:     "",
		Username:  form.Username,
		Password:  form.Password,
		Notes:     form.Notes,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: 0,
		Operation: 0,
		History:   nil,
	}
}

func clone(m Mima) Mima {
	return m
}

func (m *Mima) HideSecrets() *Mima {
	m2 := clone(*m)
	if m2.Password != "" {
		m2.Password = "******"
	}
	m2.Notes = ""
	m2.History = nil
	return &m2
}

// UpdateFrom updates m when fragment.Operation is Update.
func (m *Mima) UpdateFrom(fragment *Mima) (changeIndex bool) {
	m.Alias = fragment.Alias
	m.History = fragment.History
	if m.UpdatedAt == fragment.UpdatedAt {
		return false
	}
	m.Title = fragment.Title
	m.Username = fragment.Username
	m.Password = fragment.Password
	m.Notes = fragment.Notes
	m.UpdatedAt = fragment.UpdatedAt
	return true
}

// Delete soft-delete itself.
func (m *Mima) Delete() {
	m.DeletedAt = time.Now().UnixNano()
}

func (m *Mima) UnDelete() {
	m.DeletedAt = 0
}

func (m *Mima) IsDeleted() bool {
	return m.DeletedAt > 0
}

type History struct {
	Title     string
	Username  string
	Password  string
	Notes     string
	Timestamp int64
}

type AddForm struct {
	Title    string
	Username string
	Password string
	Notes    string
}
