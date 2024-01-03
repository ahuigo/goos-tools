package gonic

import (
	"bytes"
	"embed"
	"html/template"
	"io"
)

var (
	// tplFS is the embedded filesystem containing template files. You can write mulpitle files/dirs to tplFS.
	//go:embed tpl/* 
	tplFS embed.FS

)

func render(tpl string, data map[string]interface{})(out []byte, err error){
	fh, err := tplFS.Open(tpl)
	if err != nil {
		return
	}
	s, err := io.ReadAll(fh)
	if err != nil {
		return
	}

	t := template.Must(template.New("").Parse(string(s)))
	// builder := &strings.Builder{}
	builder := &bytes.Buffer{}
	if err = t.Execute(builder, data); err != nil {
		return
	}
	out = builder.Bytes()
	return out, nil	
}