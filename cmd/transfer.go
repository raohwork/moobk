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
	"io"
	"sync"

	"github.com/raohwork/moobk/moodrvs"
	"github.com/spf13/cobra"
)

// transferCmd represents the list command
var transferCmd = &cobra.Command{
	Aliases:   []string{"t", "trans", "sync", "send", "s"},
	Use:       "transfer local remote",
	Short:     "Transfers local snapshots to remote",
	Long:      `transfer snapshots that exist in local but missing in remote.`,
	ValidArgs: []string{"local", "remote"},
	Args:      cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		local, err := moodrvs.GetRunner(args[0], fs)
		if err != nil {
			return fmt.Errorf("cannot init local repo: %w", err)
		}

		remote, err := moodrvs.GetRunner(args[1], fs)
		if err != nil {
			return fmt.Errorf("cannot init remote repo: %w", err)
		}

		lSnaps, err := local.Snapshots()
		if err != nil {
			return fmt.Errorf("cannot gather snapshot info in local: %w", err)
		}
		lSnaps = moodrvs.Filter(lSnaps, transferFlags.Name, 0, "")
		if len(lSnaps) == 0 {
			// nothing to do
			return
		}

		rSnaps, err := remote.Snapshots()
		if err != nil {
			return fmt.Errorf("cannot gather snapshot info in remote: %w", err)
		}
		rSnaps = moodrvs.Filter(rSnaps, transferFlags.Name, 0, "")

		diffs := moodrvs.Diff(lSnaps, rSnaps)
		for _, diff := range diffs {
			err = transfer(local, remote, diff)
			if err != nil {
				fmt.Println()
				fmt.Println(err)
				return nil
			}
		}

		return
	},
}

func transfer(lRepo, rRepo moodrvs.Runner, diff moodrvs.SnapshotDiff) (err error) {
	if len(diff.Missing) < 1 {
		// nothing to do
		return
	}

	base := diff.Base
	if base == nil {
		base = &diff.Missing[0]
	}

	for idx, s := range diff.Missing {
		fmt.Printf(
			"Sending %s ~ %s from %s://%s to %s://%s ... ",
			base.RealName(),
			s.RealName(),
			lRepo.RunnerName(),
			lRepo.BackupPath(),
			rRepo.RunnerName(),
			rRepo.BackupPath(),
		)

		if transferFlags.DryRun {
			fmt.Println("skipped.")
			continue
		}

		r, w := io.Pipe()
		var rerr, werr error
		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			werr = lRepo.Send(*base, s, w)
			wg.Done()
			w.Close()
		}()
		go func() {
			rerr = rRepo.Recv(s, r)
			wg.Done()
			w.Close()
		}()
		wg.Wait()

		if werr != nil {
			err = fmt.Errorf("cannot send snapshot: %w", werr)
			return
		}
		if rerr != nil {
			err = fmt.Errorf("cannot receive snapshot: %w", rerr)
			return
		}

		base = &diff.Missing[idx]
		fmt.Println("done.")
	}

	return
}

var transferFlags = struct {
	Name   string
	DryRun bool
}{}

func init() {
	rootCmd.AddCommand(transferCmd)

	transferCmd.Flags().StringVarP(&transferFlags.Name, "name", "n", "", "optional filter. Transfers only matching snapshot")
	transferCmd.Flags().BoolVarP(&transferFlags.DryRun, "dry-run", "d", false, "do not send/receive snapshot")
}
