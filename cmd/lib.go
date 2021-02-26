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
	"strings"

	"github.com/raohwork/moobk/moodrvs"
	"github.com/spf13/cobra"
)

type repoFlags struct {
	Type string
	URI  string
}

func (r *repoFlags) Bind(cmd *cobra.Command) {
	f := cmd.Flags()
	f.StringVarP(&r.Type, "type", "t", "btrfs", "filesystem type, only btrfs and zfs are supported.")
	f.StringVarP(&r.URI, "repo", "r", "", "URL format to snapshot repo, see moobk help repo for detail")
	cmd.MarkFlagRequired("repo")
}

func (r *repoFlags) BindPersist(cmd *cobra.Command) {
	f := cmd.PersistentFlags()
	f.StringVarP(&r.Type, "type", "t", "btrfs", "filesystem type, only btrfs and zfs are supported.")
	f.StringVarP(&r.URI, "repo", "r", "", "URL format to snapshot repo, see moobk help repo for detail")
	cmd.MarkFlagRequired("repo")
}

func (r repoFlags) Repo() (ret moodrvs.Runner, err error) {
	return moodrvs.GetRunner(r.URI, strings.ToLower(r.Type))
}
