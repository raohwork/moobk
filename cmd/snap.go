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

	"github.com/raohwork/moobk/moodrvs"
	"github.com/spf13/cobra"
)

// snapCmd represents the list command
var snapCmd = &cobra.Command{
	Aliases: []string{"c", "create", "a", "add", "snapshot"},
	Use:     "snap repo target [name]",
	Short:   "Take snapshot",
	Long: `Take a snapshot of specified target, named it as name, store it to repo

Some drivers might have restrictions about where to store. Run "moobk driver" for more info.
`,
	ValidArgs: []string{
		"repo", "target", "name",
	},
	Args: cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		name := ""
		if len(args) > 2 {
			name = args[2]
		}

		repo, err := moodrvs.GetRunner(args[0], fs)
		if err != nil {
			return
		}

		snap, err := repo.Create(args[1], name)
		if err != nil {
			return
		}

		fmt.Println("snapshot created: ", snap.RealName())
		return
	},
}

func init() {
	rootCmd.AddCommand(snapCmd)
}
