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
	repositoryStructName string
	// repositoryStructConstructorName contains how repository constructor method has to be named
	repositoryStructConstructorName string
	// entityServicePackageName contains entity package name following the format: <entity_given_by_user>
	entityServicePackageName string
	// entityStructName contains entity struct name following format: <EntityGivenByUser>
	entityStructName string
	// entitiesKey contains literal "entities"
	entitiesKey string
	// fields contains all fields with useful data
	fields []fieldMeta
	// moduleName contains in which code needs to be generated
	moduleName string
	// idTypeAsJenCode contains the id type as qualifier
	idTypeAsJenCode *Statement
	// entityStructQualifier contains the quealifier code for entity
	entityStructQualifier *Statement
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
		fieldsMeta, idType              = normalizeEntityFieldsData(input)
		moduleName                      = "github.com/user/package"
	)

	return usefulData{
		repositoryStructName:            repositoryStructName,
		repositoryStructConstructorName: repositoryStructConstructorName,
		entityServicePackageName:        entityServicePackageName,
		entityStructName:                entityStructName,
		entitiesKey:                     entitiesKey,
		fields:                          fieldsMeta,
		moduleName:                      moduleName, // TODO: This must be configurable
		idTypeAsJenCode:                 idType,
		entityStructQualifier:           Qual(fmt.Sprintf("%s/%s", moduleName, entitiesKey), entityStructName),
	}
}

type fieldMeta struct {
	nameForTags            string
	nameForStructAttribute string
	nameForFunctionParams  string
	typeAsJenCode          *Statement
	field                  EntityField
	tags                   map[string]string
}

func normalizeEntityFieldsData(input CreateEntityInput) ([]fieldMeta, *Statement) {
	idField := createIdField(input)
	input.Fields = append([]EntityField{idField}, input.Fields...)
	var idType *Statement

	return underscore.Map(input.Fields, func(f EntityField) fieldMeta {
		nameForStructAttribute := strcase.ToCamel(f.Name)
		nameForFunctionParams := strcase.ToLowerCamel(f.Name)
		if f.Name == "Id" {
			idType = getType(idField, nameForStructAttribute)
		}

		return fieldMeta{
			nameForTags:            strcase.ToSnake(f.Name),
			nameForStructAttribute: nameForStructAttribute,
			nameForFunctionParams:  nameForFunctionParams,
			typeAsJenCode:          getType(f, nameForStructAttribute),
			field:                  f,
			tags:                   getFieldTags(f),
		}
	}), idType
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
	switch f.Type {
	case "decimal":
		return Qual("github.com/quagmt/udecimal", "Decimal")
	case "uuid":
		return Qual("github.com/google/uuid", "UUID")
	case "time":
		return Qual("time", "Time")
	default:
		return Id(f.Type)
	}
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
