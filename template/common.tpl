package model

import (
	"math/rand"
	"time"

	"gopkg.in/gorp.v1"
)

type Model struct {
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (m *Model) PreInsert(s gorp.SqlExecutor) error {
	now := time.Now().Round(time.Second)
	return m.preInsert(s, now)
}

func (m *Model) preInsert(_ gorp.SqlExecutor, now time.Time) error {
	m.UpdatedAt = now
	m.CreatedAt = now
	return nil
}

func (m *Model) PreUpdate(s gorp.SqlExecutor) error {
	now := time.Now().Round(time.Second)
	return m.preUpdate(s, now)
}

func (m *Model) preUpdate(_ gorp.SqlExecutor, now time.Time) error {
	m.UpdatedAt = now
	return nil
}

func NewDummyModel() Model {
	return Model{
		UpdatedAt: time.Unix(int64(rand.Uint32()), 0),
		CreatedAt: time.Unix(int64(rand.Uint32()), 0),
	}
}
