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
	"strconv"

	"github.com/raohwork/moobk/moodrvs"
)

type purge struct{}

func init() { register(purge{}) }

func (_ purge) Name() string { return "purge" }
func (_ purge) Desc() string { return "purge snapshots from local according to remote." }
func (_ purge) Help() {
	fmt.Print(`purge [-name name] remote_repo [reserve]

purge deletes snapshots that exist in both local and remote.

-name string  Optional filter. Deletes only matching snapshot.
remote_repo   Repo to check for snapshots.
reserve       Reserves some snapshots, can be:
              integer:       do not delete latest n snapshots.
              [0-9]+[hdwm]:  do not delete snapshots newer than n hour/day/week/month
                             ago from now.

purge will never delete orphan snapshot. Latest synced snapshot is also preserved.
If you specify 1 for reserve, there will be at most 2 synced snapshots.
`)
}

func (s purge) Exec(args []string) (ret int) {
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

	// prepare reserve filters
	cnt := uint(0)
	dur := ""

	if len(opt.args) > 1 {
		i, err := strconv.ParseUint(opt.args[1], 10, 32)
		if err == nil {
			cnt = uint(i)
		} else {
			dur = opt.args[1]
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
				fmt.Println(err)
				return 1
			}
			fmt.Print("remote. ")

			if err := local.Delete(d); err != nil {
				fmt.Println(err)
				return 1
			}
			fmt.Println("local.")
		}
	}

	return
}
