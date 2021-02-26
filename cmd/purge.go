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
	"strconv"

	"github.com/raohwork/moobk/moodrvs"
	"github.com/spf13/cobra"
)

// purgeCmd represents the purge command
var purgeCmd = &cobra.Command{
	Use:   "purge remote [reserve]",
	Short: "Purge snapshots from local according to remote",
	Long: `purge deletes snapshots that exist in both local and remote.

The argument "reserve" can be:
    integer:       do not delete latest n snapshots.
    [0-9]+[hdwm]:  do not delete snapshots newer than n hour/day/week/month ago from
                   now.

purge will never delete orphan snapshot. Latest synced snapshot is also preserved.
If you specify 1 for reserve, there will be at most 2 synced snapshots.
`,
	ValidArgs: []string{"reserve\t[0-9]+([hdwm])"},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		local, err := purgeFlags.Repo()
		if err != nil {
			return fmt.Errorf("cannot connect to local: %w", err)
		}

		remote, err := moodrvs.GetRunner(args[0], local.DriverName())
		if err != nil {
			return fmt.Errorf("cannot connect to remote: %w", err)
		}

		lSnaps, err := local.Snapshots()
		if err != nil {
			return fmt.Errorf("cannot gather snapshot info in local: %w", err)
		}
		lSnaps = moodrvs.Filter(lSnaps, purgeFlags.Name, 0, "")
		if len(lSnaps) == 0 {
			// nothing to do
			return
		}

		rSnaps, err := remote.Snapshots()
		if err != nil {
			return fmt.Errorf("cannot gather snapshot info in remote: %w", err)
		}
		rSnaps = moodrvs.Filter(rSnaps, purgeFlags.Name, 0, "")

		// prepare reserve filters
		cnt := uint(0)
		dur := ""

		if len(args) > 1 {
			i, err := strconv.ParseUint(args[1], 10, 32)
			if err == nil {
				cnt = uint(i)
			} else {
				dur = args[1]
			}
		}
		dupes := moodrvs.FindDupe(lSnaps, rSnaps)
		for _, arr := range dupes {
			arr = moodrvs.Reserve(arr, cnt, dur)
			for _, d := range arr {
				fmt.Printf(
					"Deleting %s ... ",
					d.RealName(),
				)
				if err := remote.Delete(d); err != nil {
					fmt.Println("error.")
					return fmt.Errorf(
						"cannot delete remote snapshot: %w",
						err,
					)
				}
				fmt.Print("remote. ")

				if err := local.Delete(d); err != nil {
					fmt.Println("error.")
					return fmt.Errorf(
						"cannot delete local snapshot: %w",
						err,
					)
				}
				fmt.Println("local.")
			}
		}

		return
	},
}

var purgeFlags = struct {
	repoFlags
	Name string
}{}

func init() {
	rootCmd.AddCommand(purgeCmd)

	purgeFlags.Bind(purgeCmd)
	f := purgeCmd.Flags()
	f.StringVarP(&purgeFlags.Name, "name", "n", "", "purges only snapshots with exactly same name")
}
