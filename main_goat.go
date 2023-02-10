package main

import (
	"github.com/tmr232/goat"
	"github.com/tmr232/goat/flags"
	"github.com/urfave/cli/v2"
)

func init() {
	goat.Register(generateBlog, goat.RunConfig{
		Flags: []cli.Flag{
			flags.MakeFlag[string]("repo", "", nil),
			flags.MakeFlag[string]("contentDir", "", nil),
			flags.MakeFlag[string]("token", "", nil),
		},
		Name:  "generateBlog",
		Usage: "",
		Action: func(c *cli.Context) error {
			return generateBlog(
				flags.GetFlag[string](c, "repo"),
				flags.GetFlag[string](c, "contentDir"),
				flags.GetFlag[string](c, "token"),
			)
		},
		CtxFlagBuilder: func(c *cli.Context) map[string]any {
			cflags := make(map[string]any)
			cflags["repo"] = flags.GetFlag[string](c, "repo")
			cflags["contentDir"] = flags.GetFlag[string](c, "contentDir")
			cflags["token"] = flags.GetFlag[string](c, "token")
			return cflags
		},
	})
}
