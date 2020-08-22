package cmd

import (
	"errors"
	"fmt"

	"github.com/risset/gobooru/backend"
	"github.com/spf13/cobra"
)

type tag struct {
	limit int
	order int
}

func newTagCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tag",
		Short: "Search for tag",
		Long:  `Search for tag. Accepts patterns in query.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("requires argument(s).")
			}

			return nil
		},
	}

	c := &tag{}

	cmd.PersistentFlags().IntVarP(&c.limit, "limit", "n", 20, "number of tags to retrieve (max 100)")
	cmd.PersistentFlags().IntVarP(&c.order, "order", "o", 0, "sort order: 0 = date, 1 = name, 2 = count")

	cmd.RunE = c.tagCmd

	return cmd
}

func (c *tag) tagCmd(cmd *cobra.Command, args []string) error {
	params := backend.BuildTagParams(api, args[0], c.order)
	if len(params) == 0 {
		return fmt.Errorf("no search parameters found.")
	}

	data, err := backend.Search("tag", params, api)
	if err != nil {
		return err
	}

	backend.ShowJSON(data)

	return nil
}
