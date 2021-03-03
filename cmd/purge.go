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
	Aliases: []string{"p", "d", "del", "delete", "rm", "remove"},
	Use:     "purge local remote",
	Short:   "Purge snapshots from local according to remote",
	Long: `purge deletes snapshots that exist in both local and remote.

The format of "reserve" flag can be:
    integer:       do not delete latest n snapshots.
    [0-9]+[hdwm]:  do not delete snapshots newer than n hour/day/week/month ago from
                   now.

Depends on existent in local and remote, a snapshot can be:

    - Synced:  the snapshot exists in both repo
    - Missing: exists in either local or remote (but not both), and "newer" than
               latest synced snapshot
    - Orphan:  exists in either local or remote (but not both), and "older" than
               latest synced snapshot

By default, purge deletes synced snapshots "on both local and remote", excepts
latest one. You may use "reserve" flag to leave more synced snapshots untouched.
Since latest synced snapshot is never purged, if you specify 1 for reserve flag,
there will be at most 2 synced snapshots after.

The "orphan" flag purges all orphan snapshots.
`,
	ValidArgs: []string{"local", "remote"},
	Args:      cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		local, err := moodrvs.GetRunner(args[0], fs)
		if err != nil {
			return fmt.Errorf("cannot connect to local: %w", err)
		}

		remote, err := moodrvs.GetRunner(args[1], fs)
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

		if purgeFlags.Reserve != "" {
			i, err := strconv.ParseUint(purgeFlags.Reserve, 10, 32)
			if err == nil {
				cnt = uint(i)
			} else {
				dur = purgeFlags.Reserve
			}
		}
		dupes, orphanL, orphanR := moodrvs.FindDupe(lSnaps, rSnaps)
		for _, arr := range dupes {
			arr = moodrvs.Reserve(arr, cnt, dur)
			for _, d := range arr {
				fmt.Printf(
					"Deleting synced %s ... ",
					d.RealName(),
				)

				if purgeFlags.DryRun {
					fmt.Println("skipped.")
					continue
				}

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

		if !purgeFlags.Orphan {
			return
		}

		for _, arr := range orphanL {
			arr = moodrvs.Reserve(arr, cnt, dur)
			for _, d := range arr {
				fmt.Printf(
					"Deleting orphan %s ... ",
					d.RealName(),
				)

				if purgeFlags.DryRun {
					fmt.Println("skipped.")
					continue
				}

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

		for _, arr := range orphanR {
			arr = moodrvs.Reserve(arr, cnt, dur)
			for _, d := range arr {
				fmt.Printf(
					"Deleting orphan %s ... ",
					d.RealName(),
				)

				if purgeFlags.DryRun {
					fmt.Println("skipped.")
					continue
				}

				if err := remote.Delete(d); err != nil {
					fmt.Println("error.")
					return fmt.Errorf(
						"cannot delete remote snapshot: %w",
						err,
					)
				}
				fmt.Println("remote. ")
			}
		}

		return
	},
}

var purgeFlags = struct {
	Name    string
	Reserve string
	Orphan  bool
	DryRun  bool
}{}

func init() {
	rootCmd.AddCommand(purgeCmd)

	f := purgeCmd.Flags()
	f.StringVarP(&purgeFlags.Name, "name", "n", "", "purges only snapshots with exactly same name")
	f.StringVarP(&purgeFlags.Reserve, "reserve", "r", "", "reserves more snapshots")
	f.BoolVarP(&purgeFlags.Orphan, "orphan", "o", false, "also deletes orphan snapshots")
	f.BoolVarP(&purgeFlags.DryRun, "dry-run", "d", false, "do not really deletes snapshots")
}
