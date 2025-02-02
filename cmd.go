package main

import (
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
		Use:   "gen --binary-id [name] [args...]",
		Short: "Create model, repository and controller",
		Long: `
Create model, repository and controller with given data.

creatore gen --binary-id User name:string:optional

If optional is omitted data is taken as required
`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := trigger(normalizeArgsAsEntityInput(args, binaryId)); err != nil {
				return err
			}

			log.Info("✅ Ready!")
			return nil
		},
	}
)
