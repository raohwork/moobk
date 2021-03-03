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

var repoLocalCmd = &cobra.Command{
	Use:   "local",
	Short: "Describes local scheme, which accesses local repo",
	Long: `local scheme is for repo hosted on the same machine you ran moobk.

The path part depends on which fs driver you are using. For btrfs, it should be a
folder managed by btrfs, like /btrfs_root/some_folder (or even /btrfs_root). For
zfs, it should be a zfs volume, like rpool/some_volume.

Here are few examples:
  local:///btrfs/backup
  local://rpool/backup
  local:///rpool/backup   (leading spash is stripped when using zfs driver)

`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(cmd.Long)
	},
}

func init() {
	repoCmd.AddCommand(repoLocalCmd)
}
