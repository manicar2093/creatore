package main

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

const banner = `
   ____   ______      _____     ____     ________     ____     ______      _____  
  / ___) (   __ \    / ___/    (    )   (___  ___)   / __ \   (   __ \    / ___/  
 / /      ) (__) )  ( (__      / /\ \       ) )     / /  \ \   ) (__) )  ( (__    
( (      (    __/    ) __)    ( (__) )     ( (     ( ()  () ) (    __/    ) __)   
( (       ) \ \  _  ( (        )    (       ) )    ( ()  () )  ) \ \  _  ( (      
 \ \___  ( ( \ \_))  \ \___   /  /\  \     ( (      \ \__/ /  ( ( \ \_))  \ \___  
  \____)  )_) \__/    \____\ /__(  )__\    /__\      \____/    )_) \__/    \____\

All code creation for model, repository and controller in one place
`

var (
	rootCmd = &cobra.Command{
		Short: "A new way to create code faster",
		Long:  banner,
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Shows installed version",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info(Version)
		},
	}

	genCmd = &cobra.Command{
		Use:     "gen [model_name] [model_fields...]",
		Short:   "Create model, repository and controller",
		Example: "creatore gen User moduleName:string:optional age:int --binary-id",
		Long: `
Create model, repository and REST API controller with given data.

You can select from all golang types to create your model and set it as optional in case:

	moduleName:string:optional

This creates an optional field called moduleName of type string. If optional is omitted data is taken as required.

	createdAt:time

This creates a field called created_at of type time.Time



`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := trigger(normalizeArgsAsEntityInput(args, binaryId)); err != nil {
				log.Error(err.Error())
				return
			}

			log.Info("✅ Ready!")
		},
	}

	initCmd = &cobra.Command{
		Use:   "init [project_url]",
		Short: "Create a new project",
		Long: `Create a new project with needed structure:

<project_name>/
├── cmd
│   └── api
│       └── main.go
├── go.mod
└── internal
    └── domain
        └── models

<project_url> should be (by convention) a url which indicates where the code is served, but you can use any valid string

From this you can start using gen command.

`,
		Example: "creatore init github.com/great-dev/revolutionary",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var isForceConfirmed bool
			if isForced {
				if err := huh.NewConfirm().
					Title("CAUTION. Are you sure of this?").
					Description("This action cannot be undone and will override existing directory content").
					Affirmative("Yes!").
					Negative("No.").
					Value(&isForceConfirmed).
					Run(); err != nil {
					log.Fatal(err.Error())
				}
				if !isForceConfirmed {
					log.Info("No action was made")
					return
				}
			}

			if err := createNewProject(newProjectData{
				moduleName:        args[0],
				isForcedConfirmed: isForceConfirmed,
			}); err != nil {
				log.Error(err.Error())
				return
			}
			log.Info("✅ Ready!")
			log.Info(`🏁 Next steps:
➡️ Run 'go mod tidy' to install deps
➡️ Run 'creatore gen' command to create your first API resource
➡️ Enjoy! 😎
`)
		},
	}
)
