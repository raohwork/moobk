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

var driverZfsRCmd = &cobra.Command{
	Use:   "zfs",
	Short: "Describes zfsr driver, a zfs driver enables recursive snapshotting",
	Long: `zfsr driver is roughly same with zfs driver, but always takes snapshots recusively.

All pros, cons and options of zfs driver applies to this driver. The only difference
is:

- It invokes "zfs snapshot" with "-r" flag when taking snapshot.
- It invokes "zfs destroy" with  "-r" flag when deleting snapshot.
- It invokes "zfs send" with "-R" flag when sending snapshot.

See zfs driver document for more info.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(cmd.Long)
	},
}

func init() {
	driverCmd.AddCommand(driverZfsRCmd)
}
