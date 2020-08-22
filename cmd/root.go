package cmd

import (
	"log"

	"github.com/risset/gobooru/backend"
	"github.com/spf13/cobra"
)

// global flags
var api int

// define a new root CLI command
func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gobooru",
		Short: "Minimal CLI booru client",
		Long:  `A simple command line client for booru sites, with support for batch downloads of images and searching for tags.`,
	}

	cmd.PersistentFlags().IntVarP(&api, "api", "a", int(backend.DANBOORU), "API to use: 0 = danbooru (default), 1 = gelbooru, 2 = konachan")

	return cmd
}

// add subcommands to the root command
func buildRoot() *cobra.Command {
	rootCmd := newRootCmd()

	subCmd := []*cobra.Command{
		newPostCmd(),
		newTagCmd(),
	}

	for _, cmd := range subCmd {
		rootCmd.AddCommand(cmd)
	}

	return rootCmd
}

// build root and subcommands and execute them
func Execute() {
	rootCmd := buildRoot()
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
