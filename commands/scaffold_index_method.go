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
			// e.g. FindByID(id) return one
			methods = append(methods, genCustomMethod(params, nil, nil, modelName, tIndex.Unique, false, false))
		} else {
			// e.g. FindByIdOrderByIDAsc(id, limit) return many
			// e.g. FindByIdOrderByIDDesc(id, limit) return many
			methods = append(methods, genCustomMethod(params, nil, orders, modelName, tIndex.Unique, true, false)) // ASC
			methods = append(methods, genCustomMethod(params, nil, orders, modelName, tIndex.Unique, true, true))  // DESC
		}
		params = convCustomMethodParams(columns[:i], true)
		rangeParam := newCustomMethodParam(column.Name, column.Type, false)
		// e.g. FindByIds(ids...) return many
		methods = append(methods, genCustomMethod(params, &rangeParam, nil, modelName, tIndex.Unique, true, false))
		// e.g. FindByIdsOrderByIDAsc(limit, rangeFunc...) return many
		// e.g. FindByIdsOrderByIDDesc(limit, rangeFunc...) return many
		methods = append(methods, genCustomMethod(params, &rangeParam, orders, modelName, tIndex.Unique, true, false)) // ASC
		methods = append(methods, genCustomMethod(params, &rangeParam, orders, modelName, tIndex.Unique, true, true))  // DESC
	}
	res := make([]CustomMethod, 0, len(methods))
	for _, m := range methods {
		if m == nil {
			continue
		}
		res = append(res, *m)
	}
	return res
}

func newCustomMethodParam(name, typ string, where bool) CustomMethodParam {
	wc := NewWordConverter(name)
	return CustomMethodParam{
		Name:             name,
		NameByCamelcase:  wc.Camelcase().Lint().ToString(),
		NameByPascalcase: wc.Pascalcase().Lint().ToString(),
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
		limit := newCustomMethodParam("limit", "uint64", false)
		method.Params = append(method.Params, limit)
		if rangeParam != nil {
			typ := getRangeFncType(rangeParam.Type)
			if typ == "" {
				return nil
			}
			rangeFncs := newCustomMethodParam("range_fncs", "..."+typ, false)
			method.Params = append(method.Params, rangeFncs)
		}
		method.Orders = convCustomMethodParams(orders, false)
	}
	if len(orders) == 0 && rangeParam != nil {
		name := NewWordConverter(rangeParam.Name).Pluralize().ToString()
		typ := fmt.Sprintf("[]%s", rangeParam.Type)
		rp := newCustomMethodParam(name, typ, true)
		rp.Name = rangeParam.Name // set original name
		method.Params = append(method.Params, rp)
		method.RangeParam = nil
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
		names = append(names, NewWordConverter(cm.RangeParam.Name).Pascalcase().Pluralize().Lint().ToString())
	}
	names = []string{"FindBy", strings.Join(names, "And")}
	if len(cm.Orders) > 0 {
		ascDesc := "Asc"
		if cm.Desc {
			ascDesc = "Desc"
		}
		names = append(names, "OrderBy", cm.Orders.joinName(""), ascDesc)
	}
	cm.Name = strings.Join(names, "")
}

func (cmps CustomMethodParams) joinName(sep string) string {
	res := make([]string, len(cmps))
	for i, cmp := range cmps {
		if i == len(cmps)-1 {
			res[i] = NewWordConverter(cmp.Name).Pascalcase().Lint().ToString()
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
		return "ranger.RangeIntFnc"
	case "uint8":
		return "ranger.RangeIntFnc"
	case "uint16":
		return "ranger.RangeIntFnc"
	case "uint32":
		return "ranger.RangeIntFnc"
	case "uint64":
		return "ranger.RangeIntFnc"
	case "time.Time":
		return "ranger.RangeTimeFnc"
	case "interface{}":
		return "" // none range
	case "bool":
		return "" // none range
	default:
		panic(fmt.Sprintf("unknown range func type [%s]", typ))
	}
}
