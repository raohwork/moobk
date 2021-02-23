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
	"sort"

	"github.com/raohwork/moobk/moodrvs"
)

type list struct{}

func init() { register(list{}) }

func (_ list) Name() string { return "list" }
func (_ list) Desc() string { return "list snapshots in the repo." }
func (_ list) Help() {
	fmt.Print(`list [-name name] [-count count] [-duration duration_spec]

list snapshots in the repo

Optional filter flags are:

-name string         snapshots with exactly same name
-count uint          latest n snapshots (0 for all, which is default)
-duration string     within duration from now. /[0-9]+[hdwm]/ for n hours, days,
                     weeks or months.
`)
}

func (s list) Exec(args []string) (ret int) {
	// filtering
	var (
		name string
		cnt  uint
		t    string
	)
	flagSet.StringVar(&name, "name", "", "")
	flagSet.UintVar(&cnt, "count", 0, "")
	flagSet.StringVar(&t, "duration", "", "")

	opt := mustG(args)
	l, err := opt.repo.Snapshots()
	if err != nil {
		fmt.Println("cannot list snapshots: ", err)
		return 1
	}

	sort.Sort(moodrvs.SnapshotSlice(l))

	l = moodrvs.Filter(l, name, cnt, t)

	for _, n := range l {
		fmt.Printf(
			"%s     %s\n",
			n.CreatedAt.Format("2006-01-02 15:04:05"),
			n.RealName(),
		)
	}

	return
}
