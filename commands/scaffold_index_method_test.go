package commands

import (
	"fmt"
	"strings"
	"testing"
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
	indexs := []TemplateDataIndex{
		TemplateDataIndex{
			Columns: []TemplateDataColumn{
				genTemplateDataColumn("id", "string"),
			},
			Unique:  true,
			Primary: true,
		},
		TemplateDataIndex{
			Columns: []TemplateDataColumn{
				genTemplateDataColumn("channel_id", "int"),
			},
			Unique:  false,
			Primary: false,
		},
		TemplateDataIndex{
			Columns: []TemplateDataColumn{
				genTemplateDataColumn("stream_name", "string"),
				genTemplateDataColumn("edition", "string"),
				genTemplateDataColumn("sequence", "int"),
			},
			Unique:  true,
			Primary: true,
		},
	}
	for _, tIndex := range indexs {
		methods := GenCustomMethods(tIndex, "Thumbnail")
		for i, method := range methods {
			params := make([]string, len(method.Params))
			for i, p := range method.Params {
				params[i] = fmt.Sprintf("%s %s", p.NameByCamelcase, p.Type)
			}
			fmt.Println(i, fmt.Sprintf("%s(%s)", method.Name, strings.Join(params, ", ")))
		}
		println("-------------")
	}
}
