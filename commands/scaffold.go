package commands

import (
	"bytes"
	"fmt"
	"path/filepath"
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
		Name                 string
		NameByCamelcase      string
		NameByPascalcase     string
		ColumnsName          string
		Columns              []TemplateDataColumn
		PrimaryKey           TemplateDataIndex
		Indexes              []TemplateDataIndex
		UsePackages          [][]string
		UseTypes             []string
		CustomMethods        []CustomMethod
		CustomMethodUseTypes []string
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
	configs := make([]TemplateExportConfig, 0, len(tmplFiles))
	for i, target := range tmplFiles {
		files[i] = filepath.Join(inputPath, target.Name)
		// check exists template file
		if !IsFileExist(files[i]) {
			return nil, errors.Errorf("not found template file, [%s]", files[i])
		}
		// check output dir
		path := filepath.Join(outputPath, target.ExportName)
		dirPath := filepath.Dir(path)
		if err := createDirIfNotExist(dirPath); err != nil {
			return nil, errors.Errorf("%s, path=[%s]", err, dirPath)
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
		name := NewWordConverter(data.Table.Name).Singularize().ToString()
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

var stdlibReg = regexp.MustCompile("^[a-z0-9/]+$")

func newTamplateParamTable(packageRoot string, table MysqlTable, commonColumns []string, customTypeMap map[string]*CustomColumnType) TemplateDataTable {
	pTable := TemplateDataTable{}
	if len(table.Columns) == 0 {
		return pTable
	}
	pTable.Name = table.Columns[0].TableName
	pTable.NameByCamelcase = NewWordConverter(table.Columns[0].TableName).Camelcase().Singularize().ToString()
	pTable.NameByPascalcase = NewWordConverter(table.Columns[0].TableName).Pascalcase().Singularize().ToString()
	pTable.Columns = make([]TemplateDataColumn, 0, len(table.Columns))

	// get column info
	names := make([]string, 0, len(table.Columns))
	packageMap := make(map[string]string, len(table.Columns))
	typeMap := map[string]bool{}
	for _, column := range table.Columns {
		name := fmt.Sprintf("%s.%s", pTable.Name, column.ColumnName)
		customType := customTypeMap[name]
		tpColumn := newTemplateParamColumn(column, commonColumns, customType)
		pTable.Columns = append(pTable.Columns, tpColumn)
		names = append(names, column.ColumnName)
		typeMap[tpColumn.Type] = true
		// set use package
		if customType != nil && customType.Package != "" {
			packageMap[customType.Package] = customType.PackageAlias
		}
		if pkg := tpColumn.getUsePackage(); pkg != "" {
			packageMap[pkg] = ""
		}
	}
	// get index info
	indexes := newTemplateParamIndex(table.Indexes, pTable.Columns)
	methods := make([]CustomMethod, 0, len(indexes))
	pTable.Indexes = make([]TemplateDataIndex, 0, len(indexes))
	for _, index := range indexes {
		if index.Primary {
			pTable.PrimaryKey = index
		} else {
			pTable.Indexes = append(pTable.Indexes, index)
		}
		methods = append(methods, GenCustomMethods(index, pTable.NameByPascalcase)...)
	}
	// deduplication
	methodMap := map[string]bool{}
	pTable.CustomMethods = make([]CustomMethod, 0, len(methods))
	for _, m := range methods {
		if methodMap[m.Name] {
			continue
		}
		methodMap[m.Name] = true
		pTable.CustomMethods = append(pTable.CustomMethods, m)
	}
	// set using type for method params
	uniquer := NewUniquer()
	for _, m := range methods {
		for _, p := range m.Params {
			uniquer.Add(p.Type)
		}
	}
	pTable.CustomMethodUseTypes = uniquer.Uniq()
	// set using types
	idx := 0
	pTable.UseTypes = make([]string, len(typeMap))
	for key := range typeMap {
		pTable.UseTypes[idx] = key
		idx++
	}
	// using packages
	if len(packageMap) > 0 {
		pTable.UsePackages = make([][]string, 3)
		for pkg, alias := range packageMap {
			v := fmt.Sprintf("%s \"%s\"", alias, pkg)
			if stdlibReg.MatchString(pkg) { // standard library
				pTable.UsePackages[0] = append(pTable.UsePackages[0], v)
			} else if m, _ := regexp.MatchString("^"+packageRoot, pkg); m { // my package
				pTable.UsePackages[2] = append(pTable.UsePackages[2], v)
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
	tpc.NameByCamelcase = NewWordConverter(column.ColumnName).Camelcase().Lint().ToString()
	tpc.NameByPascalcase = NewWordConverter(column.ColumnName).Pascalcase().Lint().ToString()
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

func newTemplateParamIndex(indexes []MysqlIndex, pColumns []TemplateDataColumn) []TemplateDataIndex {
	pIndexes := []TemplateDataIndex{}
	if len(indexes) == 0 || len(pColumns) == 0 {
		return pIndexes
	}
	columnMap := make(map[string]*TemplateDataColumn, len(pColumns))
	for i, pColumn := range pColumns {
		columnMap[pColumn.Name] = &pColumns[i]
	}
	var pIndex *TemplateDataIndex
	var beforeIndexName string
	for _, index := range indexes {
		if beforeIndexName != index.IndexName {
			if pIndex != nil {
				pIndexes = append(pIndexes, *pIndex)
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
	pIndexes = append(pIndexes, *pIndex)
	for i, pIndex := range pIndexes {
		for _, column := range pIndex.Columns {
			pIndexes[i].AutoIncrement = column.AutoIncrement
			break
		}
	}
	return pIndexes
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
			min := max / 3
			if c.Type == "null.String" {
				return fmt.Sprintf("randNullStringRange(%d, %d)", min, max)
			}
			return fmt.Sprintf("randStringRange(%d, %d)", min, max)
		} else if c.Type == "time.Time" {
			return "time.Unix(time.Now().Unix(), 0)"
		} else if strings.HasPrefix(c.Type, "null.") {
			switch c.Type {
			case "null.Int":
				r := c.MysqlColumn.DataTypeRange()
				return fmt.Sprintf("randNullInt(%d)", r.Max)
			case "null.Float":
				return "randNullFloat()"
			case "null.Time":
				return "randNullTime()"
			}
		} else if strings.Contains(c.Type, "int") {
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
		} else if strings.HasPrefix(c.Type, "float") {
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
