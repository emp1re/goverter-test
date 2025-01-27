// Code generated by github.com/jmattheis/goverter, DO NOT EDIT.
//go:build !goverter

package generated

import usezerovalueonpointerinconsistency "github.com/emp1re/goverter-test/example/use-zero-value-on-pointer-inconsistency"

type ConverterImpl struct{}

func (c *ConverterImpl) Convert(source usezerovalueonpointerinconsistency.Input) usezerovalueonpointerinconsistency.Output {
	var exampleOutput usezerovalueonpointerinconsistency.Output
	var xstring string
	if source.Name != nil {
		xstring = *source.Name
	}
	exampleOutput.Name = xstring
	exampleOutput.Age = source.Age
	return exampleOutput
}
