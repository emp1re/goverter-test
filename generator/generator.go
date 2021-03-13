package generator

import (
	"fmt"
	"go/importer"
	"go/types"
	"sort"

	"github.com/dave/jennifer/jen"
	"github.com/jmattheis/go-genconv/builder"
	"github.com/jmattheis/go-genconv/comments"
)

type Config struct {
	Name string
}

func Generate(pattern string, mapping []comments.Converter, config Config) (*jen.File, error) {
	sources, err := importer.For("source", nil).Import(pattern)
	if err != nil {
		return nil, err
	}
	file := jen.NewFile(config.Name)
	file.HeaderComment("// Code generated by github.com/jmattheis/go-genconv, DO NOT EDIT.")

	for _, converter := range mapping {
		obj := sources.Scope().Lookup(converter.Name)
		if obj == nil {
			return nil, fmt.Errorf("%s: could not find %s", pattern, converter.Name)
		}

		// create the converter struct
		file.Type().Id(converter.Config.Name).Struct()

		gen := generator{
			file:   file,
			name:   converter.Config.Name,
			lookup: map[string]map[string]GenerateMethod{},
			rules: []builder.Builder{
				&builder.BasicTargetPointerRule{},
				&builder.Pointer{},
				&builder.TargetPointer{},
				&builder.Basic{},
				&builder.Struct{},
				&builder.List{},
				&builder.Map{},
			},
		}

		// we checked in comments, that it is an interface
		interf := obj.Type().Underlying().(*types.Interface)
		for i := 0; i < interf.NumMethods(); i++ {
			method := interf.Method(i)
			if err := gen.registerMethod(method); err != nil {
				return nil, fmt.Errorf(`Error while creating converter method: %s

%s`, method.FullName(), err)
			}

		}
		if err := gen.createMethods(); err != nil {
			return nil, err
		}
	}
	return file, nil
}

type GenerateMethod struct {
	ID     string
	Name   string
	Source types.Type
	Target types.Type
}

type generator struct {
	name   string
	file   *jen.File
	rules  []builder.Builder
	lookup map[string]map[string]GenerateMethod
	// source type to target type string
}

func (g *generator) registerMethod(method *types.Func) error {
	signature, ok := method.Type().(*types.Signature)
	if !ok {
		return fmt.Errorf("expected signature %#v", method.Type())
	}
	params := signature.Params()
	if params.Len() != 1 {
		return fmt.Errorf("expected signature to have only one parameter")
	}
	result := signature.Results()
	if result.Len() != 1 {
		return fmt.Errorf("expected signature to have only one parameter")
	}
	source := params.At(0).Type()
	target := result.At(0).Type()

	inner, _ := g.lookup[source.String()]
	if inner == nil {
		inner = map[string]GenerateMethod{}
		g.lookup[source.String()] = inner
	}
	inner[target.String()] = GenerateMethod{
		ID:     method.FullName(),
		Name:   method.Name(),
		Source: source,
		Target: target,
	}
	return nil
}
func (g *generator) createMethods() error {
	methods := []GenerateMethod{}
	for _, inner := range g.lookup {
		for _, method := range inner {
			methods = append(methods, method)
		}
	}
	sort.Slice(methods, func(i, j int) bool {
		return methods[i].Name < methods[j].Name
	})
	for _, method := range methods {
		err := g.addMethod(method.Name, method.Source, method.Target)
		if err != nil {
			return fmt.Errorf("Error while creating converter method: %s\n\n%s", method.ID, builder.ToString(err))
		}
	}
	return nil
}

func (g *generator) addMethod(name string, source, target types.Type) *builder.Error {
	sourceID := jen.Id("source")
	ruleSource := builder.TypeOf(source)
	ruleTarget := builder.TypeOf(target)

	stmt, newID, err := g.Build(&builder.MethodContext{Namer: builder.NewNamer()}, sourceID.Clone(), builder.TypeOf(source), builder.TypeOf(target))
	if err != nil {
		return err.Lift(&builder.Path{
			SourceID:   "source",
			TargetID:   "target",
			SourceType: source.String(),
			TargetType: target.String(),
		})
	}

	stmt = append(stmt, jen.Return().Add(newID))
	g.file.Func().Params(jen.Id("c").Op("*").Id(g.name)).Id(name).
		Params(jen.Id("source").Add(ruleSource.TypeAsJen())).Params(ruleTarget.TypeAsJen()).
		Block(stmt...)

	return nil
}

func (g *generator) Build(ctx *builder.MethodContext, sourceID builder.JenID, source, target *builder.Type) ([]jen.Code, builder.JenID, *builder.Error) {
	for _, rule := range g.rules {
		if rule.Matches(source, target) {
			return rule.Build(g, ctx, sourceID, source, target)
		}
	}

	return nil, nil, builder.NewError(fmt.Sprintf("TypeMismatch: Cannot convert %s to %s", source.T, target.T))
}
