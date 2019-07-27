package internal

import (
	"encoding/json"
	"io"
	"text/template"

	"github.com/k0kubun/pp"
)

type Output interface {
	Print(v interface{}) error
	io.Writer
}

type prettyOutput struct {
	io.Writer
	pp *pp.PrettyPrinter
}

func (o *prettyOutput) Print(v interface{}) error {
	_, err := o.pp.Println(v)
	return err
}

func newPrettyOutput(w io.Writer) *prettyOutput {
	p := &prettyOutput{
		Writer: w,
		pp:     pp.New(),
	}

	p.pp.SetOutput(w)

	return p
}

type jsonOutput struct {
	io.Writer
}

func newJSONOutput(w io.Writer) *jsonOutput {
	return &jsonOutput{Writer: w}
}

func (o *jsonOutput) Print(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	_, err = o.Write(data)
	_, err = o.Write([]byte("\n"))

	return err
}

type templateOutput struct {
	io.Writer
	tmpl *template.Template
}

func newTemplateOutput(tmpl string, w io.Writer) *templateOutput {
	return &templateOutput{
		Writer: w,
		tmpl:   template.Must(template.New("_").Parse(tmpl + "\n")),
	}
}

func (o *templateOutput) Print(v interface{}) error {
	return o.tmpl.Execute(o, v)
}

func NewOutput(w io.Writer, format string) Output {
	switch format {
	case "pretty":
		return newPrettyOutput(w)
	case "json":
		return newJSONOutput(w)
	default:
		return newTemplateOutput(format, w)
	}
}
