package set

import (
	"fmt"
	"io"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/khulnasoft-lab/goctl/v2/internal/config"
	"github.com/khulnasoft-lab/goctl/v2/pkg/cmd/alias/shared"
	"github.com/khulnasoft-lab/goctl/v2/pkg/cmdutil"
	"github.com/khulnasoft-lab/goctl/v2/pkg/iostreams"
	"github.com/spf13/cobra"
)

type SetOptions struct {
	Config func() (config.Config, error)
	IO     *iostreams.IOStreams

	Name              string
	Expansion         string
	IsShell           bool
	OverwriteExisting bool

	validAliasName      func(string) bool
	validAliasExpansion func(string) bool
}

func NewCmdSet(f *cmdutil.Factory, runF func(*SetOptions) error) *cobra.Command {
	opts := &SetOptions{
		IO:     f.IOStreams,
		Config: f.Config,
	}

	cmd := &cobra.Command{
		Use:   "set <alias> <expansion>",
		Short: "Create a shortcut for a goctl command",
		Long: heredoc.Doc(`
			Define a word that will expand to a full goctl command when invoked.

			The expansion may specify additional arguments and flags. If the expansion includes
			positional placeholders such as "$1", extra arguments that follow the alias will be
			inserted appropriately. Otherwise, extra arguments will be appended to the expanded
			command.

			Use "-" as expansion argument to read the expansion string from standard input. This
			is useful to avoid quoting issues when defining expansions.

			If the expansion starts with "!" or if "--shell" was given, the expansion is a shell
			expression that will be evaluated througoctl the "sh" interpreter when the alias is
			invoked. This allows for chaining multiple commands via piping and redirection.
		`),
		Example: heredoc.Doc(`
			# note: Command Prompt on Windows requires using double quotes for arguments
			$ goctl alias set pv 'pr view'
			$ goctl pv -w 123  #=> goctl pr view -w 123

			$ goctl alias set bugs 'issue list --label=bugs'
			$ goctl bugs

			$ goctl alias set homework 'issue list --assignee @me'
			$ goctl homework

			$ goctl alias set 'issue mine' 'issue list --mention @me'
			$ goctl issue mine

			$ goctl alias set epicsBy 'issue list --author="$1" --label="epic"'
			$ goctl epicsBy vilmibm  #=> goctl issue list --author="vilmibm" --label="epic"

			$ goctl alias set --shell igrep 'goctl issue list --label="$1" | grep "$2"'
			$ goctl igrep epic foo  #=> goctl issue list --label="epic" | grep "foo"
		`),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]
			opts.Expansion = args[1]

			opts.validAliasName = shared.ValidAliasNameFunc(cmd)
			opts.validAliasExpansion = shared.ValidAliasExpansionFunc(cmd)

			if runF != nil {
				return runF(opts)
			}

			return setRun(opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.IsShell, "shell", "s", false, "Declare an alias to be passed througoctl a shell interpreter")
	cmd.Flags().BoolVar(&opts.OverwriteExisting, "clobber", false, "Overwrite existing aliases of the same name")

	return cmd
}

func setRun(opts *SetOptions) error {
	cs := opts.IO.ColorScheme()
	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	aliasCfg := cfg.Aliases()

	expansion, err := getExpansion(opts)
	if err != nil {
		return fmt.Errorf("did not understand expansion: %w", err)
	}

	if opts.IsShell && !strings.HasPrefix(expansion, "!") {
		expansion = "!" + expansion
	}

	isTerminal := opts.IO.IsStdoutTTY()
	if isTerminal {
		fmt.Fprintf(opts.IO.ErrOut, "- Creating alias for %s: %s\n", cs.Bold(opts.Name), cs.Bold(expansion))
	}

	var existingAlias bool
	if _, err := aliasCfg.Get(opts.Name); err == nil {
		existingAlias = true
	}

	if !opts.validAliasName(opts.Name) {
		if !existingAlias {
			return fmt.Errorf("%s Could not create alias %s: already a goctl command or extension",
				cs.FailureIcon(),
				cs.Bold(opts.Name))
		}

		if existingAlias && !opts.OverwriteExisting {
			return fmt.Errorf("%s Could not create alias %s: name already taken, use the --clobber flag to overwrite it",
				cs.FailureIcon(),
				cs.Bold(opts.Name),
			)
		}
	}

	if !opts.validAliasExpansion(expansion) {
		return fmt.Errorf("%s Could not create alias %s: expansion does not correspond to a goctl command, extension, or alias",
			cs.FailureIcon(),
			cs.Bold(opts.Name))
	}

	aliasCfg.Add(opts.Name, expansion)

	err = cfg.Write()
	if err != nil {
		return err
	}

	successMsg := fmt.Sprintf("%s Added alias %s", cs.SuccessIcon(), cs.Bold(opts.Name))
	if existingAlias && opts.OverwriteExisting {
		successMsg = fmt.Sprintf("%s Changed alias %s",
			cs.WarningIcon(),
			cs.Bold(opts.Name))
	}

	if isTerminal {
		fmt.Fprintln(opts.IO.ErrOut, successMsg)
	}

	return nil
}

func getExpansion(opts *SetOptions) (string, error) {
	if opts.Expansion == "-" {
		stdin, err := io.ReadAll(opts.IO.In)
		if err != nil {
			return "", fmt.Errorf("failed to read from STDIN: %w", err)
		}

		return string(stdin), nil
	}

	return opts.Expansion, nil
}
