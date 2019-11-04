package cmd

import (
  "fmt"
  "os"

  "github.com/spf13/cobra"

  "github.com/lakesite/reel/pkg/manager"
)

var rootCmd = &cobra.Command{
  Use:   "reel",
  Short: "reel development.",
  Long: `restore app state, for demos and development`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("use -h for help.")
  },
}

var versionCmd = &cobra.Command{
  Use:   "version",
  Short: "Print the version number of reel",
  Long:  `A number greater than 0, with prefix 'v', and possible suffixes like
          'a', 'b' or 'RELEASE'`,
  Run: func(cmd *cobra.Command, args []string) {
    // todo ...
    fmt.Println("reel v0.1a")
  },
}

var managerCmd = &cobra.Command{
  Use:   "manager",
  Short: "Run the manager.",
  Long:  `Run the management interface.`,
  Run: func(cmd *cobra.Command, args []string) {
    manager.RunManagementService()
  },
}

func init() {
  rootCmd.AddCommand(versionCmd)
  rootCmd.AddCommand(managerCmd)
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}
