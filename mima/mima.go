package mima

import (
	"fmt"

	"github.com/ahui2016/mima-web/util"
)

type Mima struct {
	ID        string
	Title     string
	Alias     string
	Username  string
	Password  string
	Notes     string
	CreatedAt string // ISO8601
	UpdatedAt string
	DeletedAt string
	Operation Operation
	History   []*History
}

func clone(m Mima) Mima {
	return m
}

func NewFrom(form *EditForm) *Mima {
	now := util.TimeNow()
	return &Mima{
		ID:        util.NewID(),
		Title:     form.Title,
		Alias:     "",
		Username:  form.Username,
		Password:  form.Password,
		Notes:     form.Notes,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: "",
		Operation: 0,
		History:   nil,
	}
}

func (m *Mima) UpdateFromForm(form *EditForm) (changeIndex, writeFrag bool) {
	if m.Alias != form.Alias {
		m.Alias = form.Alias
		writeFrag = true
	}
	if m.noChangeFrom(form) {
		return
	}

	updatedAt := util.TimeNow()
	m.makeHistory(updatedAt)

	m.Title = form.Title
	m.Username = form.Username
	m.Password = form.Password
	m.Notes = form.Notes
	m.UpdatedAt = updatedAt
	return true, true
}

func (m *Mima) makeHistory(updatedAt string) {
	h := &History{
		Title:    m.Title,
		Username: m.Username,
		Password: m.Password,
		Notes:    m.Notes,
		DateTime: updatedAt,
	}
	m.History = append(m.History, h)
}

func (m *Mima) DeleteHistory(datetime string) error {
	if i := m.getHistoryIndex(datetime); i < 0 {
		return fmt.Errorf("History not found: %s", datetime)
	} else {
		m.History = append(m.History[:i], m.History[i+1:]...)
		return nil
	}
}

func (m *Mima) getHistoryIndex(datetime string) int {
	for i, h := range m.History {
		if h.DateTime == datetime {
			return i
		}
	}
	return -1
}

func (m *Mima) noChangeFrom(form *EditForm) bool {
	s1 := m.Title + m.Username + m.Password + m.Notes
	s2 := form.Title + form.Username + form.Password + form.Notes
	return s1 == s2
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

// UpdateFromFrag updates m when fragment.Operation is Update.
func (m *Mima) UpdateFromFrag(fragment *Mima) (changeIndex bool) {
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
	m.DeletedAt = util.TimeNow()
}

func (m *Mima) UnDelete() {
	m.DeletedAt = ""
}

func (m *Mima) IsDeleted() bool {
	return m.DeletedAt != ""
}

type History struct {
	Title    string
	Username string
	Password string
	Notes    string
	DateTime string
}

type EditForm struct {
	ID       string
	Title    string
	Alias    string
	Username string
	Password string
	Notes    string
}
