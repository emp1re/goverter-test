// Code generated by github.com/jmattheis/goverter, DO NOT EDIT.
//go:build !goverter

package b

type RootBImpl struct{}

func (c *RootBImpl) Convert(source []string) []string {
	var stringList []string
	if source != nil {
		stringList = make([]string, len(source))
		for i := 0; i < len(source); i++ {
			stringList[i] = source[i]
		}
	}
	return stringList
}