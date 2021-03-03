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

	"github.com/spf13/cobra"
)

// repoCmd represents the repo command
var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Describes what 'repository' is",
	Long:  `Show documentation about what "repository" is, what kind of repositories is supported, and its key behavier.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(`moobk support snapshot management both locally and remotely. Snapshots are stored in "repository" which can vary in different fs driver. 

All repos are specified in URL format: scheme://repo/location?option1=value1&option2=value2. Refer to specific scheme for more info about supported options.

Driver options are pass through query string, too. Run "moobk driver" for more info.

Run "moobk help repo" to see available drivers.
`)
	},
}

func init() {
	rootCmd.AddCommand(repoCmd)
}
