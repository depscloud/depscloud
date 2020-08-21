package completion

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	return &cobra.Command{
		Use:   "completion <bash|powershell|zsh>",
		Short: "Generate command completion for different shells",
		RunE: func(cmd *cobra.Command, args []string) error {
			sh := "bash"
			if len(args) > 0 {
				sh = args[0]
			}

			switch sh {
			case "zsh":
				return cmd.GenZshCompletion(os.Stdout)
			case "powershell":
				return cmd.GenPowerShellCompletion(os.Stdout)
			case "bash":
				return cmd.GenBashCompletion(os.Stdout)
			}

			return fmt.Errorf("unrecognized shell: %s", sh)
		},
	}
}
