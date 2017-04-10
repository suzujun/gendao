package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type (
	Command struct {
		Config Config
		ReadAt time.Time
	}
)

// NewCommandFromJSON new command from json file
func NewCommandFromJSON(configPath, dbName string) (*Command, error) {
	b, err := readFile(configPath)
	if err != nil {
		return nil, err
	}
	com := Command{}
	if err := com.Config.ParseJSON(b); err != nil {
		return nil, err
	}
	if dbName != "" {
		com.Config.MysqlConfig.DbName = dbName
	}
	com.ReadAt = time.Now()
	return &com, nil
}

// GenerateJSON generate json file
func (cmd Command) GenerateJSON() error {
	dbconf := cmd.Config.MysqlConfig
	con, err := NewConnection(dbconf.User, dbconf.Password, dbconf.DbName, false)
	if err != nil {
		return err
	}
	defer con.Close()
	outputPath := strings.Replace(cmd.Config.OutputJSONPath, "{dbname}", cmd.Config.MysqlConfig.DbName, -1)
	err = writeTablesJSON(con, outputPath)
	return err
}

func (cmd Command) GenerateSourceFromJSON(table string) error {

	config := cmd.Config
	targetTables := []string{}
	if table != "" {
		targetTables = strings.Split(table, ",")
	}

	// check json path
	dbname := config.MysqlConfig.DbName
	if dbname == "" {
		return errors.New("No database name selected in config")
	}
	path := strings.Replace(config.OutputJSONPath, "{dbname}", dbname, -1)
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("json path must be directory")
	}

	// check and craete output path
	if err := createDirIfNotExist(config.OutputSourcePath); err != nil {
		return err
	}

	myTemplate, err := NewTemplate(config.InputTemplatePath, config.TemplateToTableLoop, config.OutputSourcePath)
	if err != nil {
		return err
	}

	var pTables []TemplateDataTable
	var outputSource = func(path string) error {
		var table MysqlTable
		if err := readFileJSON(path, &table); err != nil {
			return err
		}
		fmt.Println("file:", path)
		pTable := newTamplateParamTable(cmd.Config.PackageRoot, table, config.CommonColumns, config.CustomColumnType)
		data := TemplateData{
			Config: cmd.Config,
			Table:  pTable,
		}
		if err := myTemplate.outputSourceFileTable(data); err != nil {
			return err
		}
		pTables = append(pTables, pTable)
		return nil
	}

	if len(targetTables) > 0 {
		// Only specified file
		for _, table := range targetTables {
			path := fmt.Sprintf("%s/%s.json", path, table)
			if err := outputSource(path); err != nil {
				return err
			}
		}
	} else {
		// Target all files on the specified path
		if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if info == nil || info.IsDir() || err != nil {
				return err
			}
			name := info.Name()
			tableName := name[:len(name)-5] // remove ".json"
			if len(targetTables) == 0 && stringsContains(cmd.Config.IgnoreTableNames, tableName) {
				fmt.Println("file:", name, "[ignore]")
				return nil
			}
			return outputSource(path)
		}); err != nil {
			return err
		}
	}

	// --------------
	// template by once
	// --------------

	if len(config.TemplateByOnce) == 0 {
		return nil
	}

	myTemplate, err = NewTemplate(config.InputTemplatePath, config.TemplateByOnce, config.OutputSourcePath)
	if err != nil {
		return err
	}

	data := TemplateData{
		Config: cmd.Config,
	}

	// set common common column
	ccLen := len(cmd.Config.CommonColumns)
	if ccLen > 0 {
		for _, table := range pTables {
			if cols := table.CommonColumns(); len(cols) == ccLen {
				data.CommonColumns = cols
				break
			}
		}
	}
	return myTemplate.outputSourceFileTable(data)
}
