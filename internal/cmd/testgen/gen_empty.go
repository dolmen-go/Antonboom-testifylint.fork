package main

import "text/template"

type EmptyCasesGenerator struct {
}

func (g EmptyCasesGenerator) Template() *template.Template {
	return emptyCasesTmpl
}

func (g EmptyCasesGenerator) Data() any {
	type test struct {
		InvalidChecks []Check
		ValidChecks   []Check
	}

	return struct {
		Pkgs     []string
		VarSets  [][]any
		Empty    test
		NotEmpty test
	}{
		Pkgs: []string{"assert", "require"},
		VarSets: [][]any{
			{"a"}, {"aPtr"}, {"s"}, {"m"}, {"ss"}, {"c"},
		},
		Empty: test{
			InvalidChecks: []Check{
				{Fn: "Len", ArgsTmpl: "t, %s, 0", ReportedMsg: "use %s.Empty"},
				{Fn: "Equal", ArgsTmpl: "t, len(%s), 0", ReportedMsg: "use %s.Empty"},
				{Fn: "Equal", ArgsTmpl: "t, 0, len(%s)", ReportedMsg: "use %s.Empty"},
				{Fn: "True", ArgsTmpl: "t, len(%s) == 0", ReportedMsg: "use %s.Empty"},
				{Fn: "True", ArgsTmpl: "t, 0 == len(%s)", ReportedMsg: "use %s.Empty"},
			},
			ValidChecks: []Check{
				{Fn: "Empty", ArgsTmpl: "t, %s"},
			},
		},
		NotEmpty: test{
			InvalidChecks: []Check{
				{Fn: "NotEqual", ArgsTmpl: "t, len(%s), 0", ReportedMsg: "use %s.NotEmpty"},
				{Fn: "NotEqual", ArgsTmpl: "t, 0, len(%s)", ReportedMsg: "use %s.NotEmpty"},
				{Fn: "Greater", ArgsTmpl: "t, len(%s), 0", ReportedMsg: "use %s.NotEmpty"},
				{Fn: "GreaterOrEqual", ArgsTmpl: "t, len(%s), 1", ReportedMsg: "use %s.NotEmpty"},
				{Fn: "True", ArgsTmpl: "t, len(%s) != 0", ReportedMsg: "use %s.NotEmpty"},
				{Fn: "True", ArgsTmpl: "t, 0 != len(%s)", ReportedMsg: "use %s.NotEmpty"},
				{Fn: "True", ArgsTmpl: "t, len(%s) > 0", ReportedMsg: "use %s.NotEmpty"},
				{Fn: "True", ArgsTmpl: "t, 0 < len(%s)", ReportedMsg: "use %s.NotEmpty"},
				{Fn: "True", ArgsTmpl: "t, len(%s) >= 1", ReportedMsg: "use %s.NotEmpty"},
				{Fn: "True", ArgsTmpl: "t, 1 <= len(%s)", ReportedMsg: "use %s.NotEmpty"},
			},
			ValidChecks: []Check{
				{Fn: "NotEmpty", ArgsTmpl: "t, %s"},
			},
		},
	}
}

var emptyCasesTmpl = template.Must(template.New("emptyCasesTmpl").
	Funcs(template.FuncMap{
		"ExpandCheck": ExpandCheck,
	}).
	Parse(`// Code generated by testifylint/internal/cmd/testgen. DO NOT EDIT.

package basic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmptyAsserts(t *testing.T) {
	var (
		a    [0]int
		aPtr *[0]int
		s    []int
		m    map[int]int
		ss   string
		c    chan int
	)
	{{ range $pi, $pkg := .Pkgs }}
	t.Run("{{ $pkg }}", func(t *testing.T) {
		{{- range $vi, $vars := $.VarSets }}
		{
			{{- range $ci, $check := $.Empty.InvalidChecks }}
			{{ ExpandCheck $check $pkg $vars }}
			{{ end }}}
		{{ end }}
		// Valid {{ $pkg }}s.
		{{ range $vi, $vars := $.VarSets }}
		{
			{{- range $ci, $check := $.Empty.ValidChecks }}
			{{ ExpandCheck $check $pkg $vars }}
			{{ end }}}
		{{ end -}}
	})
	{{ end }}}

func TestNotEmptyAsserts(t *testing.T) {
	var (
		a    [0]int
		aPtr *[0]int
		s    []int
		m    map[int]int
		ss   string
		c    chan int
	)
	{{ range $pi, $pkg := .Pkgs }}
	t.Run("{{ $pkg }}", func(t *testing.T) {
		{{- range $vi, $vars := $.VarSets }}
		{
		{{- range $ci, $check := $.NotEmpty.InvalidChecks }}
			{{ ExpandCheck $check $pkg $vars }}
			{{ end }}}
		{{ end }}
		// Valid {{ $pkg }}s.
		{{ range $vi, $vars := $.VarSets }}
		{
			{{- range $ci, $check := $.NotEmpty.ValidChecks }}
			{{ ExpandCheck $check $pkg $vars }}
			{{ end }}}
		{{ end -}}
	})
	{{ end }}}
`))
