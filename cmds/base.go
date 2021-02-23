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
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/raohwork/moobk/moodrvs"
)

var available = map[string]cmd{}

type cmd interface {
	Name() string
	Desc() string
	Help()
	Exec(args []string) int
}

func register(c cmd) {
	n := c.Name()
	l := strings.ToLower(n)
	x, has := available[l]
	if has {
		log.Fatalf("command %s has been registered for %+v (%T)", n, x, x)
	}
	available[l] = c
}

type globalOptions struct {
	args []string
	repo moodrvs.Runner
}

var flagSet *flag.FlagSet

func init() {
	flagSet = flag.NewFlagSet("", flag.ContinueOnError)
	flagSet.Usage = func() {}
	flagSet.SetOutput(ioutil.Discard)
	return
}

func g(args []string) (ret *globalOptions, err error) {
	var fs, repo string
	flagSet.StringVar(&fs, "t", "", "")
	flagSet.StringVar(&repo, "r", "", "")
	flagSet.Parse(args)

	if fs == "" {
		fs = os.Getenv("MOOBK_FS")
	}
	if fs == "" {
		fs = "btrfs"
	}

	if repo == "" {
		repo = os.Getenv("MOOBK_REPO")
	}
	if repo == "" {
		err = errors.New("-r is required")
		return
	}

	r, err := moodrvs.GetRunner(repo, fs)
	if err != nil {
		return
	}

	ret = &globalOptions{
		args: flagSet.Args(),
		repo: r,
	}
	return
}

func want(n uint, args []string) (ok bool) {
	if uint(len(args)) < n {
		return
	}

	return true
}

func mustG(args []string) (ret *globalOptions) {
	ret, err := g(args)
	if err != nil {
		available["help"].Exec([]string{})
		fmt.Println()
		fmt.Println("error: ", err)
		os.Exit(1)
	}

	return
}

func Run(args []string) {
	if len(args) < 1 {
		available["help"].Exec([]string{})
		os.Exit(1)
	}

	c := strings.ToLower(args[0])
	args = args[1:]
	x, has := available[c]
	if !has {
		fmt.Println("Command ", c, " is not supported")
		fmt.Println()
		x = available["help"]
	}
	os.Exit(x.Exec(args))
}
