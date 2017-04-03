package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type (
	Command struct {
		Config Config
		ReadAt time.Time
	}
	Config struct {
		PackageRoot         string                       `json:"packageRoot"`
		CommonColumns       []string                     `json:"commonColumns"`
		MysqlConfig         MysqlConfig                  `json:"mysqlConfig"`
		OutputJSONPath      string                       `json:"outputJsonPath"`
		OutputSourcePath    string                       `json:"outputSourcePath"`
		InputTemplatePath   string                       `json:"inputTemplatePath"`
		TemplateByOnce      []TemplateFile               `json:"templateByOnce"`
		TemplateToTableLoop []TemplateFile               `json:"templateToTableLoop"`
		IgnoreTableNames    []string                     `json:"ignoreTableNames"`
		CustomColumnType    map[string]*CustomColumnType `json:"customColumnTypes"`
	}
	TemplateFile struct {
		Name       string `json:"name"`
		ExportName string `json:"exportName"`
		Overwrite  bool   `json:"overwrite"`
	}
	MysqlConfig struct {
		User     string `json:"user"`
		Password string `json:"password"`
		DbName   string `json:"dbName"`
	}
	CustomColumnType struct {
		Type         string `json:"type"`
		SampleValue  string `json:"sampleValue"`
		Package      string `json:"package"`
		PackageAlias string `json:"packageAlias"`
	}
)

func newConfig() Config {
	// default setting values
	return Config{
		MysqlConfig: MysqlConfig{
			User:     "root",
			Password: "",
			DbName:   "",
		},
		PackageRoot: getPackageRoot(),
		IgnoreTableNames: []string{
			"goose_db_version",
		},
		CommonColumns: []string{
			"created_at",
			"updated_at",
		},
		OutputJSONPath:    "./out",
		OutputSourcePath:  "./src",
		InputTemplatePath: "./template",
		TemplateByOnce:    []TemplateFile{
		// {Name: "model.tpl", ExportName: "model/model.go"}, // dao/model.go
		},
		TemplateToTableLoop: []TemplateFile{
			{Name: "dao_xxx.tpl", ExportName: "dao/{name}.go", Overwrite: false},          // dao/channel.go
			{Name: "dao_xxx_gen.tpl", ExportName: "dao/{name}_gen.go", Overwrite: true},   // dao/channel_gen.go
			{Name: "dao_xxx_mock.tpl", ExportName: "dao/{name}_mock.go", Overwrite: true}, // dao/channel_mock.go
			{Name: "model_xxx.tpl", ExportName: "model/{name}.go", Overwrite: true},       // model/channel.go
			{Name: "part_method_name.tpl"},
		},
		CustomColumnType: map[string]*CustomColumnType{
		// "table_name.column_name": &CustomColumnType{
		// 	Type:        "bool",
		// 	SampleValue: "true",
		// 	Package:     "github.com/path/to",
		// 	PackageAlias:     "alias_name",
		// },
		},
	}
}

func (c Config) Write(path string) error {
	b, err := c.ExportJSON()
	if err != nil {
		return err
	}
	_, err = createFile(path, b)
	return err
}

func (c Config) ExportJSON() ([]byte, error) {
	return json.MarshalIndent(c, "", "  ")
}

func (c *Config) ParseJSON(data []byte) error {
	return json.Unmarshal(data, c)
}

func GenerateConfigJSON(user, password, database string) ([]byte, error) {
	conf := newConfig()
	if user != "" {
		conf.MysqlConfig.User = user
	}
	if password != "" {
		conf.MysqlConfig.Password = password
	}
	if database != "" {
		conf.MysqlConfig.DbName = database
	}
	return json.MarshalIndent(conf, "", "  ")
}

func GetCommand(configPath, dbName string) (*Command, error) {
	b, err := readFile(configPath)
	if err != nil {
		return nil, err
	}
	var com Command
	if err := json.Unmarshal(b, &com.Config); err != nil {
		return nil, err
	}
	if dbName != "" {
		com.Config.MysqlConfig.DbName = dbName
	}
	com.ReadAt = time.Now()
	return &com, nil
}

func (cmd Command) GenerateJSON() error {
	dbconf := cmd.Config.MysqlConfig
	con, err := NewConnection(dbconf.User, dbconf.Password, dbconf.DbName, false)
	if err != nil {
		return err
	}
	defer con.Close()
	outputPath := cmd.Config.OutputJSONPath + "/" + cmd.Config.MysqlConfig.DbName
	err = getAndOutputJSON(con, outputPath)
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
	path := config.OutputJSONPath + "/" + dbname
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

	myTemplate, tmpErr := NewTemplate(config.InputTemplatePath, config.TemplateToTableLoop, config.OutputSourcePath)
	if tmpErr != nil {
		return tmpErr
	}

	var pTables []TemplateDataTable
	if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		name := info.Name()
		tableName := name[:len(name)-5]
		if len(targetTables) > 0 && !stringsContains(targetTables, tableName) {
			fmt.Println("file:", name, "[skip]")
			return nil
		}
		if stringsContains(cmd.Config.IgnoreTableNames, tableName) {
			fmt.Println("file:", name, "[ignore]")
			return nil
		}
		fmt.Println("file:", name)

		// read json file
		var table MysqlTableJSON
		if err := readFileJSON(path, &table); err != nil {
			return err
		}

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
	}); err != nil {
		return err
	}

	// --------------
	// template by once
	// --------------

	myTemplate, tmpErr = NewTemplate(config.InputTemplatePath, config.TemplateByOnce, config.OutputSourcePath)
	if tmpErr != nil {
		return tmpErr
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

func filterCommonColumns(tables []TemplateDataTable) []TemplateDataColumn {
	for _, table := range tables {
		res := []TemplateDataColumn{}
		for _, column := range table.Columns {
			if column.Common {
				res = append(res, column)
			}
		}
		if len(res) > 0 {
			return res
		}
	}
	return nil
}

func getPackageRoot() string {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("error: ", err.Error())
		return ""
	}
	gopath := os.Getenv("GOPATH")
	re := regexp.MustCompile(fmt.Sprintf("^%s/src/", gopath))
	return re.ReplaceAllString(pwd, "")
}
