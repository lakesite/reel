package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/lakesite/ls-governor"

	"github.com/lakesite/reel/pkg/api"
	"github.com/lakesite/reel/pkg/job"
	"github.com/lakesite/reel/pkg/reel"
)

var (
	config      string
	application string

	rootCmd = &cobra.Command{
		Use:   "reel -c [config.toml] -a [application name]",
		Short: "run reel with a config against an app.",
		Long:  `restore app state, for demos and development`,
		Run: func(cmd *cobra.Command, args []string) {
			gms := &governor.ManagerService{}
			gms.InitManager(config)
			gapi := gms.CreateAPI(application)
			reel.InitApp(application, gapi)
			reel.Rewind(application, "", gapi)
		},
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of reel",
		Long: `A number greater than 0, with prefix 'v', and possible suffixes like
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
			// setup the job queue
			rw := &job.ReelWorker{}
			rw.Start()

			gms := &governor.ManagerService{}
			gms.InitManager(config)
			gms.InitDatastore("reel")
			gapi := gms.CreateAPI("reel")
			api.SetupRoutes(gapi)
			gms.Daemonize(gapi)
		},
	}

	listsourcesCmd = &cobra.Command{
		Use:   "listsources",
		Short: "run reel with a config against an app to list database sources.",
		Long:  `run reel with a config against an app to list database sources.`,
		Run: func(cmd *cobra.Command, args []string) {
			gms := &governor.ManagerService{}
			gms.InitManager(config)
			gapi := gms.CreateAPI(application)	
			reel.InitApp(application, gapi)
			reel.PrintSources(application, gapi)
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

	listsourcesCmd.Flags().StringVarP(&config, "config", "c", "", "config file")
	listsourcesCmd.Flags().StringVarP(&application, "application", "a", "", "application name")
	listsourcesCmd.MarkFlagRequired("config")
	listsourcesCmd.MarkFlagRequired("application")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(managerCmd)
	rootCmd.AddCommand(listsourcesCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
