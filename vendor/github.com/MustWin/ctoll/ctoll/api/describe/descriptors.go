package describe

import (
	"regexp"

	"github.com/MustWin/ctoll/ctoll/api/errcode"
)

type RouteDescriptor struct {
	Name        string
	Path        string
	Entity      string
	Description string
	Methods     []MethodDescriptor
}

type MethodDescriptor struct {
	Method      string
	Description string
	Requests    []RequestDescriptor
}

type RequestDescriptor struct {
	Name            string
	Description     string
	Headers         []ParameterDescriptor
	PathParameters  []ParameterDescriptor
	QueryParameters []ParameterDescriptor
	Body            BodyDescriptor
	Successes       []ResponseDescriptor
	Failures        []ResponseDescriptor
}

type ResponseDescriptor struct {
	Name        string
	Description string
	StatusCode  int
	Headers     []ParameterDescriptor
	Fields      []ParameterDescriptor
	ErrorCodes  []errcode.ErrorCode
	Body        BodyDescriptor
}

type BodyDescriptor struct {
	ContentType string
	Format      string
}

type ParameterDescriptor struct {
	Name        string
	Type        string
	Description string
	Required    bool
	Format      string
	Regexp      *regexp.Regexp
	Examples    []string
}
