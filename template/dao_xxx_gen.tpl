// Automatically generated by gendao.
// Source: valencia_media/{{ .Table.Name }}

// ********************
// *** DO NOT EDIT! ***
// ********************

package dao

import (
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v1"{{range .Table.UseTypes}}{{if eq . "null.Time"}}
	"gopkg.in/guregu/null.v3"{{end}}{{end}}

	"{{ .Config.PackageRoot }}/dao/ranger"
	"{{ .Config.PackageRoot }}/model"
)
{{$TableNamePascal := .Table.NameByPascalcase}}
{{$TableNameCamel := .Table.NameByCamelcase}}
type (
	inner{{ $TableNamePascal }} interface { {{range .Table.CustomMethods}}
		{{template "part_method_name.tpl" .}}{{end}}
		Insert({{ $TableNameCamel }} *model.{{ $TableNamePascal }}) error
		Update({{ $TableNameCamel }} *model.{{ $TableNamePascal }}) error
		DeleteBy{{range .Table.PrimaryKey.Columns}}{{.NameByPascalcase}}{{end}}({{range .Table.PrimaryKey.Columns}}{{print .NameByCamelcase " " .Type}}{{end}}) error
	}
	// {{ $TableNamePascal }}Dao {{ $TableNameCamel }} dao struct
	{{ $TableNamePascal }}Dao struct {
		baseDao
	}
)

func new{{$TableNamePascal}}(dbm, dbs *gorp.DbMap) *{{$TableNamePascal}}Dao {
	m := model.{{$TableNamePascal}}{}
	tableName := m.TableName()
	pks := m.PrimaryKeys()
	dbm.AddTableWithName(m, tableName).SetKeys({{.Table.PrimaryKey.AutoIncrement}}, pks...)
	dbs.AddTableWithName(m, tableName).SetKeys({{.Table.PrimaryKey.AutoIncrement}}, pks...)
	dao := {{$TableNamePascal}}Dao{}
	dao.baseDao = newBaseDao(dbm, dbs)
	dao.tableName = tableName
	dao.columnsName = strings.Join(m.ColumnNames(), ",")
	return &dao
}

// ------------------------------
// Global Methods for interface
// ------------------------------

{{range .Table.CustomMethods}}
// {{.Name}} get {{$TableNameCamel}} with {{range $i, $p := .Params}}{{if ne $i 0}} and {{end}}{{.NameByCamelcase}}{{end}}
func (dao {{ $TableNamePascal }}Dao) {{template "part_method_name.tpl" .}} {
	builder := dao.newSelectBuilder(){{range .Params}}{{if .Where}}.
		Where(sq.Eq{"{{.Name}}": {{.NameByCamelcase}}}){{end}}{{end}}{{if .Orders}}.{{$desc := .Desc}}
		OrderBy({{range $i, $p := .Orders}}{{if ne $i 0}}, {{end}}"{{.Name}}{{if $desc}} desc{{end}}"{{end}}){{end}}{{range .Params}}{{if eq .Name "limit"}}.
		Limit(limit){{end}}{{end}}{{if .RangeParam}}{{if contains .RangeParam.Type "int"}}
	builder = ranger.SetWhereInt(builder, "{{.RangeParam.Name}}", rangeFncs){{else if contains .RangeParam.Type "string"}}
	builder = ranger.SetWhereStr(builder, "{{.RangeParam.Name}}", rangeFncs){{else if contains .RangeParam.Type "time"}}
	builder = ranger.SetWhereTime(builder, "{{.RangeParam.Name}}", rangeFncs){{end}}{{end}}{{if .ReturnMany}}
	return dao.findManyByBuilder(&builder){{else}}
	return dao.findOneByBuilder(&builder){{end}}
}
{{end}}

// Insert insert {{$TableNameCamel}}
func (dao {{ $TableNamePascal }}Dao) Insert({{$TableNameCamel}} *model.{{$TableNamePascal}}) error {
	return dao.insert({{$TableNameCamel}})
}

// Update update {{$TableNameCamel}}
func (dao {{ $TableNamePascal }}Dao) Update({{$TableNameCamel}} *model.{{$TableNamePascal}}) error {
	return dao.update({{$TableNameCamel}})
}

// DeleteBy{{range .Table.PrimaryKey.Columns}}{{.NameByPascalcase}}{{end}} delete {{$TableNameCamel}} by {{range .Table.PrimaryKey.Columns}}{{.Name}}{{end}}
func (dao {{ $TableNamePascal }}Dao) DeleteBy{{range .Table.PrimaryKey.Columns}}{{.NameByPascalcase}}{{end}}({{range .Table.PrimaryKey.Columns}}{{print .NameByCamelcase " " .Type}}{{end}}) error {
	m := &model.{{$TableNamePascal}}{ {{range .Table.PrimaryKey.Columns}}{{.NameByPascalcase}}{{end}}: {{range .Table.PrimaryKey.Columns}}{{.NameByCamelcase}}{{end}} }
	return dao.delete(m)
}

// ------------------
// Private Methods
// ------------------

func (dao {{$TableNamePascal}}Dao) findOneByBuilder(builder *sq.SelectBuilder) (*model.{{$TableNamePascal}}, error) {
	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "build sql failed [sql='%s'][args='%+v']", sql, args)
	}
	var {{$TableNameCamel}} model.{{$TableNamePascal}}
	if err := dao.dbs.SelectOne(&{{$TableNameCamel}}, sql, args...); err != nil {
		return nil, errors.Wrapf(err, "fetch data failed [sql='%s'][args='%+v']", sql, args)
	}
	return &{{$TableNameCamel}}, nil
}

func (dao {{$TableNamePascal}}Dao) findManyByBuilder(builder *sq.SelectBuilder) (model.{{$TableNamePascal}}Slice, error) {
	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "build sql failed [sql='%s'][args='%+v']", sql, args)
	}
	var {{$TableNameCamel}}s model.{{$TableNamePascal}}Slice
	if _, err := dao.dbs.Select(&{{$TableNameCamel}}s, sql, args...); err != nil {
		return nil, errors.Wrapf(err, "fetch data failed [sql='%s'][args='%+v']", sql, args)
	}
	return {{$TableNameCamel}}s, nil
}

func (dao {{$TableNamePascal}}Dao) insert({{$TableNameCamel}} *model.{{$TableNamePascal}}) error {
	return errors.Wrapf(dao.dbm.Insert({{$TableNameCamel}}), "insert failed [%+v]", {{$TableNameCamel}})
}

func (dao {{$TableNamePascal}}Dao) update({{$TableNameCamel}} *model.{{$TableNamePascal}}) error {
	_, err := dao.dbm.Update({{$TableNameCamel}})
	return errors.Wrapf(err, "update failed [%+v]", {{$TableNameCamel}})
}

func (dao {{$TableNamePascal}}Dao) delete({{$TableNameCamel}} *model.{{$TableNamePascal}}) error {
	_, err := dao.dbm.Delete({{$TableNameCamel}})
	return errors.Wrapf(err, "delete failed [%+v]", {{$TableNameCamel}})
}

//----------------------------------------
// Compiler Check
//----------------------------------------

var _ {{$TableNamePascal}} = &{{$TableNamePascal}}Dao{}
