package main

import (
	jen "github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"os"
)

func genDeclAt19() jen.Code {
	return jen.Null()
}
func genDeclAt129() jen.Code {
	return jen.Null().Type().Id("Advisor").Struct(jen.Id("CreateAdvisorInput"), jen.Id("Id").Id("uuid").Dot("UUID"), jen.Id("CreatedAt").Qual("time", "Time"), jen.Id("UpdatedAt").Qual("time", "Time"), jen.Id("LatestWeight").Id("goption").Dot("Optional").Index(jen.Id("udecimal").Dot("Decimal")), jen.Id("LatestHeight").Id("goption").Dot("Optional").Index(jen.Id("int")), jen.Id("CurrentAge").Id("int"))
}
func genFile() *jen.File {
	ret := jen.NewFile("advisors")
	ret.Add(genDeclAt19())
	ret.Add(genDeclAt129())
	return ret
}

type EntityField struct {
	Name       string
	Type       string
	IsOptional bool
}

type CreateEntityInput struct {
	EntityName string
	IsUuid     bool
	Fields     []EntityField
}

func createEntity(input CreateEntityInput) jen.Code {
	entityName := strcase.ToCamel(input.EntityName)
	jen.Type().Id(entityName).Struct()
}

func main() {
	file := genFile()
	err := file.Render(os.Stdout)
	if err != nil {
		panic(err)
	}
}
