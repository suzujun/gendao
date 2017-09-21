package scaffold

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/suzujun/gendao/dependency"
	"github.com/suzujun/gendao/helper"
	"github.com/suzujun/gendao/helper/mysql"
)

func TestScaffold_NewTemplate(t *testing.T) {
	require := require.New(t)
	dir, err := ioutil.TempDir("", "")
	require.NoError(err)

	// generate template file
	filenames := make([]string, 3)
	for i := range filenames {
		filenames[i] = fmt.Sprintf("template_%05d.tmp", i)
		_, err := helper.CreateFile(filepath.Join(dir, filenames[i]), fmt.Sprintf("template%05d", i))
		require.NoError(err)
	}

	assert := assert.New(t)

	// sample templateFile
	tmpfiles := []dependency.TemplateFile{
		{
			Name:       filenames[0],
			ExportName: "template1_result.go",
			Overwrite:  true,
		},
		{
			Name:       filenames[1],
			ExportName: "template2_result.go",
			Overwrite:  true,
		},
		{
			Name:       filenames[2],
			ExportName: "",
			Overwrite:  false,
		},
	}

	// new template
	outputPath := filepath.Join(dir, "out")
	ts, err := NewTemplate(dir, tmpfiles, outputPath)
	assert.NoError(err)
	assert.NotNil(ts)
	assert.Len(ts.Template.Templates(), 3)
	assert.Len(ts.ExportConfigs, 2)
}

func TestScaffold_setType(t *testing.T) {
	require := require.New(t)
	dir, err := ioutil.TempDir("", "")
	require.NoError(err)

	tests := []struct {
		dataType   string
		columnType string
		nullable   bool
		want       string
	}{
		{dataType: "char", want: "string"},
		{dataType: "char", nullable: true, want: "null.String"},
		{dataType: "varchar", want: "string"},
		{dataType: "varchar", nullable: true, want: "null.String"},
		{dataType: "enum", want: "string"},
		{dataType: "enum", nullable: true, want: "null.String"},
		{dataType: "set", want: "string"},
		{dataType: "set", nullable: true, want: "null.String"},
		{dataType: "tinyint", nullable: true, want: "null.Int"},
		{dataType: "tinyint", columnType: "unsigned", want: "uint8"},
		{dataType: "tinyint", nullable: false, want: "int8"},
		{dataType: "smallint", nullable: true, want: "null.Int"},
		{dataType: "smallint", columnType: "unsigned", want: "uint16"},
		{dataType: "smallint", nullable: false, want: "int16"},
		{dataType: "mediumint", nullable: true, want: "null.Int"},
		{dataType: "mediumint", columnType: "unsigned", want: "uint32"},
		{dataType: "mediumint", nullable: false, want: "int32"},
		{dataType: "bigint", nullable: true, want: "null.Int"},
		{dataType: "bigint", columnType: "unsigned", want: "uint64"},
		{dataType: "bigint", nullable: false, want: "int64"},
		{dataType: "int", nullable: true, want: "null.Int"},
		{dataType: "int", columnType: "unsigned", want: "uint32"},
		{dataType: "int", nullable: false, want: "int32"},
		{dataType: "integer", nullable: true, want: "null.Int"},
		{dataType: "integer", columnType: "unsigned", want: "uint32"},
		{dataType: "integer", nullable: false, want: "int32"},
		{dataType: "float", nullable: true, want: "null.Float"},
		{dataType: "float", nullable: false, want: "float32"},
		{dataType: "double", nullable: true, want: "null.Float"},
		{dataType: "double", nullable: false, want: "float64"},
		{dataType: "decimal", nullable: true, want: "null.Float"},
		{dataType: "decimal", nullable: false, want: "float64"},
		{dataType: "dec", nullable: true, want: "null.Float"},
		{dataType: "dec", nullable: false, want: "float64"},
		{dataType: "date", nullable: true, want: "null.Time"},
		{dataType: "date", nullable: false, want: "time.Time"},
		{dataType: "datetime", nullable: true, want: "null.Time"},
		{dataType: "datetime", nullable: false, want: "time.Time"},
		{dataType: "timestamp", nullable: true, want: "null.Time"},
		{dataType: "timestamp", nullable: false, want: "time.Time"},
		{dataType: "time", nullable: true, want: "null.Time"},
		{dataType: "time", nullable: false, want: "time.Time"},
		{dataType: "dummy", nullable: false, want: "interface{}"},
	}
	for _, test := range tests {
		null := "not_null"
		if test.nullable {
			null = "null"
		}
		title := fmt.Sprintf("%s_%s_%s", test.dataType, test.columnType, null)
		t.Run(title, func(t *testing.T) {
			assert := assert.New(t)
			mc := mysql.Column{DataType: test.dataType, IsNullable: test.nullable, ColumnType: test.columnType}
			v := &TemplateDataColumn{Column: mc}
			v.setType()
			assert.Equal(test.want, v.Type)
		})
	}
	// generate template file
	filenames := make([]string, 3)
	for i := range filenames {
		filenames[i] = fmt.Sprintf("template_%05d.tmp", i)
		_, err := helper.CreateFile(filepath.Join(dir, filenames[i]), fmt.Sprintf("template%05d", i))
		require.NoError(err)
	}
}
