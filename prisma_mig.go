package main

import (
	"embed"
	"fmt"
	"github.com/charmbracelet/log"
	"os"
	"text/template"
)

//go:embed templates/*
var templatesFS embed.FS

func createPrismaMigration(input usefulData) error {
	log.Infof("Creating prisma migration...")
	tpl := template.Must(template.ParseFS(templatesFS, "templates/*"))

	schemaFile, err := os.Create(fmt.Sprintf("prisma/schema/%s.prisma", input.ModelSnakeCase))
	if err != nil {
		return err
	}

	if err := tpl.ExecuteTemplate(schemaFile, "prisma_model", input); err != nil {
		return err
	}
	return nil
}
