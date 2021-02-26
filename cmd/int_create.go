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
	"github.com/spf13/cobra"
)

// internalCreateCmd represents the internalCreate command
var internalCreateCmd = &cobra.Command{
	Use:   "create target [name]",
	Short: "Create wraps COW.Create() method",
	Long: `as the name suggests, it is for internal use only.

End-users *SHOULD NEVER* use this directly.
`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := internalFlags.Repo()
		if err != nil {
			intWriteErr(err)
			return
		}

		name := ""
		if len(args) > 1 {
			name = args[1]
		}
		intWriteJson(repo.Create(args[0], name))
	},
}

func init() {
	internalCmd.AddCommand(internalCreateCmd)
}
