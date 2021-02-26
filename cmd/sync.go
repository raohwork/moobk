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

// syncCmd represents the list command
var syncCmd = &cobra.Command{
	Aliases:   []string{"send", "t", "trans", "transfer", "s"},
	Use:       "sync remote",
	Short:     "Transfers local snapshots to remote",
	Long:      `sync send snapshots that exist in local but missing in remote.`,
	ValidArgs: []string{"remote"},
	Args:      cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		local, err := syncFlags.Repo()
		if err != nil {
			return fmt.Errorf("cannot init local repo: %w", err)
		}

		remote, err := moodrvs.GetRunner(args[0], local.DriverName())
		if err != nil {
			return fmt.Errorf("cannot init remote repo: %w", err)
		}

		lSnaps, err := local.Snapshots()
		if err != nil {
			return fmt.Errorf("cannot gather snapshot info in local: %w", err)
		}
		lSnaps = moodrvs.Filter(lSnaps, syncFlags.Name, 0, "")
		if len(lSnaps) == 0 {
			// nothing to do
			return
		}

		rSnaps, err := remote.Snapshots()
		if err != nil {
			return fmt.Errorf("cannot gather snapshot info in remote: %w", err)
		}
		rSnaps = moodrvs.Filter(rSnaps, syncFlags.Name, 0, "")

		diffs := moodrvs.Diff(lSnaps, rSnaps)
		for _, diff := range diffs {
			transfer(local, remote, diff)
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
			"Sending %s ~ %s from %s://%s to %s://%s\n",
			base.RealName(),
			s.RealName(),
			lRepo.RunnerName(),
			lRepo.BackupPath(),
			rRepo.RunnerName(),
			rRepo.BackupPath(),
		)

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

		ok := true
		if werr != nil {
			ok = false
			fmt.Println("Error sending snapshot: ", werr)
			err = werr
		}
		if rerr != nil {
			ok = false
			fmt.Println("Error sending snapshot: ", rerr)
			err = rerr
		}
		if !ok {
			return
		}

		base = &diff.Missing[idx]
	}

	return
}

var syncFlags = struct {
	repoFlags
	Name string
}{}

func init() {
	rootCmd.AddCommand(syncCmd)

	syncFlags.Bind(syncCmd)
	syncCmd.Flags().StringVarP(&syncFlags.Name, "name", "n", "", "optional filter. Sends only matching snapshot")
}
