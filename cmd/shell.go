// This file is part of moobk.
//
// moobk is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 2 of the License, or
// (at your option) any later version.
//
// moobk is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with moobk.  If not, see <https://www.gnu.org/licenses/>.

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var shellCmd = &cobra.Command{
	Use:   "shell [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:

  $ source <(moobk shell bash)

  # To load completions for each session, execute once:
  # Linux:
  $ moobk shell bash > /etc/bash_completion.d/moobk
  # macOS:
  $ moobk shell bash > /usr/local/etc/bash_completion.d/moobk

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ moobk shell zsh > "${fpath[1]}/_moobk"

  # You will need to start a new shell for this setup to take effect.

fish:

  $ moobk shell fish | source

  # To load completions for each session, execute once:
  $ moobk shell fish > ~/.config/fish/completions/moobk.fish

PowerShell:

  PS> moobk shell powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> moobk shell powershell > moobk.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
}
