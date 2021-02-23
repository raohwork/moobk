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
)

type snap struct{}

func init() { register(snap{}) }

func (_ snap) Name() string { return "snap" }
func (_ snap) Desc() string { return "Take a snapshot." }
func (_ snap) Help() {
	fmt.Print(`snap path [name]

Takes a snapshot at specified path

path   path to take snapshot
name   use name instead of basename(path)
`)
}

func (s snap) Exec(args []string) (ret int) {
	opt := mustG(args)
	if !want(1, opt.args) {
		s.Help()
		return 1
	}

	name := ""
	if len(opt.args) > 1 {
		name = opt.args[1]
	}

	snap, err := opt.repo.Create(opt.args[0], name)
	if err != nil {
		fmt.Println("cannot take snapshot: ", err)
		return 1
	}

	fmt.Println("snapshot created: ", snap.RealName())
	return
}
