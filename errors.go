package main

import (
	"fmt"
)

var (
	notInAValidProjectStructureError = func(dir string) error {
		return fmt.Errorf(`you need use me in a valid project structure which contents '%s' directory`, dir)
	}
	errPackageAlreadyExists = func(dir, pkg string) error {
		return fmt.Errorf("dir %s already exists and package for %s model seems to be already created", dir, pkg)
	}
)
