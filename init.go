package main

import (
	"errors"
	"fmt"
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

	jf := NewFile("main")
	jf.ImportAlias(echoQual, "echo")

	//
	//e.GET("/", func(c echo.Context) error {
	//	return c.String(http.StatusOK, "Hello, World!")
	//})

	jf.Func().Id("main").Params().Block(
		Id("e").Op(":=").Qual(echoQual, "New").Call(),
		Id("e").Dot("Get").Call(
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
