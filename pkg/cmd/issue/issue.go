package issue

import (
	"github.com/MakeNowJust/heredoc"
	cmdClose "github.com/khulnasoft-lab/goctl/v2/pkg/cmd/issue/close"
	cmdComment "github.com/khulnasoft-lab/goctl/v2/pkg/cmd/issue/comment"
	cmdCreate "github.com/khulnasoft-lab/goctl/v2/pkg/cmd/issue/create"
	cmdDelete "github.com/khulnasoft-lab/goctl/v2/pkg/cmd/issue/delete"
	cmdDevelop "github.com/khulnasoft-lab/goctl/v2/pkg/cmd/issue/develop"
	cmdEdit "github.com/khulnasoft-lab/goctl/v2/pkg/cmd/issue/edit"
	cmdList "github.com/khulnasoft-lab/goctl/v2/pkg/cmd/issue/list"
	cmdLock "github.com/khulnasoft-lab/goctl/v2/pkg/cmd/issue/lock"
	cmdPin "github.com/khulnasoft-lab/goctl/v2/pkg/cmd/issue/pin"
	cmdReopen "github.com/khulnasoft-lab/goctl/v2/pkg/cmd/issue/reopen"
	cmdStatus "github.com/khulnasoft-lab/goctl/v2/pkg/cmd/issue/status"
	cmdTransfer "github.com/khulnasoft-lab/goctl/v2/pkg/cmd/issue/transfer"
	cmdUnpin "github.com/khulnasoft-lab/goctl/v2/pkg/cmd/issue/unpin"
	cmdView "github.com/khulnasoft-lab/goctl/v2/pkg/cmd/issue/view"
	"github.com/khulnasoft-lab/goctl/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdIssue(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue <command>",
		Short: "Manage issues",
		Long:  `Work with GitHub issues.`,
		Example: heredoc.Doc(`
			$ goctl issue list
			$ goctl issue create --label bug
			$ goctl issue view 123 --web
		`),
		Annotations: map[string]string{
			"help:arguments": heredoc.Doc(`
				An issue can be supplied as argument in any of the following formats:
				- by number, e.g. "123"; or
				- by URL, e.g. "https://github.com/OWNER/REPO/issues/123".
			`),
		},
		GroupID: "core",
	}

	cmdutil.EnableRepoOverride(cmd, f)

	cmdutil.AddGroup(cmd, "General commands",
		cmdList.NewCmdList(f, nil),
		cmdCreate.NewCmdCreate(f, nil),
		cmdStatus.NewCmdStatus(f, nil),
	)

	cmdutil.AddGroup(cmd, "Targeted commands",
		cmdView.NewCmdView(f, nil),
		cmdComment.NewCmdComment(f, nil),
		cmdClose.NewCmdClose(f, nil),
		cmdReopen.NewCmdReopen(f, nil),
		cmdEdit.NewCmdEdit(f, nil),
		cmdDevelop.NewCmdDevelop(f, nil),
		cmdLock.NewCmdLock(f, cmd.Name(), nil),
		cmdLock.NewCmdUnlock(f, cmd.Name(), nil),
		cmdPin.NewCmdPin(f, nil),
		cmdUnpin.NewCmdUnpin(f, nil),
		cmdTransfer.NewCmdTransfer(f, nil),
		cmdDelete.NewCmdDelete(f, nil),
	)

	return cmd
}
