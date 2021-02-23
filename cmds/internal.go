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
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/raohwork/moobk/moodrvs"
)

type internal struct{}

func init() { register(internal{}) }

func (_ internal) Name() string { return "internal" }
func (_ internal) Desc() string { return "internal wrapper to use with ssh and ssh+sudo" }
func (_ internal) Help() {
	fmt.Print(`internal sub-command

internal, as the name shows, is for internal use only.

It wraps fs driver to ensure only hard-coded actions are performed.

End-users *SHOULD NOT* use this directly.
`)
}

func (i internal) WriteJson(data interface{}, err error) (code int) {
	if err != nil {
		return i.WriteErr(err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.Encode(data)
	return
}

func (_ internal) WriteErr(err error) (code int) {
	fmt.Fprint(os.Stderr, err)
	return 1
}

func (i internal) Exec(args []string) (ret int) {
	opt := mustG(args)
	if !want(1, opt.args) {
		i.Help()
		return 1
	}

	switch strings.ToLower(opt.args[0]) {
	case "snapshots":
		return i.WriteJson(opt.repo.Snapshots())
	case "test":
		if !want(2, opt.args) {
			return i.WriteJson(nil, errors.New("missing path"))
		}
		return i.WriteJson(opt.repo.Test(opt.args[1]))
	case "create":
		if !want(2, opt.args) {
			return i.WriteJson(nil, errors.New("missing path"))
		}
		name := ""
		if len(opt.args) > 2 {
			name = opt.args[2]
		}
		return i.WriteJson(opt.repo.Create(opt.args[1], name))
	case "delete":
		if !want(2, opt.args) {
			return i.WriteJson(nil, errors.New("missing snapsoht"))
		}
		s, ok := moodrvs.ParseSnapshot(opt.args[1])
		if !ok {
			return i.WriteJson(nil, errors.New("incorrect snapshot"))
		}

		return i.WriteJson(nil, opt.repo.Delete(s))
	case "send":
		if !want(3, opt.args) {
			return i.WriteErr(errors.New("need base/s snapshot name"))
		}
		base, ok := moodrvs.ParseSnapshot(opt.args[1])
		if !ok {
			return i.WriteErr(errors.New("invalid base snapshot"))
		}
		s, ok := moodrvs.ParseSnapshot(opt.args[2])
		if !ok {
			return i.WriteErr(errors.New("invalid s snapshot"))
		}
		err := opt.repo.Send(base, s, os.Stdout)
		if err != nil {
			return i.WriteErr(err)
		}
		return
	case "recv":
		if !want(2, opt.args) {
			return i.WriteErr(errors.New("need snapshot name"))
		}
		s, ok := moodrvs.ParseSnapshot(opt.args[1])
		if !ok {
			return i.WriteErr(errors.New("invalid s snapshot"))
		}
		err := opt.repo.Recv(s, os.Stdin)
		if err != nil {
			return i.WriteErr(err)
		}
		return
	}

	return i.WriteErr(errors.New("invalid action"))
}
