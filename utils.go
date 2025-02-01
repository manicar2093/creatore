package main

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/rjNemo/underscore"
)

type structNames struct {
	// repositoryStructName contains struct name for repository following this format: <Entity>Repository
	repositoryStructName string
	// repositoryStructName contains var name for repository following this format: <entity>Repository
	repositoryStructVarName string
	// repositoryStructConstructorName contains how repository constructor method has to be named
	repositoryStructConstructorName string
	// modelServicePackageName contains entity package name following the format: <entity_given_by_user>
	modelServicePackageName string
	// modelStructName contains entity struct name following format: <EntityGivenByUser>
	modelStructName string
	// moduleName contains in which code needs to be generated
	moduleName string
	// controllerStructName contains the name controller struct will use
	controllerStructName string
	// controllerStructConstructorName contains the name for controller constructor
	controllerStructConstructorName string
	// receiverVarName contains the name of receiver var name for struct methods
	receiverVarName string
}

type reposMethodsNames struct {
	saveMethodKey                string
	getByIdMethodKey             string
	getAllPaginatedMethodKey     string
	updateSelectiveByIdMethodKey string
	deleteByIdMethodKey          string
}

type usefulData struct {
	structNames
	reposMethodsNames
	// modelsKey contains literal "entities"
	modelsKey string
	// fields contains all fields with useful data
	fields []fieldMeta
	// idTypeAsJenCode contains the id type as qualifier
	idTypeAsJenCode *Statement
	// modelStructQualifier contains the quealifier code for entity
	modelStructQualifier *Statement
}

func createUsefulData(input CreateEntityInput) usefulData {
	var (
		pc = pluralize.NewClient()

		entityNameAsPlural              = pc.Plural(input.EntityName)
		entityPluralNameAsCamelCase     = strcase.ToCamel(entityNameAsPlural)
		repositoryStructName            = fmt.Sprintf("%sRepository", entityPluralNameAsCamelCase)
		repositoryStructVarName         = fmt.Sprintf("%sRepository", strcase.ToLowerCamel(entityPluralNameAsCamelCase))
		controllerStructName            = fmt.Sprintf("%sController", entityPluralNameAsCamelCase)
		repositoryStructConstructorName = fmt.Sprintf("New%s", repositoryStructName)
		controllerStructConstructorName = fmt.Sprintf("New%s", controllerStructName)
		modelServicePackageName         = strcase.ToSnake(entityNameAsPlural)
		modelsKey                       = "models"
		modelStructName                 = strcase.ToCamel(input.EntityName)
		fieldsMeta, idType              = normalizeEntityFieldsData(input)
		moduleName                      = "github.com/user/package"
	)

	return usefulData{
		structNames: structNames{
			repositoryStructName:            repositoryStructName,
			repositoryStructConstructorName: repositoryStructConstructorName,
			repositoryStructVarName:         repositoryStructVarName,
			modelServicePackageName:         modelServicePackageName,
			modelStructName:                 modelStructName,
			moduleName:                      moduleName, // TODO: This must be configurable
			controllerStructName:            controllerStructName,
			controllerStructConstructorName: controllerStructConstructorName,
			receiverVarName:                 "c",
		},
		// FIXME: NOT USED BY NOW :/ I don't even know if it is useful
		reposMethodsNames: reposMethodsNames{
			saveMethodKey:                "Save",
			getByIdMethodKey:             "GetById",
			getAllPaginatedMethodKey:     "GetAllPaginated",
			updateSelectiveByIdMethodKey: "UpdateSelectiveById",
			deleteByIdMethodKey:          "DeleteById",
		},
		modelsKey:            modelsKey,
		fields:               fieldsMeta,
		idTypeAsJenCode:      idType,
		modelStructQualifier: Qual(fmt.Sprintf("%s/%s", moduleName, modelsKey), modelStructName),
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
			idType = getType(idField)
		}

		return fieldMeta{
			nameForTags:            strcase.ToSnake(f.Name),
			nameForStructAttribute: nameForStructAttribute,
			nameForFunctionParams:  nameForFunctionParams,
			typeAsJenCode:          getType(f),
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

func getType(f EntityField) *Statement {
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
