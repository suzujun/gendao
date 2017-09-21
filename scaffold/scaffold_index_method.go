package scaffold

import (
	"fmt"
	"strings"

	"github.com/suzujun/gendao/helper"
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
		last := i == max-1
		unique := last && tIndex.Unique
		params := convCustomMethodParams(columns[:i+1], true)
		var orders []TemplateDataColumn
		if !last {
			orders = columns[i+1:]
		}
		// -------------------
		// Single designation
		// -------------------
		if last {
			// e.g. FindByID(id) return one
			methods = append(methods, genCustomMethod(params, nil, nil, modelName, unique, false))
		} else {
			// e.g. FindByIdOrderByIDAsc(id, limit) return many
			// e.g. FindByIdOrderByIDDesc(id, limit) return many
			methods = append(methods, genCustomMethod(params, nil, orders, modelName, false, false)) // ASC
			methods = append(methods, genCustomMethod(params, nil, orders, modelName, false, true))  // DESC
		}
		// -------------------
		// Multiple designation
		// -------------------
		params = convCustomMethodParams(columns[:i], true)
		rangeParam := newCustomMethodParam(column.Name, column.Type, false)
		orders = columns[i:]
		if last {
			// e.g. FindByIds(ids...) return many
			methods = append(methods, genCustomMethod(params, &rangeParam, nil, modelName, false, false))
		}
		// e.g. FindByIdsOrderByIDAsc(limit, rangeFunc...) return many
		// e.g. FindByIdsOrderByIDDesc(limit, rangeFunc...) return many
		methods = append(methods, genCustomMethod(params, &rangeParam, orders, modelName, false, false)) // ASC
		methods = append(methods, genCustomMethod(params, &rangeParam, orders, modelName, false, true))  // DESC
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
	wc := helper.NewWordConverter(name)
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

func genCustomMethod(params CustomMethodParams, rangeParam *CustomMethodParam, orders []TemplateDataColumn, modelName string, unique, desc bool) *CustomMethod {
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
			rangeFncs := newCustomMethodParam(rangeParam.Name+"_range_fncs", "..."+typ, false)
			method.Params = append(method.Params, rangeFncs)
		}
		method.Orders = convCustomMethodParams(orders, false)
	}
	if len(orders) == 0 && rangeParam != nil {
		name := helper.NewWordConverter(rangeParam.Name).Pluralize().ToString()
		typ := fmt.Sprintf("[]%s", rangeParam.Type)
		rp := newCustomMethodParam(name, typ, true)
		rp.Name = rangeParam.Name // set original name
		method.Params = append(method.Params, rp)
		method.RangeParam = nil
	}
	method.ReturnModel = modelName
	method.Unique = unique
	method.Desc = desc
	method.ReturnMany = !unique
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
		names = append(names, helper.NewWordConverter(cm.RangeParam.Name).Pascalcase().Pluralize().Lint().ToString())
	}
	names = []string{"FindBy", strings.Join(names, "And")}
	if len(cm.Orders) > 0 {
		ascDesc := "Asc"
		if cm.Desc {
			ascDesc = "Desc"
		}
		names = append(names, "OrderBy", cm.Orders.joinName("And"), ascDesc)
	}
	cm.Name = strings.Join(names, "")
}

func (cmps CustomMethodParams) joinName(sep string) string {
	res := make([]string, len(cmps))
	for i, cmp := range cmps {
		if i == len(cmps)-1 {
			res[i] = helper.NewWordConverter(cmp.Name).Pascalcase().Lint().ToString()
		} else {
			res[i] = cmp.NameByPascalcase
		}
	}
	return strings.Join(res, sep)
}

func getRangeFncType(typ string) string {
	typ = strings.ToLower(typ)
	if typ == "interface{}" {
		return "" // none range
	} else if typ == "bool" {
		return "" // none range
	} else if strings.Contains(typ, "string") {
		return "ranger.RangeStrFnc"
	} else if strings.Contains(typ, "int") {
		return "ranger.RangeIntFnc"
	} else if strings.Contains(typ, "float") {
		return "ranger.RangeFloatFnc"
	} else if strings.Contains(typ, "time") {
		return "ranger.RangeTimeFnc"
	} else {
		panic(fmt.Sprintf("unknown range func type [%s]", typ))
	}
}
