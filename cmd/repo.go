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

// repoCmd represents the repo command
var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Describes what 'repository' is",
	Long:  `Show documentation about what "repository" is, what kind of repositories is supported, and its key behavier.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(`moobk support snapshot management both locally and remotely. Snapshots are stored in "repository" which can vary in different fs driver. For btrfs, a repository is actually a folder; For zfs, a repository is a file system or volume.

All repos are specified in URL format. Currently there are 4 schemes supported:

- local

    local means "the repo is hosted at localhost". The path part depends on which fs
    driver you are using.
    For btrfs, it should be a folder managed by btrfs, like /btrfs_root/some_folder
    (or even /btrfs_root).
    For zfs, it should be a zfs volume, like rpool/some_volume.

    Here are few examples:
      local:///btrfs/backup
      local://rpool/backup
      local:///rpool/backup   (leading spash is stripped when using zfs driver)


- ssh

    ssh accesses remote repo using ssh. To use this scheme, you *MUST* setup remote
    environment first:
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

    A common pitfall is, running moobk with sudo on local, which runs ssh as root.
    You may use query string to tell ssh correct key file or, just setup ssh env for
    root.

- ssh+sudo (recommended for remote accessing)

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
`)
	},
}

func init() {
	rootCmd.AddCommand(repoCmd)
}
