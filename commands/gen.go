package commands

import (
	"fmt"
	"path/filepath"
)

func writeTablesJSON(con *MysqlConnection, outputPath string) error {
	if err := createDirIfNotExist(outputPath); err != nil {
		return err
	}
	tables, err := con.GetTableNames()
	if err != nil {
		return err
	}
	for _, tname := range tables {
		mt, err := con.GetMysqlTable(tname)
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
