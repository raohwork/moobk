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
	"sort"

	"github.com/raohwork/moobk/moodrvs"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List recognized snapshots in the repo",
	Long: `List recognized snapshots in the repo.

moobk recognizes snapshot by just naming convention: "name"-"timestamp". If it does
not recognize other snapshot you made manually, that't totally normal.

Currently the order is hard-coded. Custom sorting might be added in later version.
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		repo, err := listFlags.Repo()
		if err != nil {
			return
		}

		l, err := repo.Snapshots()
		if err != nil {
			return
		}

		sort.Sort(moodrvs.SnapshotSlice(l))
		l = moodrvs.Filter(l, listFlags.Name, listFlags.Count, listFlags.T)

		for _, n := range l {
			fmt.Printf(
				"%s     %s\n",
				n.CreatedAt.Format("2006-01-02 15:04:05"),
				n.RealName(),
			)
		}
		return
	},
}

var listFlags = struct {
	repoFlags
	Name  string
	Count uint
	T     string
}{}

func init() {
	rootCmd.AddCommand(listCmd)

	listFlags.Bind(listCmd)
	f := listCmd.Flags()
	f.StringVarP(&listFlags.Name, "name", "n", "", "lists only snapshots with exactly same name")
	f.UintVarP(&listFlags.Count, "count", "c", 0, "latest n snapshots (0 for all, which is default)")
	f.StringVarP(&listFlags.T, "duration", "d", "", "within duration from now. /[0-9]+[hdwm]/ for n hours, days, weeks or months")
}
