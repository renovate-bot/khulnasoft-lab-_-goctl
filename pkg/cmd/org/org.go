package org

import (
	"github.com/MakeNowJust/heredoc"
	orgListCmd "github.com/khulnasoft-lab/goctl/v2/pkg/cmd/org/list"
	"github.com/khulnasoft-lab/goctl/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdOrg(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "org <command>",
		Short: "Manage organizations",
		Long:  "Work with Github organizations.",
		Example: heredoc.Doc(`
			$ goctl org list
		`),
		GroupID: "core",
	}

	cmdutil.AddGroup(cmd, "General commands", orgListCmd.NewCmdList(f, nil))

	return cmd
}
