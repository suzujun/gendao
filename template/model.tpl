package model

import (
	"math/rand"
	"time"

	"gopkg.in/gorp.v1"
)

// Model ...
type Model struct { {{range .CommonColumns}}
  {{ print .NameByPascalcase " " .Type "`db:\"" .Name "\"`" }}{{end}}
}

const baseString = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// PreInsert is previous insert func
func (m *Model) PreInsert(s gorp.SqlExecutor) error {
	now := time.Now().Round(time.Second)
	return m.preInsert(s, now)
}

func (m *Model) preInsert(_ gorp.SqlExecutor, now time.Time) error {
	m.UpdatedAt = now
	m.CreatedAt = now
	return nil
}

// PreUpdate is previous update func
func (m *Model) PreUpdate(s gorp.SqlExecutor) error {
	now := time.Now().Round(time.Second)
	return m.preUpdate(s, now)
}

func (m *Model) preUpdate(_ gorp.SqlExecutor, now time.Time) error {
	m.UpdatedAt = now
	return nil
}

// NewDummyModel is generate new dummy model
func NewDummyModel() Model {
	return Model{
		UpdatedAt: time.Unix(int64(rand.Uint32()), 0),
		CreatedAt: time.Unix(int64(rand.Uint32()), 0),
	}
}

func randIntn(max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max)
}

func randString(length int) string {
	if length <= 0 {
		return ""
	}
	b := make([]byte, int(length))
	for i := range b {
		rand.Seed(time.Now().UnixNano())
		b[i] = baseString[int(rand.Int63()%int64(len(baseString)))]
	}
	return string(b)
}

func randStringRange(min, max int) string {
	if min >= max {
		return ""
	}
	rand.Seed(time.Now().UnixNano())
	length := int(rand.Int63() % int64(max - min))
	if length > 10000 {
		length = 10000 // limiter
	}
	return randString(length)
}
