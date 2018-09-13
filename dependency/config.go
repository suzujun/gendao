package dependency

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/suzujun/gendao/helper"
)

type (
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
		Host     string `json:"host"`
		Port     string `json:"port"`
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

func NewConfig(host, port, user, password, database string) Config {
	conf := newConfig()
	if host != "" {
		conf.MysqlConfig.Host = host
	}
	if port != "" {
		conf.MysqlConfig.Port = port
	}
	if user != "" {
		conf.MysqlConfig.User = user
	}
	if password != "" {
		conf.MysqlConfig.Password = password
	}
	if database != "" {
		conf.MysqlConfig.DbName = database
	}
	return conf
}

func newConfig() Config {
	// default setting values
	return Config{
		MysqlConfig: MysqlConfig{
			Host:     "localhost",
			Port:     "3306",
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
		OutputJSONPath:    "./out/{dbname}",
		OutputSourcePath:  "./src",
		InputTemplatePath: "./template",
		TemplateByOnce:    []TemplateFile{
		// {Name: "model.tpl", ExportName: "model/model.go"}, // dao/model.go
		},
		TemplateToTableLoop: []TemplateFile{
			{Name: "dao_xxx.tpl", ExportName: "dao/{name}.go", Overwrite: false},        // dao/channel.go
			{Name: "dao_xxx_gen.tpl", ExportName: "dao/{name}_gen.go", Overwrite: true}, // dao/channel_gen.go
			{Name: "model_xxx.tpl", ExportName: "model/{name}.go", Overwrite: true},     // model/channel.go
			{Name: "part_method_name.tpl"},
		},
		CustomColumnType: map[string]*CustomColumnType{},
	}
}

func (c Config) Write(path string) error {
	b, err := c.ExportJSON()
	if err != nil {
		return err
	}
	_, err = helper.CreateFile(path, b)
	return err
}

func (c Config) ExportJSON() ([]byte, error) {
	return json.MarshalIndent(c, "", "  ")
}

func (c *Config) ParseJSON(data []byte) error {
	return json.Unmarshal(data, c)
}

func getPackageRoot() string {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("error: ", err.Error())
		return ""
	}
	gopath := os.Getenv("GOPATH")
	srcPath := filepath.Join(gopath, "src")
	rel, err := filepath.Rel(srcPath, pwd)
	if err != nil {
		fmt.Println("error: ", err.Error())
		return ""
	}
	return filepath.ToSlash(rel)
}
