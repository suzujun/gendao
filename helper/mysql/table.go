package mysql

import (
	"encoding/json"

	"github.com/suzujun/gendao/helper"
)

type (
	Table struct {
		Catalog string   `json:"catalog"`
		Schema  string   `json:"schema"`
		Name    string   `json:"name"`
		Columns []Column `json:"columns"`
		Indexes []Index  `json:"indexes"`
	}
)

func (mt Table) WriteJSON(path string) error {
	jsonBytes, err := json.MarshalIndent(mt, "", "  ")
	if err != nil {
		return err
	}
	_, err = helper.CreateFile(path, jsonBytes)
	return err
}
