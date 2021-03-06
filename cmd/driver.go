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
	Long:  `Show documentation about what "driver" is.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(`moobk abstracts whole backup procedure into few actions:

- test:   Tests if specified target can be managed by this driver.
- create: Creates a new snapshot.
- list:   Lists snapshots.
- delete: Deletes a snapshot.
- send:   Sends a snapshot to somewhere.
- recv:   Receives a snapshot from somewhere.

A driver is piece of code to execute these actions on something. Like, btrfs driver
can take snapshot on a btrfs filesystem.

Run "moobk help driver" to see supported drivers.
`)
	},
}

func init() {
	rootCmd.AddCommand(driverCmd)
}
