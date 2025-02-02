package main

import (
	"errors"
	"github.com/charmbracelet/log"
	"github.com/rjNemo/underscore"
	"golang.org/x/mod/modfile"
	"os"
	"strings"
	"time"
)

var binaryId bool

func main() {
	rootCmd.AddCommand(genCmd)
	rootCmd.AddCommand(versionCmd)
	genCmd.Flags().BoolVarP(&binaryId, "binary-id", "", false, "indicates if model has a UUID as id. If omitted Id will be an int")

	configLogger()
	rootCmd.Execute()
}

func normalizeArgsAsEntityInput(args []string, isBinaryId bool) ModelCreationInput {
	return ModelCreationInput{
		EntityName: args[0],
		IsUuid:     isBinaryId,
		Fields: underscore.Map(args[1:len(args)], func(item string) ModelFieldData {
			var (
				splitted = strings.Split(item, ":")
				name     = splitted[0]
				typ      = splitted[1]
				isOpt    = isOptionalFromArgs(args)
			)

			return ModelFieldData{
				Name:       name,
				Type:       typ,
				IsOptional: isOpt,
			}

		}),
	}
}

func isOptionalFromArgs(args []string) bool {
	if len(args) < 3 {
		return false
	}

	return args[2] == "optional"
}

func trigger(input ModelCreationInput) error {
	mod, err := getGoModName()
	if err != nil {
		return err
	}

	data := createUsefulData(input, mod)

	handlers := []func(usefulData) error{
		validatesHasNeededDirs,
		createNewDirectories,
		createModelFile,
		createRepositoryFile,
		createControllerFile,
	}

	for _, handler := range handlers {
		if err := handler(data); err != nil {
			return err
		}
	}

	return nil
}

func validatesHasNeededDirs(input usefulData) error {
	log.Info("Validating needed directories to work...")
	if _, err := os.Open(input.dirNames.internalBaseDir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return notInAValidProjectStructureError(input.dirNames.internalBaseDir)
		}
	}

	if _, err := os.Open(input.dirNames.baseModelDir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return notInAValidProjectStructureError(input.dirNames.baseModelDir)
		}
	}

	return nil
}

func createNewDirectories(input usefulData) error {
	log.Infof("Creating new '%s' package", input.dirNames.serviceDir)

	if err := os.Mkdir(input.serviceDir, os.ModePerm); err != nil {
		if errors.Is(err, os.ErrExist) {
			return errPackageAlreadyExists(input.dirNames.serviceDir, input.modelServicePackageName)
		}
	}

	return nil
}

func getGoModName() (string, error) {
	goModBytes, err := os.ReadFile("go.mod")
	if err != nil {
		if os.IsNotExist(err) {
			return "", errors.New("go.mod not found")
		}
		return "", err
	}

	modName := modfile.ModulePath(goModBytes)
	log.Infof("Go module detected: %s", modName)

	return modName, nil
}

func configLogger() {
	log.SetTimeFormat(time.Kitchen)
	log.SetPrefix("🔨")
}
