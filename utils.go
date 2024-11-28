package main

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/rjNemo/underscore"
)

type usefulData struct {
	// repositoryStructName contains struct name for repository following this format: <Entity>Repository
	repositoryStructName            string
	repositoryStructConstructorName string
	// entityServicePackageName contains entity package name following the format: <entity_given_by_user>
	entityServicePackageName string
	// entityStructName contains entity struct name following format: <EntityGivenByUser>
	entityStructName string
	// entitiesKey contains literal "entities"
	entitiesKey string
	// fields contains all fields with useful data
	fields []fieldMeta
}

func createUsefulData(input CreateEntityInput) usefulData {
	var (
		pc = pluralize.NewClient()

		entityNameAsPlural              = pc.Plural(input.EntityName)
		repositoryStructName            = fmt.Sprintf("%sRepository", strcase.ToCamel(entityNameAsPlural))
		repositoryStructConstructorName = fmt.Sprintf("New%s", repositoryStructName)
		entityServicePackageName        = strcase.ToSnake(entityNameAsPlural)
		entitiesKey                     = "entities"
		entityStructName                = strcase.ToCamel(input.EntityName)
	)

	return usefulData{
		repositoryStructName:            repositoryStructName,
		repositoryStructConstructorName: repositoryStructConstructorName,
		entityServicePackageName:        entityServicePackageName,
		entityStructName:                entityStructName,
		entitiesKey:                     entitiesKey,
		fields:                          normalizeEntityFieldsData(input),
	}
}

type fieldMeta struct {
	nameForTags            string
	nameForStructAttribute string
	typeAsJenCode          jen.Code
	field                  EntityField
	isId                   bool
}

func normalizeEntityFieldsData(input CreateEntityInput) []fieldMeta {
	input.Fields = append([]EntityField{createIdField(input)}, input.Fields...)

	return underscore.Map(input.Fields, func(f EntityField) fieldMeta {
		return fieldMeta{
			nameForTags:            strcase.ToSnake(f.Name),
			nameForStructAttribute: strcase.ToCamel(f.Name),
			typeAsJenCode:          getType(f.Type),
			field:                  f,
			isId:                   f.Name == "id",
		}
	})
}

func createIdField(input CreateEntityInput) EntityField {
	if input.IsUuid {
		return EntityField{
			Name: idKey,
			Type: "uuid",
		}
	}
	return EntityField{
		Name: idKey,
		Type: "uint",
	}
}

func getType(t string) jen.Code {
	switch t {
	case "decimal":
		return jen.Qual("github.com/quagmt/udecimal", "Decimal")
	case "uuid":
		return jen.Qual("github.com/google/uuid", "UUID")
	case "time":
		return jen.Qual("time", "Time")
	default:
		return jen.Id(t)
	}
}
