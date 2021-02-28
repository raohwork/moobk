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

var driverZfsCmd = &cobra.Command{
	Use:   "zfs",
	Short: "Describes zfs driver",
	Long: `zfs driver is a wrapper to "zfs" program, the main binary to manage zfs datasets.

zfs driver manages snapshots on a zfs dataset. It does not operates on the dataset
directly. Instead, it calls "zfs" program to do the job.

Driver actions are mapped to following commands:

- test:   zfs get -H type (grabs only type == "filesystem")
- create: zfs snap
- list:   zfs list -H -d 1 -t snapshot
- delete: zfs destroy
- send:   zfs send -i
- recv:   zfs recv -duF

Snapshots in zfs are forced to store under the same dataset. Through it has builtin
support to recursive snapshotting in "zfs" program, this driver *DOES NOT* support
it to keep same behavier with btrfs driver. Use zfsr driver if you need recursive
snapshotting support.

To specify repo location, use the same syntax when you run "zfs" program
(pool_name/dataset/path).

SPECIAL NOTICE:
This driver receives snapshots with "-F" flag enabled, YOU HAVE TO TAKE CARE OF IT!
(See zfs manual for detailed description about "zfs receive -F")

This driver supports only one option: bin, which specifies path to "zfs" program. For example, add following url query string to your repo URL

  drv_bin=/opt/zfs/bin/zfs

will make this driver executes "/opt/zfs/bin/zfs" instead of default "zfs".
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(cmd.Long)
	},
}

func init() {
	driverCmd.AddCommand(driverZfsCmd)
}
