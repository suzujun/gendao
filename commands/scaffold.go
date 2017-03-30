package commands

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/pkg/errors"
)

type (
	// MyTemplate ...
	MyTemplate struct {
		Template      *template.Template
		ExportConfigs []TemplateExportConfig
	}
	TemplateExportConfig struct {
		TemplateName   string
		ExportPathName string
		Overwrite      bool
	}
	// TemplateData ...
	TemplateData struct {
		Config        Config
		Table         TemplateDataTable
		CommonColumns []TemplateDataColumn
	}
	// TemplateDataTable ...
	TemplateDataTable struct {
		Name             string
		NameByCamelcase  string
		NameByPascalcase string
		ColumnsName      string
		Columns          []TemplateDataColumn
		PrimaryKey       TemplateDataIndex
		Indexs           []TemplateDataIndex
		UsePackages      [][]string
		CustomMethods    []CustomMethod
	}
	// TemplateDataColumn ...
	TemplateDataColumn struct {
		Name             string
		NameByCamelcase  string
		NameByPascalcase string
		Type             string
		Primary          bool
		Common           bool
		AutoIncrement    bool
		Unique           bool
		SampleValue      string
		MysqlColumn
	}
	// TemplateDataIndex ...
	TemplateDataIndex struct {
		Name          string
		Columns       []TemplateDataColumn
		Primary       bool
		AutoIncrement bool
		Unique        bool
	}
)

const randomBaseCharcter = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func NewTemplate(inputPath string, tmplFiles []TemplateFile, outputPath string) (*MyTemplate, error) {
	funcMap := template.FuncMap{
		"title": strings.Title,
		"now": func() string {
			return time.Now().Format(time.RFC3339)
		},
		"contains": stringsContains,
	}
	// load template
	files := make([]string, len(tmplFiles))
	configs := []TemplateExportConfig{}
	for i, target := range tmplFiles {
		ps := []string{inputPath, target.Name}
		files[i] = strings.Join(ps, "/")
		// check output dir
		path := outputPath + "/" + target.ExportName
		last := strings.LastIndex(path, "/")
		if err := createDirIfNotExist(path[:last]); err != nil {
			return nil, errors.Errorf("%s, path=[%s]", err, path[:last])
		}
		if target.ExportName == "" {
			continue
		}
		// add config
		configs = append(configs, TemplateExportConfig{
			TemplateName:   target.Name,
			ExportPathName: path,
			Overwrite:      target.Overwrite,
		})
	}
	tp := MyTemplate{
		Template:      template.Must(template.New("default").Funcs(funcMap).ParseFiles(files...)),
		ExportConfigs: configs,
	}
	return &tp, nil
}

func (my MyTemplate) outputSourceFileTable(data TemplateData) error {
	for _, config := range my.ExportConfigs {
		tmpl := my.Template.Lookup(config.TemplateName)
		buff := bytes.NewBuffer([]byte{})
		if err := tmpl.Execute(buff, data); err != nil {
			return err
		}
		name := ConvSingularize(data.Table.Name)
		path := strings.Replace(config.ExportPathName, "{name}", name, -1)
		var res int
		var err error
		if config.Overwrite {
			res, err = createFile(path, buff.Bytes())
		} else {
			res, err = createFileIfNotExist(path, buff.Bytes())
		}
		if err != nil {
			return err
		}
		if res == 0 {
			fmt.Println("generate:", path, "[skip if exist]")
		} else {
			fmt.Println("generate:", path)
		}
	}
	return nil
}

