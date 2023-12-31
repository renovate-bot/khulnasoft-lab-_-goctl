package ruleset

import (
	"github.com/MakeNowJust/heredoc"
	cmdCheck "github.com/khulnasoft-lab/goctl/v2/pkg/cmd/ruleset/check"
	cmdList "github.com/khulnasoft-lab/goctl/v2/pkg/cmd/ruleset/list"
	cmdView "github.com/khulnasoft-lab/goctl/v2/pkg/cmd/ruleset/view"
	"github.com/khulnasoft-lab/goctl/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdRuleset(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ruleset <command>",
		Short: "View info about repo rulesets",
		Long: heredoc.Doc(`
			Repository rulesets are a way to define a set of rules that apply to a repository.
			These commands allow you to view information about them.
		`),
		Aliases: []string{"rs"},
		Example: heredoc.Doc(`
			$ goctl ruleset list
			$ goctl ruleset view --repo OWNER/REPO --web
			$ goctl ruleset check branch-name
		`),
	}

	cmdutil.EnableRepoOverride(cmd, f)
	cmd.AddCommand(cmdList.NewCmdList(f, nil))
	cmd.AddCommand(cmdView.NewCmdView(f, nil))
	cmd.AddCommand(cmdCheck.NewCmdCheck(f, nil))

	return cmd
}
