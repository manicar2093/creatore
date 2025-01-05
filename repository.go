package main

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"github.com/rjNemo/underscore"

	"os"
)

const (
	saveMethodKey                = "Save"
	getByIdMethodKey             = "GetById"
	getAllMethodKey              = "GetAll"
	updateSelectiveByIdMethodKey = "UpdateSelectiveById"
	deleteByIdMethodKey          = "DeleteById"
)

type methodGenerators func(input usefulData, fnTargetKeyword *Statement) []Code

var (
	RepoSupportedMethods = []string{
		saveMethodKey, getByIdMethodKey, getAllMethodKey, updateSelectiveByIdMethodKey, deleteByIdMethodKey,
	}
	repoMethodsGeneratorsMap = map[string]methodGenerators{
		saveMethodKey: func(input usefulData, fnTargetKeyword *Statement) []Code {

			return []Code{
				Commentf("%s can Create and Update an entity. You can use this for http PATH method. Check https://gorm.io/docs/update.html#Save-All-Fields for more info", saveMethodKey).Line(),
				fnTargetKeyword.Id(saveMethodKey).Params(
					Id("input").Op("*").Qual(fmt.Sprintf("%s/%s", input.moduleName, input.entitiesKey), input.entityStructName),
				).Error().
					Block(
						If(
							Id("res").Op(":=").Id("c").Dot("db").Dot("Save").Call(
								Id("input"),
							),
							Id("res").Dot("Error").Op("!=").Nil(),
						).Block(
							Return(Id("res").Dot("Error")),
						).Line().
							Return(
								Nil(),
							),
					).Line().Line(),
			}
		},
		getByIdMethodKey: func(input usefulData, fnTargetKeyword *Statement) []Code {
			return []Code{fnTargetKeyword.Id(getByIdMethodKey).Params(
				Id("id").Add(input.idTypeAsJenCode),
			).Params(
				Op("*").Qual(fmt.Sprintf("%s/%s", input.moduleName, input.entitiesKey), input.entityStructName),
				Error(),
			).Block(
				Var().Id("found").Qual(fmt.Sprintf("%s/%s", input.moduleName, input.entitiesKey), input.entityStructName),
				If(
					Id("res").Op(":=").Id("c").Dot("db").Dot("First").Call(
						Op("&").Id("found"),
						Id("id"),
					),
					Id("res").Dot("Error").Op("!=").Nil(),
				).Block(
					Return(Nil(), Id("res").Dot("Error")),
				),
				Return(Op("&").Id("found"), Nil()),
			).Line().Line()}
		},
		getAllMethodKey: func(input usefulData, fnTargetKeyword *Statement) []Code {
			return []Code{fnTargetKeyword.Id(getAllMethodKey).Params().Add(
				Index().Qual(fmt.Sprintf("%s/%s", input.moduleName, input.entitiesKey), input.entityStructName),
			).Block(
				Var().Id("found").Index().Qual(fmt.Sprintf("%s/%s", input.moduleName, input.entitiesKey), input.entityStructName),
				If(
					Id("res").Op(":=").Id("c").Dot("db").Dot("Model").Call(
						Op("&").Qual(fmt.Sprintf("%s/%s", input.moduleName, input.entitiesKey), input.entityStructName).Block(),
					).Dot("Find").Call(Op("&").Id("found")),
					Id("res").Dot("Error").Op("!=").Nil(),
				).Block(
					Return(Nil(), Id("res").Dot("Error")),
				),
				Return(Op("&").Id("found"), Nil()),
			).Line().Line()}
		},
		updateSelectiveByIdMethodKey: func(input usefulData, fnTargetKeyword *Statement) []Code {
			updateInputStructName := fmt.Sprintf("Update%sInput", input.entityStructName)

			strct := Null().Type().Id(updateInputStructName).Struct(
				underscore.Map(input.fields, func(meta fieldMeta) Code {
					if meta.field.Name == "Id" {
						return nil
					}
					return Null().Id(meta.nameForStructAttribute).Qual(optionalQual, optionalName).Index(meta.typeAsJenCode).Tag(meta.tags)
				})...,
			)

			mapWithUpdateData := Id("updates").Op("=").New(Map(Id("string")).Id("interface{}"))
			resultVar := Id("result").Op("=").Qual(fmt.Sprintf("%s/%s", input.moduleName, input.entitiesKey), input.entityStructName).Block()

			varDeclarations := Var().Defs(resultVar, mapWithUpdateData)

			optionValidations := underscore.Map(input.fields, func(meta fieldMeta) Code {
				if meta.field.Name == "Id" {
					return nil
				}
				return If(Id("changes").Dot(meta.nameForStructAttribute).Dot("IsPresent").Call()).Block(
					Id("updates").
						Index(Lit(meta.nameForTags)).
						Op("=").
						Id("changes").
						Dot(meta.nameForStructAttribute).
						Dot("MustGet").
						Call(),
				)
			})

			updateBody := []Code{
				varDeclarations,
				Line(),
			}
			updateBody = append(updateBody, optionValidations...)
			updateGormCall := If(
				Id("res").Op(":=").Id("c").Dot("db").Dot("Model").Call(Id(fmt.Sprintf("&%s", "result"))).Dot("Clauses").Call(Id("clause.Returning{}")).Dot("Where").Call(Lit("id = ?"), Id("id")).Dot("Updates").Call(Id("updates")),
				Id("res").Dot("Error").Op("!=").Nil(),
			).Block(
				Return(Nil(), Id("res").Dot("Error")),
			)
			updateBody = append(updateBody,
				Line(),
				updateGormCall,
				Return(Id("result"), Nil()),
			)

			return []Code{
				strct,
				Line(),
				Line(),
				Commentf("%s can select which field has to be updated from given input", updateSelectiveByIdMethodKey),
				Line(),
				fnTargetKeyword.Id(updateSelectiveByIdMethodKey).Params(
					Id("id").Add(input.idTypeAsJenCode),
					Id("changes").Id(updateInputStructName),
				).Params(
					Op("*").Qual(fmt.Sprintf("%s/%s", input.moduleName, input.entitiesKey), input.entityStructName), Error(),
				).Block(
					updateBody...,
				).Line().Line()}
		},
		deleteByIdMethodKey: func(input usefulData, fnTargetKeyword *Statement) []Code {
			return []Code{fnTargetKeyword.Id(deleteByIdMethodKey).Params(
				Id("input").Id("string"),
			).Params(
				Id("string"), Error(),
			).Block(
				Comment("TODO: implement me!"),
				Id("panic").Call(Lit("implement me!")),
			).Line().Line()}
		},
	}
)

func createRepositoryFile(data usefulData) error {
	jf := NewFile(data.entityServicePackageName)
	jf.PackageComment("Code generated by creatore")

	jf.Add(generateRepositoryCode(data)...)

	return jf.Render(os.Stdout)
}

func generateRepositoryCode(data usefulData) []Code {
	var result []Code
	result = append(result, Null().Type().Id(data.repositoryStructName).Struct(
		Id("db").Op("*").Qual("github.com/manicar2093/winter/connections", "connections"),
	).Line().Line())
	result = append(result, Null().Func().Id(data.repositoryStructConstructorName).Params(
		Id("db").Op("*").Qual("github.com/manicar2093/winter/connections", "connections"),
	).Block(
		Return(
			Id(fmt.Sprintf("&%s", data.repositoryStructName)).Values(
				Dict{
					Id("db"): Id("db"),
				},
			),
		),
	).Line().Line())

	underscore.Each(RepoSupportedMethods, func(s string) {
		funcContext := Func().Params(
			Id("c").Op("*").Id(data.repositoryStructName),
		)
		result = append(result, repoMethodsGeneratorsMap[s](data, funcContext)...)
	})

	return result

}
