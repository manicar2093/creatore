package main

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
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
	typeAsJenCode          *Statement
	field                  EntityField
	isId                   bool
	tags                   map[string]string
}

func normalizeEntityFieldsData(input CreateEntityInput) []fieldMeta {
	input.Fields = append([]EntityField{createIdField(input)}, input.Fields...)

	return underscore.Map(input.Fields, func(f EntityField) fieldMeta {
		nameForStructAttribute := strcase.ToCamel(f.Name)
		return fieldMeta{
			nameForTags:            strcase.ToSnake(f.Name),
			nameForStructAttribute: nameForStructAttribute,
			typeAsJenCode:          getType(f, nameForStructAttribute),
			field:                  f,
			isId:                   f.Name == "id",
			tags:                   getFieldTags(f),
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

func getType(f EntityField, nameForStructAttribute string) *Statement {
	builder := Null().Id(nameForStructAttribute)

	var getTypeQual *Statement

	switch f.Type {
	case "decimal":
		getTypeQual = Qual("github.com/quagmt/udecimal", "Decimal")
	case "uuid":
		getTypeQual = Qual("github.com/google/uuid", "UUID")
	case "time":
		getTypeQual = Qual("time", "Time")
	default:
		getTypeQual = Id(f.Type)
	}

	if f.IsOptional {
		return builder.Qual("github.com/manicar2093/goption", "goption").Index(getTypeQual)
	}

	return builder.Add(getTypeQual)
}

func getFieldTags(field EntityField) map[string]string {
	tags := make(map[string]string)
	tags["json"] = strcase.ToSnake(field.Name)
	tags["mapstructure"] = strcase.ToSnake(field.Name)
	switch {
	case field.Name == idKey && field.Type == "uuid":
		tags["gorm"] = "default:gen_random_uuid()"
	case field.Name == idKey && field.Type == "uint":
		tags["gorm"] = "primaryKey"
	}

	return tags
}
