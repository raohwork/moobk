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
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// internalCmd represents the internal command
var internalCmd = &cobra.Command{
	Use:   "internal",
	Short: "internal wrapper to use with remote control",
	Long: `internal, as the name suggests, is for internal use only.

It wraps fs driver to ensure only hard-coded actions are performed.

End-users *SHOULD NEVER* use this directly.
`,
}

var internalFlags = repoFlags{}

func init() {
	rootCmd.AddCommand(internalCmd)

	internalFlags.BindPersist(internalCmd)
}

func intWriteJson(data interface{}, err error) (code int) {
	if err != nil {
		return intWriteErr(err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.Encode(data)
	return
}

func intWriteErr(err error) (code int) {
	fmt.Fprint(os.Stderr, err)
	return 1
}
