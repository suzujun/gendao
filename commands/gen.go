package commands

import (
	"fmt"
	"path/filepath"

	"github.com/suzujun/gendao/helper"
	"github.com/suzujun/gendao/helper/mysql"
)

func writeTablesJSON(con *mysql.Connection, outputPath string) error {
	if err := helper.CreateDirIfNotExist(outputPath); err != nil {
		return err
	}
	tables, err := con.GetTableNames()
	if err != nil {
		return err
	}
	for _, tname := range tables {
		mt, err := con.GetTable(tname)
		if err != nil {
			return err
		}
		path := filepath.Join(outputPath, tname+".json")
		if err := mt.WriteJSON(path); err != nil {
			return err
		}
		fmt.Println("genereate:", path)
	}
	return nil
}
