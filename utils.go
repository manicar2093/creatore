package main

import (
	"fmt"
	"github.com/charmbracelet/log"
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
	// goModName contains in which code needs to be generated
	goModName string
	// controllerStructName contains the name controller struct will use
	controllerStructName string
	// controllerStructConstructorName contains the name for controller constructor
	controllerStructConstructorName string
	// receiverVarName contains the name of receiver var name for struct methods
	receiverVarName                    string
	updateSelectiveByIdInputStructName string
}

type dirNames struct {
	// baseModelDir is the directory creatore is looking to save the new model file domain/models/
	baseModelDir string
	// internalBaseDir is where creatore will create the new package for the new model internal/
	internalBaseDir string
	// serviceDir is the path of the new package internal/<model_name>/
	serviceDir string
	// modelFile is the path of the file where the model will be saved
	modelFile string
	// controllerFile is the path where controller will be saved
	controllerFile string
	// repositoryFile is the path where repository will be saved
	repositoryFile string
}

type usefulData struct {
	structNames
	dirNames
	// modelsKey contains literal "models"
	modelsKey string
	// fields contains all fields with useful data
	fields []fieldMeta
	// idTypeAsJenCode contains the id type as qualifier
	idTypeAsJenCode *Statement
	// modelStructQualifier contains the quealifier code for entity
	modelStructQualifier *Statement
	// isIdUUID indicates if id is type UUID. If false is considered id is an int
	isIdUUID bool
}

func createUsefulData(input ModelCreationInput, goModName string) usefulData {
	log.Info("Creating data to generate code...")
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
		fieldsMeta, idType, idIsUUID    = normalizeEntityFieldsData(input)
	)

	return usefulData{
		structNames: structNames{
			repositoryStructName:               repositoryStructName,
			repositoryStructConstructorName:    repositoryStructConstructorName,
			repositoryStructVarName:            repositoryStructVarName,
			modelServicePackageName:            modelServicePackageName,
			modelStructName:                    modelStructName,
			goModName:                          goModName,
			controllerStructName:               controllerStructName,
			controllerStructConstructorName:    controllerStructConstructorName,
			receiverVarName:                    "c",
			updateSelectiveByIdInputStructName: "UpdateSelectiveByIdInput",
		},
		dirNames: dirNames{
			baseModelDir:    "internal/domain/models",
			internalBaseDir: "internal/",
			serviceDir:      fmt.Sprintf("internal/%s", modelServicePackageName),
			modelFile:       fmt.Sprintf("internal/domain/models/%s_creatore.go", modelServicePackageName),
			controllerFile:  fmt.Sprintf("internal/%s/controller_creatore.go", modelServicePackageName),
			repositoryFile:  fmt.Sprintf("internal/%s/repository_creatore.go", modelServicePackageName),
		},
		modelsKey:            modelsKey,
		fields:               fieldsMeta,
		idTypeAsJenCode:      idType,
		isIdUUID:             idIsUUID,
		modelStructQualifier: Qual(fmt.Sprintf("%s/%s", goModName, modelsKey), modelStructName),
	}
}

type fieldMeta struct {
	nameForTags            string
	nameForStructAttribute string
	nameForFunctionParams  string
	typeAsJenCode          *Statement
	field                  ModelFieldData
	tags                   map[string]string
}

func normalizeEntityFieldsData(input ModelCreationInput) ([]fieldMeta, *Statement, bool) {
	idField, isUUID := createIdField(input)
	input.Fields = append([]ModelFieldData{idField}, input.Fields...)
	var idType *Statement

	return underscore.Map(input.Fields, func(f ModelFieldData) fieldMeta {
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
	}), idType, isUUID
}

func createIdField(input ModelCreationInput) (ModelFieldData, bool) {
	if input.IsUuid {
		return ModelFieldData{
			Name: idKey,
			Type: "uuid",
		}, true
	}
	return ModelFieldData{
		Name: idKey,
		Type: "uint",
	}, false
}

func getType(f ModelFieldData) *Statement {
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

func getFieldTags(field ModelFieldData) map[string]string {
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
