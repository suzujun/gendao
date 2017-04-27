package commands

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScaffoldIndexMethod_GenCustomMethods(t *testing.T) {
	genTemplateDataColumn := func(name, typ string) TemplateDataColumn {
		return TemplateDataColumn{
			Name:             name,
			NameByCamelcase:  NewWordConverter(name).Camelcase().Lint().ToString(),
			NameByPascalcase: NewWordConverter(name).Pascalcase().Lint().ToString(),
			Type:             typ,
		}
	}
	tests := []struct {
		title string
		index TemplateDataIndex
	}{
		{
			title: "uniq:index{id}",
			index: TemplateDataIndex{
				Columns: []TemplateDataColumn{
					genTemplateDataColumn("id", "string"),
				},
				Unique:  true,
				Primary: true,
			},
		},
		{
			title: "uniq:index{channel_id}",
			index: TemplateDataIndex{
				Columns: []TemplateDataColumn{
					genTemplateDataColumn("channel_id", "int"),
				},
				Unique:  true,
				Primary: false,
			},
		},
		{
			title: "no-uniq:index{channel_id}",
			index: TemplateDataIndex{
				Columns: []TemplateDataColumn{
					genTemplateDataColumn("channel_id", "int"),
				},
				Unique:  false,
				Primary: false,
			},
		},
		{
			title: "uniq:index{stream_name,edition,sequence}",
			index: TemplateDataIndex{
				Columns: []TemplateDataColumn{
					genTemplateDataColumn("stream_name", "string"),
					genTemplateDataColumn("edition", "string"),
					genTemplateDataColumn("sequence", "int"),
				},
				Unique:  true,
				Primary: true,
			},
		},
		{
			title: "uniq:index{program_id,start_questioning_at}",
			index: TemplateDataIndex{
				Columns: []TemplateDataColumn{
					genTemplateDataColumn("program_id", "string"),
					genTemplateDataColumn("start_questioning_at", "null.Time"),
				},
				Unique:  true,
				Primary: true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			methods := GenCustomMethods(test.index, "Thumbnail")
			for i, method := range methods {
				params := make([]string, len(method.Params))
				for i, p := range method.Params {
					params[i] = fmt.Sprintf("%s %s", p.NameByCamelcase, p.Type)
				}
				var result string
				if method.ReturnMany {
					result = fmt.Sprintf("model.%sSlice, error", method.ReturnModel)
				} else {
					result = fmt.Sprintf("*model.%s, error", method.ReturnModel)
				}
				fmt.Println(i, fmt.Sprintf("%s(%s) (%s)", method.Name, strings.Join(params, ", "), result))
			}
			fmt.Println()
		})
	}
}

func TestScaffoldIndexMethod_getRangeFncType(t *testing.T) {

	tests := []struct {
		typ    string
		result string
	}{
		{typ: "string", result: "ranger.RangeStrFnc"},
		{typ: "null.String", result: "ranger.RangeStrFnc"},
		{typ: "int", result: "ranger.RangeIntFnc"},
		{typ: "int8", result: "ranger.RangeIntFnc"},
		{typ: "int16", result: "ranger.RangeIntFnc"},
		{typ: "int32", result: "ranger.RangeIntFnc"},
		{typ: "int64", result: "ranger.RangeIntFnc"},
		{typ: "uint", result: "ranger.RangeIntFnc"},
		{typ: "uint8", result: "ranger.RangeIntFnc"},
		{typ: "uint16", result: "ranger.RangeIntFnc"},
		{typ: "uint32", result: "ranger.RangeIntFnc"},
		{typ: "uint64", result: "ranger.RangeIntFnc"},
		{typ: "null.Int", result: "ranger.RangeIntFnc"},
		{typ: "float32", result: "ranger.RangeFloatFnc"},
		{typ: "float64", result: "ranger.RangeFloatFnc"},
		{typ: "null.Float", result: "ranger.RangeFloatFnc"},
		{typ: "time.Time", result: "ranger.RangeTimeFnc"},
		{typ: "null.Time", result: "ranger.RangeTimeFnc"},
		{typ: "interface{}", result: ""},
		{typ: "bool", result: ""},
	}
	for _, test := range tests {
		t.Run(test.typ, func(t *testing.T) {
			assert := assert.New(t)
			assert.Equal(test.result, getRangeFncType(test.typ))
		})
	}
}
