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
