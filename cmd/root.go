package cmd

import (
  "fmt"
  "os"

  "github.com/spf13/cobra"

  "github.com/lakesite/reel/pkg/manager"
)

var (
  config string
  application string

  rootCmd = &cobra.Command{
    Use:   "reel -c [config.toml] -a [application name]",
    Short: "run reel with a config against an app.",
    Long: `restore app state, for demos and development`,
    Run: func(cmd *cobra.Command, args []string) {
      ms := &manager.ManagerService{}
      ms.Init(config)
      ms.InitApp(application)
      ms.Rewind(application)
    },
  }

  versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Print the version number of reel",
    Long:  `A number greater than 0, with prefix 'v', and possible suffixes like
            'a', 'b' or 'RELEASE'`,
    Run: func(cmd *cobra.Command, args []string) {
      // todo ...
      fmt.Println("reel v0.1a")
    },
  }

  managerCmd = &cobra.Command{
    Use:   "manager",
    Short: "Run the manager.",
    Long:  `Run the management interface.`,
    Run: func(cmd *cobra.Command, args []string) {
      ms := &manager.ManagerService{}
      ms.Init(config)
      ms.RunManagementService()
    },
  }
)

func init() {
  rootCmd.Flags().StringVarP(&config, "config", "c", "", "config file")
  rootCmd.Flags().StringVarP(&application, "application", "a", "", "application name")
  rootCmd.MarkFlagRequired("config")
  rootCmd.MarkFlagRequired("application")

  managerCmd.Flags().StringVarP(&config, "config", "c", "", "config file")
  managerCmd.MarkFlagRequired("config")

  rootCmd.AddCommand(versionCmd)
  rootCmd.AddCommand(managerCmd)
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}
