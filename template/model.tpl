package model

import (
	"math/rand"
	"time"

	"gopkg.in/gorp.v1"
	null "gopkg.in/guregu/null.v3"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

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
	return rand.Intn(max)
}

func randString(length int) string {
	if length <= 0 {
		return ""
	}
	b := make([]byte, int(length))
	for i := range b {
		b[i] = baseString[int(rand.Int63()%int64(len(baseString)))]
	}
	return string(b)
}

func randStringRange(min, max int) string {
	if min >= max {
		return ""
	}
	length := int(rand.Int63() % int64(max - min))
	if length > 10000 {
		length = 10000 // limiter
	}
	return randString(length)
}

func randTime() time.Time {
	return time.Unix(rand.Int63n(int64(3000*365*24*60*60)), rand.Int63n(int64(time.Second)))
}

func randNullTime() null.Time {
  valid := rand.Intn(2) == 0
  if !valid {
    return null.Time{}
  }
  return null.TimeFrom(randTime())
}
