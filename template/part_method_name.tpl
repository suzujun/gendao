{{.Name}}({{range $i, $p := .Params}}{{if ne $i 0}}, {{end}}{{$p.NameByCamelcase}} {{$p.Type}}{{end}}) ({{if .ReturnMany}}model.{{.ReturnModel}}Slice{{else}}*model.{{.ReturnModel}}{{end}}, error)