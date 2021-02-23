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

package cmds

import (
	"fmt"
	"io"
	gsync "sync"

	"github.com/raohwork/moobk/moodrvs"
)

type sync struct{}

func init() { register(sync{}) }

func (_ sync) Name() string { return "sync" }
func (_ sync) Desc() string { return "sync snapshots from local to remote." }
func (_ sync) Help() {
	fmt.Print(`sync [-name name] remote_repo

sync send snapshots that exist in local but missing in remote.

-name string  Optional filter. Sends only matching snapshot.
remote_repo   Repo to receive missing snapshots, see moobk help repo for detail.
`)
}

func (s sync) Exec(args []string) (ret int) {
	var name string
	flagSet.StringVar(&name, "name", "", "")
	opt := mustG(args)
	if !want(1, opt.args) {
		s.Help()
		return 1
	}

	local := opt.repo
	remote, err := moodrvs.GetRunner(opt.args[0], local.DriverName())
	if err != nil {
		fmt.Println("cannot connect to remote: ", err)
		return 1
	}

	lSnaps, err := local.Snapshots()
	if err != nil {
		fmt.Println("cannot gather snapshot info in local: ", err)
		return 1
	}
	lSnaps = moodrvs.Filter(lSnaps, name, 0, "")
	if len(lSnaps) == 0 {
		// nothing to do
		return
	}

	rSnaps, err := remote.Snapshots()
	if err != nil {
		fmt.Println("cannot gather snapshot info in remote: ", err)
		return 1
	}
	rSnaps = moodrvs.Filter(rSnaps, name, 0, "")

	diffs := moodrvs.Diff(lSnaps, rSnaps)
	for _, diff := range diffs {
		s.transfer(local, remote, diff)
	}

	return
}

func (s sync) transfer(lRepo, rRepo moodrvs.Runner, diff moodrvs.SnapshotDiff) (err error) {
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
		wg := gsync.WaitGroup{}
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
