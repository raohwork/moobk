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
	"errors"
	"os"

	"github.com/raohwork/moobk/moodrvs"
	"github.com/spf13/cobra"
)

// internalRecvCmd represents the internalRecv command
var internalRecvCmd = &cobra.Command{
	Use:   "recv snap",
	Short: "Recv wraps COW.Recv() method",
	Long: `as the name suggests, it is for internal use only.

End-users *SHOULD NEVER* use this directly.
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := internalFlags.Repo()
		if err != nil {
			intWriteErr(err)
			return
		}

		s, ok := moodrvs.ParseSnapshot(args[0])
		if !ok {
			intWriteErr(errors.New("invalid s snapshot"))
			return
		}

		if err = repo.Recv(s, os.Stdin); err != nil {
			intWriteErr(err)
			return
		}
	},
}

func init() {
	internalCmd.AddCommand(internalRecvCmd)
}
