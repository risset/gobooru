package cmd

import (
	"strings"

	"github.com/risset/gobooru/backend"
	"github.com/spf13/cobra"
)

type post struct {
	dl     bool
	limit  int
	random bool
}

func newPostCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "post",
		Short: "Get post data",
		Long:  `Get post data for specified tags. Optionally, download images.`,
	}

	c := &post{}

	cmd.PersistentFlags().BoolVarP(&c.dl, "dl", "d", false, "toggle downloading of images")
	cmd.PersistentFlags().IntVarP(&c.limit, "limit", "n", 20, "number of posts to retrieve (max 200)")
	cmd.PersistentFlags().BoolVarP(&c.random, "random", "r", false, "ignore tags argument and retrieve random post(s)")

	cmd.RunE = c.postCmd

	return cmd
}

func (c *post) postCmd(cmd *cobra.Command, args []string) error {
	tags := strings.Join(args[:], " ")
	s := backend.BuildPostSearch(backend.API(api), tags, c.limit, c.random)

	data, err := backend.GetData(s)
	if err != nil {
		return err
	}

	if c.dl {
		dir := backend.GetImgDirName(c.random, tags)
		err := backend.GetAllImages(data, dir)
		if err != nil {
			return err
		}
	} else {
		backend.ShowJSON(data)
	}

	return nil
}
