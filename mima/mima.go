package mima

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

type History struct {
	Title    string
	Username string
	Password string
	Notes    string
	DateTime string
}
