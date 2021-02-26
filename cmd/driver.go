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

// driverCmd represents the driver command
var driverCmd = &cobra.Command{
	Use:   "driver",
	Short: "Describes what 'driver' is",
	Long:  `Show documentation about what "driver" is, which driver is supported, and its key behavier.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(`moobk supports two CoWfs now: btrfs and zfs

moobk *DOES NOT* interact with the filesystem directly. It executes utilities
provided by the fs to do management job (You have to install it by yourself). There's
a piece of code to "translate" jobs to command executions, which calls "driver", or
"COW" in source code.

- btrfs driver

    btrfs is really flexible. That makes it problematic sometimes. The fs structure
    can be in any kind, so it's not realistic to auto-detect it.

    Since btrfs lets users to decide where their backup should place, moobk follows
    it. moobk does not force users to "put something at somewhere", user takes the
    responsibility to tell moobk where the repo is. As the result, btrfs driver
    *DOES NOT* support recursive snapshotting nor sending.

    To specify repo location, use filesystem path like when your use "btrfs" program.
    moobk *DOES NOT* support relative path.

- zfs driver

    As opposite, zfs is not that flexible. Snapshots are forces to store under same
    zfs dataset. To the bright side, recursive snapshotting and sending are built in
    It is an error to save snapshots into different dataset when using moobk. To keep
    same behavier with btrfs driver, this driver *DOES NOT* support recursive
    snapshotting nor sending.

    To specify repo location, use same syntex "zfs" program accepts
    (pool_name/dataset/path).

    SPECIAL NOTICE:
    This driver receives snapshots with "-F" flag enabled, YOU HAVE TO TAKE CARE OF
    IT! (See zfs manual for detailed description about "zfs receive -F")

- zfsr driver

    This driver is another zfs driver, but follows zfs convntion: recursive
    snapshotting and sending. It enables "-r" flag when creating and sending snashot.

    Affected actions are
      - Create: enables "-r"
      - Delete: enables "-r"
      - Send:   enables "-R"
`)
	},
}

func init() {
	rootCmd.AddCommand(driverCmd)
}
