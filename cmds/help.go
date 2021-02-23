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
	"os"
	"sort"
	"strings"
)

type help struct{}

func init() { register(help{}) }

func (_ help) Name() string { return "help" }
func (_ help) Desc() string { return "Display usage." }
func (_ help) Help() {
	fmt.Printf(`Usage: %s command [global options] [command arguments...]`+"\n", os.Args[0])
	fmt.Println()
	fmt.Print(`Supported global options:

-t fs         Specify filesystem type. btrfs and zfs are supported. You may use envvar
              MOOBK_FS instead. See moobk help driver for detail.
-r repo       URL format to snapshot repo. There're 3 schema supported: local, ssh,
              ssh+sudo. See moobk help repo for detail. You may use envvar MOOBK_REPO
              instead.

Supported commands:

`)
	keys := make([]string, 0, len(available))
	for k := range available {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		c := available[k]
		fmt.Printf("%-13s %s\n", c.Name(), c.Desc())
	}
}

func (h help) Exec(args []string) (ret int) {
	if len(args) > 0 {
		c := strings.ToLower(args[0])
		x, ok := available[c]
		if ok {
			x.Help()
			return
		}
	}

	h.Help()
	return
}
