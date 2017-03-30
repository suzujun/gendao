package commands

import (
	"fmt"
	"strings"
)

type (
	// CustomMethod ...
	CustomMethod struct {
		Name        string
		Params      CustomMethodParams
		RangeParam  *CustomMethodParam
		Orders      CustomMethodParams
		Unique      bool
		ReturnMany  bool
		ReturnModel string
		Desc        bool
	}
	// CustomMethodParam ...
	CustomMethodParam struct {
		Name             string
		NameByCamelcase  string
		NameByPascalcase string
		Type             string
		Where            bool
	}
	// CustomMethodParams ...
	CustomMethodParams []CustomMethodParam
)

// GenCustomMethods ...
func GenCustomMethods(tIndex TemplateDataIndex, modelName string) []CustomMethod {
	max := len(tIndex.Columns)
	columns := tIndex.Columns
	var methods []*CustomMethod
	for i, column := range columns {
		params := convCustomMethodParams(columns[:i+1], true)
		orders := columns
		if i == max-1 && tIndex.Unique {
			methods = append(methods, genCustomMethod(params, nil, nil, modelName, tIndex.Unique, false, false))
		} else {
			methods = append(methods, genCustomMethod(params, nil, orders, modelName, tIndex.Unique, true, false)) // ASC
			methods = append(methods, genCustomMethod(params, nil, orders, modelName, tIndex.Unique, true, true))  // DESC
		}
		params = convCustomMethodParams(columns[:i], true)
		rangeParam := newCustomMethodParam(column.Name, column.Type, false)
		methods = append(methods, genCustomMethod(params, &rangeParam, orders, modelName, tIndex.Unique, true, false)) // ASC
		methods = append(methods, genCustomMethod(params, &rangeParam, orders, modelName, tIndex.Unique, true, true))  // DESC
	}
	var res []CustomMethod
	for _, m := range methods {
		if m == nil {
			continue
		}
		res = append(res, *m)
	}
	return res
}

func newCustomMethodParam(name, typ string, where bool) CustomMethodParam {
	return CustomMethodParam{
		Name:             name,
		NameByCamelcase:  ConvCamelcase(name, true),
		NameByPascalcase: ConvPascalcase(name, true),
		Type:             typ,
		Where:            where,
	}
}
func convCustomMethodParams(cols []TemplateDataColumn, where bool) CustomMethodParams {
	params := make(CustomMethodParams, len(cols))
	for i, col := range cols {
		params[i] = newCustomMethodParam(col.Name, col.Type, where)
	}
	return params
}

func genCustomMethod(params CustomMethodParams, rangeParam *CustomMethodParam, orders []TemplateDataColumn, modelName string, unique, returnMany, desc bool) *CustomMethod {
	method := CustomMethod{}
	method.Params = params
	if rangeParam != nil {
		method.RangeParam = rangeParam
	}
	if len(orders) > 0 {
		limit := CustomMethodParam{
			Name:             "limit",
			NameByPascalcase: "Pimit",
			NameByCamelcase:  "limit",
			Type:             "uint64",
		}
		method.Params = append(method.Params, limit)
		if rangeParam != nil {
			typ := getRangeFncType(rangeParam.Type)
			if typ == "" {
				return nil
			}
			rangeFncs := CustomMethodParam{
				Name:             "range_fncs",
				NameByPascalcase: "RangeFncs",
				NameByCamelcase:  "rangeFncs",
				Type:             "..." + typ,
			}
			method.Params = append(method.Params, rangeFncs)
		}
		method.Orders = convCustomMethodParams(orders, false)
	}
	method.ReturnModel = modelName
	method.Unique = unique
	method.Desc = desc
	method.ReturnMany = returnMany
	method.setName()
	return &method
}

func (cm *CustomMethod) setName() {
	var names []string
	for _, param := range cm.Params {
		if param.Where {
			names = append(names, param.NameByPascalcase)
		}
	}
	if cm.RangeParam != nil {
		names = append(names, ConvPascalcasePluralize(cm.RangeParam.Name, true))
	}
	names = []string{"FindBy", strings.Join(names, "And")}
	if len(cm.Orders) > 0 {
		if cm.Desc {
			names = append(names, "OrderBy", cm.Orders.joinName(""), "Desc")
		} else {
			names = append(names, "OrderBy", cm.Orders.joinName(""), "Asc")
		}
	}
	cm.Name = strings.Join(names, "")
}

func (cmps CustomMethodParams) joinName(sep string) string {
	res := make([]string, len(cmps))
	for i, cmp := range cmps {
		if i == len(cmps)-1 {
			res[i] = ConvPascalcase(cmp.Name, true)
		} else {
			res[i] = cmp.NameByPascalcase
		}
	}
	return strings.Join(res, sep)
}

func getRangeFncType(typ string) string {
	switch typ {
	case "string":
		return "ranger.RangeStrFnc"
	case "int":
		return "ranger.RangeIntFnc"
	case "int8":
		return "ranger.RangeIntFnc"
	case "int16":
		return "ranger.RangeIntFnc"
	case "int32":
		return "ranger.RangeIntFnc"
	case "int64":
		return "ranger.RangeIntFnc"
	case "uint":
		return "ranger.RangeUintFnc"
	case "uint8":
		return "ranger.RangeUintFnc"
	case "uint16":
		return "ranger.RangeUintFnc"
	case "uint32":
		return "ranger.RangeUintFnc"
	case "uint64":
		return "ranger.RangeUintFnc"
	case "time.Time":
		return "ranger.RangeTimeFnc"
	case "interface{}":
		return "ranger.RangeFnc"
	case "bool":
		return "" // none range
	default:
		panic(fmt.Sprintf("unknown range func type [%s]", typ))
	}
}
