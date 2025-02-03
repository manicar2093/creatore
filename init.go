package main

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/log"
	. "github.com/dave/jennifer/jen"
	"net/url"
	"os"
	"strings"
)

type newProjectData struct {
	moduleName        string
	isForcedConfirmed bool
}

func createNewProject(input newProjectData) (string, error) {
	log.Info("Creating project structure...")
	asUrl, err := url.Parse(input.moduleName)
	if err != nil {
		return "", errors.Join(err, fmt.Errorf("%s is not a valid module moduleName. follow convention of use urls as github.com/<user>/<repo>, gitlab.com/<user>/<repo>, bitbucket.com/<user>/<repo>, etc", input.moduleName))
	}

	pathSlice := strings.Split(asUrl.Path, string(os.PathSeparator))
	projectDirName := pathSlice[len(pathSlice)-1]

	if err := os.Mkdir(projectDirName, os.ModePerm); err != nil {
		if !os.IsExist(err) {
			return "", err
		}
		if !input.isForcedConfirmed {

			return "", fmt.Errorf("%s already exists, use -f to force creation", projectDirName)
		}
	}

	for _, dir := range []string{
		fmt.Sprintf("%s/internal/domain/models", projectDirName),
		fmt.Sprintf("%s/cmd/api", projectDirName),
		fmt.Sprintf("%s/prisma/schema", projectDirName),
	} {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return "", err
		}
	}

	if err := os.WriteFile(
		fmt.Sprintf("%s/go.mod", projectDirName),
		[]byte(fmt.Sprintf("module %s\n", input.moduleName)),
		0755,
	); err != nil {
		return "", err
	}

	log.Info("Creating prisma version: 6.3.0 migration sys...")
	var prismaContent = `// This is your Prisma schema file,
// learn more about it in the docs: https://pris.ly/d/prisma-schema

generator client {
  provider        = "prisma-client-js"
  previewFeatures = ["prismaSchemaFolder"]
}

datasource db {
  provider = "postgresql"
  url      = env("DATABASE_URL")
}`
	if err := os.WriteFile(fmt.Sprintf("%s/prisma/schema/schema.prisma", projectDirName), []byte(prismaContent), 0755); err != nil {
		return "", err
	}

	var packageJsonContent = `{
  "devDependencies": {
    "prisma": "^6.3.0"
  },
  "dependencies": {
    "@prisma/client": "^6.3.0"
  }
}`
	if err := os.WriteFile(fmt.Sprintf("%s/package.json", projectDirName), []byte(packageJsonContent), 0755); err != nil {
		return "", err
	}

	log.Info("Creating env files...")
	envs := []string{
		"dev",
		"test",
	}
	envContent := func(env, project string) string {
		return fmt.Sprintf(`DATABASE_URL="postgresql://development:development@localhost:5432/%s_%s?sslmode=disable"
ENVIRONMENT=%s
PORT=3000`, project, env, env)

	}
	for _, env := range envs {
		if err := os.WriteFile(fmt.Sprintf("%s/.env.%s", projectDirName, env), []byte(envContent(env, projectDirName)), 0755); err != nil {
			return "", err
		}
	}

	jf := NewFile("main")
	jf.ImportAlias(echoQual, "echo")

	jf.Func().Id("main").Params().Block(
		Id("e").Op(":=").Qual(echoQual, "New").Call(),
		Id("e").Dot("GET").Call(
			Lit("/"),
			Func().Params(echoContextParam()).Error().Block(
				Return(
					Id(ctxKeyword).Dot("String").Call(
						Qual(netHttpQual, "StatusCreated"),
						Lit("Hello from creatore"),
					),
				),
			),
		),
		Id("e").Dot("Logger").Dot("Fatal").Call(
			Id("e").Dot("Start").Call(Lit(":3000")),
		),
	)

	return projectDirName, jf.Save(fmt.Sprintf("%s/cmd/api/main.go", projectDirName))
}
