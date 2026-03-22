package core

import "time"

type Note struct {
	ID        string
	Title     string
	Path      string
	Tags      []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Workspace struct {
	RootPath string}