func newTamplateParamTable(table MysqlTableJSON, commonColumns []string, customTypeMap map[string]*CustomColumnType) TemplateDataTable {
	pTable := TemplateDataTable{}
	if len(table.Columns) == 0 {
		return pTable
	}
	pTable.Name = table.Columns[0].TableName
	pTable.NameByCamelcase = ConvCamelcaseSingularize(table.Columns[0].TableName, false)
	pTable.NameByPascalcase = ConvPascalcaseSingularize(table.Columns[0].TableName, false)
	pTable.Columns = []TemplateDataColumn{}

	// get column info
	names := []string{}
	packageMap := map[string]string{}
	for _, column := range table.Columns {
		name := fmt.Sprintf("%s.%s", pTable.Name, column.ColumnName)
		customType := customTypeMap[name]
		tpColumn := newTemplateParamColumn(column, commonColumns, customType)
		pTable.Columns = append(pTable.Columns, tpColumn)
		names = append(names, column.ColumnName)
		// set use package
		if customType != nil && customType.Package != "" {
			packageMap[customType.Package] = customType.PackageAlias
		}
		if pkg := tpColumn.getUsePackage(); pkg != "" {
			packageMap[pkg] = ""
		}
	}
	// get index info
	indexs := newTemplateParamIndex(table.Indexs, pTable.Columns)
	for _, index := range indexs {
		if index.Primary {
			pTable.PrimaryKey = index
		} else {
			pTable.Indexs = append(pTable.Indexs, index)
		}
		pTable.CustomMethods = append(pTable.CustomMethods, GenCustomMethods(index, pTable.NameByPascalcase)...)
	}
	// using packages
	if len(packageMap) > 0 {
		pTable.UsePackages = make([][]string, 3)
		for pkg, alias := range packageMap {
			v := fmt.Sprintf("%s \"%s\"", alias, pkg)
			if m, _ := regexp.MatchString("^[a-z0-9/]+$", pkg); m { // standard library
				pTable.UsePackages[0] = append(pTable.UsePackages[0], v)
			// } else if m, _ := regexp.MatchString("^gopath/", pkg); m { // TODO match gopath
			// 	pTable.UsePackages[2] = append(pTable.UsePackages[2], v)
			} else {
				pTable.UsePackages[1] = append(pTable.UsePackages[1], v) // other library
			}
		}
	}
	// set columns name
	pTable.ColumnsName = strings.Join(names, ",")

	return pTable
}

func newTemplateParamColumn(column MysqlColumn, commonColumns []string, customType *CustomColumnType) TemplateDataColumn {
	tpc := TemplateDataColumn{MysqlColumn: column}
	tpc.Name = column.ColumnName
	tpc.NameByCamelcase = ConvCamelcase(column.ColumnName, true)
	tpc.NameByPascalcase = ConvPascalcase(column.ColumnName, true)
	tpc.Primary = column.Primary()
	tpc.Unique = column.Unique()
	tpc.Common = stringsContains(commonColumns, column.ColumnName)
	if customType != nil {
		tpc.Type = customType.Type
		tpc.SampleValue = customType.SampleValue
	} else {
		tpc.setType()
		tpc.setSampleValue()
	}
	tpc.AutoIncrement = column.AutoIncrement()
	return tpc
}

func newTemplateParamIndex(indexs []MysqlIndex, pColumns []TemplateDataColumn) []TemplateDataIndex {
	pIndexs := []TemplateDataIndex{}
	if len(indexs) == 0 || len(pColumns) == 0 {
		return pIndexs
	}
	columnMap := map[string]*TemplateDataColumn{}
	for i, pColumn := range pColumns {
		columnMap[pColumn.Name] = &pColumns[i]
	}
	var pIndex *TemplateDataIndex
	var beforeIndexName string
	for _, index := range indexs {
		if beforeIndexName != index.IndexName {
			if pIndex != nil {
				pIndexs = append(pIndexs, *pIndex)
			}
			pIndex = &TemplateDataIndex{Name: index.IndexName}
			pIndex.Unique = index.NonUnique == 0
			beforeIndexName = index.IndexName
		}
		pIndex.Primary = index.IndexName == "PRIMARY"
		if pColumn := columnMap[index.ColumnName]; pColumn != nil {
			pIndex.Columns = append(pIndex.Columns, *pColumn)
		}
	}
	pIndexs = append(pIndexs, *pIndex)
	for i, pIndex := range pIndexs {
		for _, column := range pIndex.Columns {
			pIndexs[i].AutoIncrement = column.AutoIncrement
			break
		}
	}
	return pIndexs
}

func (tdt *TemplateDataTable) CommonColumns() []TemplateDataColumn {
	res := []TemplateDataColumn{}
	for _, column := range tdt.Columns {
		if column.Common {
			res = append(res, column)
		}
	}
	return res
}

