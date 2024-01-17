package gonic

import (
	"bytes"
	"embed"
	"html/template"
	"io"
	"reflect"
)

var (
	// tplFS is the embedded filesystem containing template files. You can write mulpitle files/dirs to tplFS.
	//go:embed tpl/*
	tplFS         embed.FS
	templateFuncs = template.FuncMap{"structToMap": structToMap}
)

// RangeStructer takes the first argument, which must be a struct, and
// returns the value of each field in a slice. It will return nil
// if there are no arguments or first argument is not a struct
/*
Usage:
{{ $data := . | structToMap }}
{{ range $key, $value := $data }}
    <tr>
        <th scope="col">{{ $key }}</th><td> {{ $value }}</td>
    </tr>
{{ end }}
*/
func structToMap(item interface{}) map[string]interface{} {
    out := make(map[string]interface{})
    v := reflect.ValueOf(item)

    for i := 0; i < v.NumField(); i++ {
        key := v.Type().Field(i).Name
        value := v.Field(i).Interface()
        out[key] = value
    }
    return out
}

func render(tpl string, data map[string]interface{}) (out []byte, err error) {
	fh, err := tplFS.Open(tpl)
	if err != nil {
		return
	}
	s, err := io.ReadAll(fh)
	if err != nil {
		return
	}

	t := template.Must(
		template.New("").Funcs(templateFuncs).Parse(string(s)),
		// template.New("").Parse(string(s)),
	)
	// builder := &strings.Builder{}
	builder := &bytes.Buffer{}
	if err = t.Execute(builder, data); err != nil {
		return
	}
	out = builder.Bytes()
	return out, nil
}
