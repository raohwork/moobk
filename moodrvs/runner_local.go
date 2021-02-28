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

package moodrvs

import (
	"errors"
	"net/url"
)

type localRunner struct {
	COW
	backupPath string
}

func (l *localRunner) BackupPath() (ret string) { return l.backupPath }
func (l *localRunner) RunnerName() (ret string) { return "local" }

// run everything defined by COW, directly on local machine
func newLocalRunner(uri *url.URL, fs string, opts url.Values) (ret Runner, err error) {
	cow, ok := GetCOW(fs, opts)
	if !ok {
		err = errors.New("unsupported fs type: " + fs)
		return
	}
	if !cow.SetRepo(uri.Path) {
		err = errors.New("cannot set repo path")
		return
	}

	x := &localRunner{COW: cow, backupPath: uri.Path}
	ret = x
	return
}

func init() {
	addRunner("local", newLocalRunner)
}
