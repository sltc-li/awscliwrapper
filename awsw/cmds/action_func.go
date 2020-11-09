package cmds

import (
	"github.com/urfave/cli"

	"github.com/li-go/awscliwrapper"
)

func ActionFunc(fn func(w *awscliwrapper.Wrapper) error) cli.ActionFunc {
	return ActionFuncWithContext(func(c *cli.Context, w *awscliwrapper.Wrapper) error {
		return fn(w)
	})
}

func ActionFuncWithContext(fn func(c *cli.Context, w *awscliwrapper.Wrapper) error) cli.ActionFunc {
	return func(c *cli.Context) error {
		w, err := awscliwrapper.New()
		if err != nil {
			return err
		}

		return fn(c, w)
	}
}
