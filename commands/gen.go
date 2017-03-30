package commands

import (
	"encoding/json"
	"fmt"
)

func getAndOutputJSON(con *MysqlConnection, outputPath string) error {
	if err := createDirIfNotExist(outputPath); err != nil {
		return err
	}
	tables, err := con.GetTables()
	if err != nil {
		return err
	}
	var columns []MysqlColumn
	var indexs []MysqlIndex
	for _, tname := range tables {
		columns, err = con.GetColumns(tname)
		if err != nil {
			return err
		}
		indexs, err = con.GetIndexs(tname)
		if err != nil {
			return err
		}
		data := map[string]interface{}{}
		data["columns"] = columns
		data["indexs"] = indexs
		if len(columns) > 0 {
			data["catalog"] = columns[0].TableCatalog
			data["schema"] = columns[0].TableSchema
			data["name"] = columns[0].TableName
		}
		jsonBytes := []byte{}
		jsonBytes, err = json.MarshalIndent(data, "", "  ")
		if err != nil {
			return err
		}
		path := outputPath + "/" + tname + ".json"
		fmt.Println("genereate:", path)
		if _, err := createFile(path, jsonBytes); err != nil {
			return err
		}
	}
	return nil
}
