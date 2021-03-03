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

var driverBtrfsCmd = &cobra.Command{
	Use:   "btrfs",
	Short: "Describes btrfs driver",
	Long: `btrfs driver is a wrapper to "btrfs" program, the toolbox to manage btrfs filesystem.

btrfs driver manages snapshots on a btrfs filesystem. It does not operates on the
filesystem directly. Instead, it calls "btrfs" program to do the job.

Driver actions are mapped to following commands:

- test:   btrfs sub show
- create: btrfs snap
- list:   filepath.Glob() in Golang
- delete: btrfs sub del
- send:   btrfs send
- recv:   btrfs receive

btrfs is really flexible. That makes it problematic sometimes. The fs structure can
be in any kind, so it's not realistic to auto-detect it.

Since btrfs lets users to decide where their backup should place, moobk follows it.
moobk does not force users to "put something at somewhere", user takes the
responsibility to tell moobk where the repo is. As the result, btrfs driver *DOES
NOT* support recursive snapshotting nor sending.

To specify repo location, use filesystem path like when you're using "btrfs" program.
moobk *DOES NOT* support relative path.

This driver supports only one option: bin, which specifies path to "btrfs" program. For example, add following url query string to your repo URL

  drv_bin=/opt/btrfs/bin/btrfs

will make this driver executes "/opt/btrfs/bin/btrfs" instead of default "btrfs".
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(cmd.Long)
	},
}

func init() {
	driverCmd.AddCommand(driverBtrfsCmd)
}
