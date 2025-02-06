package main

import (
	"fmt"
	"github.com/charmbracelet/log"
	. "github.com/dave/jennifer/jen"
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/rjNemo/underscore"
)

type StructNames struct {
	// ModelSnakeCase contains model name as received in snake_case
	ModelSnakeCase string
	// RepositoryStructName contains struct name for repository following this format: <Entity>Repository
	RepositoryStructName string
	// RepositoryStructName contains var name for repository following this format: <entity>Repository
	RepositoryStructVarName string
	// RepositoryStructConstructorName contains how repository constructor method has to be named
	RepositoryStructConstructorName string
	// ModelServicePackageName contains entity package name following the format: <entity_given_by_user>
	ModelServicePackageName string
	// ModelStructName contains entity struct name following format: <EntityGivenByUser>
	ModelStructName string
	// GoModName contains in which code needs to be generated
	GoModName string
	// ControllerStructName contains the name controller struct will use
	ControllerStructName string
	// ControllerStructConstructorName contains the name for controller constructor
	ControllerStructConstructorName string
	// ReceiverVarName contains the name of receiver var name for struct methods
	ReceiverVarName                    string
	UpdateSelectiveByIdInputStructName string
}

type DirNames struct {
	// BaseModelDir is the directory creatore is looking to save the new model file domain/models/
	BaseModelDir string
	// InternalBaseDir is where creatore will create the new package for the new model internal/
	InternalBaseDir string
	// ServiceDir is the path of the new package internal/<model_name>/
	ServiceDir string
	// ModelFile is the path of the file where the model will be saved
	ModelFile string
	// ControllerFile is the path where controller will be saved
	ControllerFile string
	// RepositoryFile is the path where repository will be saved
	RepositoryFile string
}

type usefulData struct {
	StructNames
	DirNames
	// ModelsKey contains literal "models"
	ModelsKey string
	// Fields contains all Fields with useful data
	Fields []FieldMeta
	// IdTypeAsJenCode contains the id type as qualifier
	IdTypeAsJenCode *Statement
	// ModelStructQualifier contains the quealifier code for entity
	ModelStructQualifier *Statement
	// IsIdUUID indicates if id is type UUID. If false is considered id is an int
	IsIdUUID bool
}

func createUsefulData(input ModelCreationInput, goModName string) usefulData {
	log.Info("Creating data to generate code...")
	var (
		pc = pluralize.NewClient()

		entityNameAsPlural              = pc.Plural(input.EntityName)
		modelSnakeCase                  = strcase.ToSnake(input.EntityName)
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
		baseModelDir                    = "internal/domain/models"
	)

	return usefulData{
		StructNames: StructNames{
			ModelSnakeCase:                     modelSnakeCase,
			RepositoryStructName:               repositoryStructName,
			RepositoryStructConstructorName:    repositoryStructConstructorName,
			RepositoryStructVarName:            repositoryStructVarName,
			ModelServicePackageName:            modelServicePackageName,
			ModelStructName:                    modelStructName,
			GoModName:                          goModName,
			ControllerStructName:               controllerStructName,
			ControllerStructConstructorName:    controllerStructConstructorName,
			ReceiverVarName:                    "c",
			UpdateSelectiveByIdInputStructName: "UpdateSelectiveByIdInput",
		},
		DirNames: DirNames{
			BaseModelDir:    baseModelDir,
			InternalBaseDir: "internal/",
			ServiceDir:      fmt.Sprintf("internal/%s", modelServicePackageName),
			ModelFile:       fmt.Sprintf("internal/domain/models/%s_creatore.go", modelServicePackageName),
			ControllerFile:  fmt.Sprintf("internal/%s/controller_creatore.go", modelServicePackageName),
			RepositoryFile:  fmt.Sprintf("internal/%s/repository_creatore.go", modelServicePackageName),
		},
		ModelsKey:            modelsKey,
		Fields:               fieldsMeta,
		IdTypeAsJenCode:      idType,
		IsIdUUID:             idIsUUID,
		ModelStructQualifier: Qual(fmt.Sprintf("%s/%s", goModName, baseModelDir), modelStructName),
	}
}

type FieldMeta struct {
	// NameForTags contains Field name as snake_case
	NameForTags string
	// NameForStructAttribute contains Field name as CamelCase
	NameForStructAttribute string
	// contains Field name as lowerCamelCase
	NameForFunctionParams string
	// TypeAsJenCode contains Field type as *Statement useful to generate code with jennifer
	TypeAsJenCode *Statement
	Type          string
	// Field contains data as received
	Field ModelFieldData
	// Tags contains all need Tags for this Field
	Tags     map[string]string
	IsId     bool
	IsIdUuid bool
}

func normalizeEntityFieldsData(input ModelCreationInput) ([]FieldMeta, *Statement, bool) {
	idField, isIdUUID := createIdField(input)
	input.Fields = append([]ModelFieldData{idField}, input.Fields...)
	var idType *Statement

	return underscore.Map(input.Fields, func(f ModelFieldData) FieldMeta {
		nameForStructAttribute := strcase.ToCamel(f.Name)
		nameForFunctionParams := strcase.ToLowerCamel(f.Name)
		isId := f.Name == "Id"
		jenType, strType := getType(f)
		if isId {
			idType = jenType
		}

		return FieldMeta{
			NameForTags:            strcase.ToSnake(f.Name),
			NameForStructAttribute: nameForStructAttribute,
			NameForFunctionParams:  nameForFunctionParams,
			TypeAsJenCode:          jenType,
			Type:                   strType,
			Field:                  f,
			Tags:                   getFieldTags(f, isId),
			IsId:                   isId,
			IsIdUuid:               isIdUUID,
		}
	}), idType, isIdUUID
}

func createIdField(input ModelCreationInput) (ModelFieldData, bool) {
	if input.IsIdUuid {
		return ModelFieldData{
			Name: idKey,
			Type: "uuid",
		}, true
	}
	return ModelFieldData{
		Name: idKey,
		Type: "int",
	}, false
}

func getType(f ModelFieldData) (*Statement, string) {
	switch f.Type {
	case "decimal":
		return Qual("github.com/quagmt/udecimal", "Decimal"), "Decimal"
	case "uuid":
		return Qual("github.com/google/uuid", "UUID"), "String"
	case "time":
		return Qual("time", "Time"), "DateTime"
	default:
		return Id(f.Type), strcase.ToCamel(f.Type)
	}
}

func getFieldTags(field ModelFieldData, isId bool) map[string]string {
	tags := make(map[string]string)
	tags["json"] = strcase.ToSnake(field.Name)
	tags["mapstructure"] = strcase.ToSnake(field.Name)
	if isId {
		tags["param"] = strcase.ToSnake(field.Name)
	}
	if !field.IsOptional {
		tags["validate"] = "required"
	}
	switch {
	case field.Name == idKey && field.Type == "uuid":
		tags["gorm"] = "default:gen_random_uuid()"
	case field.Name == idKey && field.Type == "uint":
		tags["gorm"] = "primaryKey"
	}

	return tags
}
