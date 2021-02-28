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

var repoSSHCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Describes ssh scheme, which uses ssh to access remote repo",
	Long: `ssh scheme uses ssh to access repo hosted on another machine.

To use this scheme, you *MUST* setup remote environment first:
1. Create a user that capable to runs fs tools (like "btrfs" or "zfs"), or use
   root (not recommended).
2. Setup public key autheticating on remote, allowing passwordless login (optional).
3. Install moobk on remote user, add it to PATH envvar.
4. Prepare local ssh environment to connect to remote.

After these steps, you are able to manage remote repos from local.

You can pass some options using query string:

- moobk:  Specify remote moobk binary name
- ssh:    Specify local ssh binary name
- ssh_X:  Passes "-X" to ssh

Here are few url examples:
  ssh://user@example.com/btrfs/backup
  ssh://user@example.com:1234/btrfs/backup (using different port)
  ssh://example.com/rpool/backup           (using default/current user)
  ssh://example.com/rpool/backup?moobk=/bin/moo
                                           (moobk binary is located at /bin/moo)

Here is complicated example of query string:

  ssh=/opt/ssh/bin/ssh&moobk=/usr/local/bin/moobk&ssh_4&ssh_i=/path/to/private.key

It will execute local ssh client at /opt/ssh/bin/ssh, passes
"-4 -i /path/to/private.key" to it, to execute /usr/local/bin/moobk on remote.

A common pitfall is, running moobk with sudo on local, which runs ssh as root. You
may use query string to tell ssh correct key file or, just setup ssh env for root.

`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(cmd.Long)
	},
}

func init() {
	repoCmd.AddCommand(repoSSHCmd)
}
