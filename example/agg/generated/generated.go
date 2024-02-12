package generated

import agg "github.com/emp1re/goverter-test/example/agg"

type ConverterImpl struct{}

func (c *ConverterImpl) Convert(source agg.Input1) agg.Output1 {
	var mainOutput1 agg.Output1
	mainOutput1.ID = source.ID
	mainOutput1.Name = source.Name
	mainOutput1.Root = source.Root
	mainOutput1.Admin = source.Admin
	mainOutput1.Color = source.Color
	mainOutput1.Tax = source.Tax
	mainOutput1.Phone = source.Phone
	stringList := []string{source.Address}
	mainOutput1.Address = stringList
	return mainOutput1
}
func (c *ConverterImpl) ConvertAgg(source []agg.Input1) []agg.Output1 {
	var result []agg.Output1
	m := map[string]agg.Output1{}
	for _, v := range source {
		if _, ok := m[v.Address]; !ok {
			m[v.Address] = agg.Output1{
				ID:      v.ID,
				Name:    v.Name,
				Root:    v.Root,
				Admin:   v.Admin,
				Color:   v.Color,
				Tax:     v.Tax,
				Phone:   v.Phone,
				Address: []string{},
			}
		}
		if v.Address == "" {
			continue
		}
		obj := m[v.Address]
		obj.Address = append(obj.Address, v.Address)
		m[v.Address] = obj
	}
	for _, v := range source {
		found, ok := m[v.Address]
		if !ok {
			continue
		}
		delete(m, v.Address)
		result = append(result, found)
	}
	return result
}