func (tdc *TemplateDataColumn) getUsePackage() string {
	if tdc.Type == "time.Time" {
		return "time"
	} else if strings.Index(tdc.Type, "null.") == 0 {
		return "gopkg.in/guregu/null.v3"
	} else if tdc.Type == "string" && tdc.Primary {
		return "fmt"
	}
	return ""
}

func (tdc *TemplateDataColumn) setType() {
	tdc.Type = (func(mc MysqlColumn) string {
		unsigned := mc.Unsigned()
		switch mc.DataType {
		case "char", "varchar", "enum", "set":
			if mc.IsNullable {
				return "null.String"
			}
			return "string" // TODO unsigned で "*string" しなくてよいか？
		case "tinyint":
			if mc.IsNullable {
				return "null.Int"
			} else if unsigned {
				return "uint8" // 0 to 255
			}
			return "int8" // -128 to 127
		case "smallint":
			if mc.IsNullable {
				return "null.Int"
			} else if unsigned {
				return "uint16" // 0 to 65535
			}
			return "int16" // -32768 to 32767
		case "mediumint":
			if mc.IsNullable {
				return "null.Int"
			} else if unsigned {
				return "uint32" // unsigned: 0 to 16777215
			}
			return "int32" // signed: -8388608 to 8388607
		case "bigint":
			if mc.IsNullable {
				return "null.Int"
			} else if unsigned {
				return "uint64"
			}
			return "int64"
		case "int", "integer":
			if mc.IsNullable {
				return "null.Int"
			} else if unsigned {
				return "uint32" // 0 to 4294967295
			}
			return "int32" // -2147483648 to 2147483647
		case "float":
			if mc.IsNullable {
				return "null.Float"
			}
			return "float32"
		case "double", "decimal", "dec":
			if mc.IsNullable {
				return "null.Float"
			}
			return "float64"
		case "date", "datetime", "timestamp", "time":
			if mc.IsNullable {
				return "null.Time"
			}
			return "time.Time"
		default:
			return "interface{}"
		}
	})(tdc.MysqlColumn)
}

func (tdc *TemplateDataColumn) setSampleValue() {
	tdc.SampleValue = (func(c *TemplateDataColumn) string {
		if c.Type == "string" || c.Type == "null.String" {
			max := int(*c.MysqlColumn.CharacterMaximumLength)
			min := int(max / 3)
			if c.Type == "null.String" {
				return fmt.Sprintf("randNullStringRange(%d, %d)", min, max)
			}
			return fmt.Sprintf("randStringRange(%d, %d)", min, max)
		} else if c.Type == "time.Time" {
			return "time.Unix(time.Now().Unix(), 0)"
		} else if strings.Index(c.Type, "null.") == 0 {
			switch c.Type {
			case "null.Int":
				r := c.MysqlColumn.DataTypeRange()
				return fmt.Sprintf("randNullInt(%d)", r.Max)
			case "null.Float":
				return "randNullFloat()"
			case "null.Time":
				return "randNullTime()"
			}
		} else if strings.Index(c.Type, "int") >= 0 {
			switch c.Type {
			case "int":
				return "rand.Int()"
			case "int64":
				return "rand.Int63()"
			case "int32":
				return "rand.Int31()"
			case "int16":
				return "int16(rand.Intn(32767))"
			case "int8":
				return "int8(rand.Intn(127))"
			case "uint":
				return "uint(rand.Uint64())"
			case "uint64":
				// return "rand.Uint64()" // fail at mysql: uint64 values with high bit set are not supported
				return "uint64(rand.Uint32())"
			case "uint32":
				return "uint32(rand.Intn(4294967295))"
			case "uint16":
				return "uint16(rand.Intn(65535))"
			case "uint8":
				return "uint8(rand.Intn(255))"
			default:
				return "1"
			}
		} else if strings.Index(c.Type, "float") == 0 {
			switch c.Type {
			case "float32":
				return "rand.Float32()" // TODO mysql variable type に応じた値を設定する
			case "float64":
				return "rand.Float64()" // TODO mysql variable type に応じた値を設定する
			default:
				return "1.01"
			}
		}
		return "" // TODO バイナリとかの対応
	})(tdc)
}
