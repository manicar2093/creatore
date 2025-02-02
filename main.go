package main

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/rjNemo/underscore"
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
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func normalizeArgsAsEntityInput(args []string, isBinaryId bool) ModelCreationInput {
	return ModelCreationInput{
		EntityName: args[0],
		IsUuid:     isBinaryId,
		Fields: underscore.Map(args[1:len(args)-1], func(item string) ModelFieldData {
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
	data := createUsefulData(input)
	if err := createNewDirectories(data); err != nil {
		return err
	}
	if err := createModelFile(data); err != nil {
		return err
	}
	if err := createRepositoryFile(data); err != nil {
		return err
	}
	if err := createControllerFile(data); err != nil {
		return err
	}
	return nil
}

func createNewDirectories(input usefulData) error {
	log.Infof("Creating '%s' directory", input.modelServicePackageName)
	return os.MkdirAll(fmt.Sprintf("./internal/%s", input.modelServicePackageName), os.ModePerm)
}

func configLogger() {
	log.SetTimeFormat(time.Kitchen)
	log.SetPrefix("🔨")
}
