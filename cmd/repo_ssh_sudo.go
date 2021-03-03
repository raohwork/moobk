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
	"fmt"

	"github.com/spf13/cobra"
)

var repoSSHSUDOCmd = &cobra.Command{
	Use:   "ssh+sudo",
	Short: "Describes ssh+sudo scheme, which uses ssh and sudo to access remote repo",
	Long: `ssh+sudo scheme uses ssh and sudo to access repo hosted on another machine.

This is roughly same as ssh, with enhanced security. The user on remote gains
priveledge using sudo, so you may prevent it from access your device directly.

You *MUST* add "NOPASSWD" tag on the user. It's also recommended to "tighten"
sudo settings such as verify hash digest of moobk.

Here's a simple example in sudoers:
  user ALL= PASSWD:ALL, NOPASSWD:/path/moobk

Here's a more complex example:
  Cmnd_Alias MOO= ssh224:xxxxxxx /path/moobk internal *
  user ALL= PASSWD:ALL, NOPASSWD:MOO

Here are few url examples:
  ssh+sudo://user@example.com/btrfs/backup
  ssh+sudo://user@example.com:1234/btrfs/backup (using different port)
  ssh+sudo://example.com/rpool/backup           (using default/current user)

Query string options for ssh scheme are supported.

`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(cmd.Long)
	},
}

func init() {
	repoCmd.AddCommand(repoSSHSUDOCmd)
}
