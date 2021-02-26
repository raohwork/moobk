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
	"fmt"
	"log"
	"net/url"
	"strings"
)

// Runner defines how to run tools the driver needed.
//
// It is a abstract layer above COW interface, to unify codes accessing local/remote
// together. For a "serious" runner, it should provide some degree of security.
//
// Same philosophy applies, PRs to "hipster" runners are also welcome, like
//
//     - Remote execute by commanding robot to press some buttons on a RaspberryPi,
//       and sends data by normal tcp link (you have to describe how to build the
//       robot and needed helper promgram on RPi)
//     - Plain HTTP (no S), since it's extremely dangerous.
//     - Access computer on Mars by engraving qr code on surface of a shuttle, shot
//       it to Mars, and scan the qrcode, if you're Elon Musk.
type Runner interface {
	// To compitable with COW interface
	COW
	// repository path
	BackupPath() string
	// name of this runner like local/ssh/ssh+sudo
	RunnerName() string
}

var availableRunner = map[string]func(*url.URL, string) (Runner, error){}

func addRunner(n string, f func(*url.URL, string) (Runner, error)) {
	n = strings.ToLower(n)
	_, ok := availableRunner[n]
	if ok {
		log.Fatalf("runner type %s has been registered", n)
	}
	availableRunner[n] = f
}

// alias support
var aliases = map[string]string{}

// SetAlias sets repo alias
func SetAlias(a map[string]string) {
	if a != nil {
		aliases = a
	}
}

// GetRunner creates a runner by providing driver and repo specs
func GetRunner(uri, fs string) (ret Runner, err error) {
	u, err := url.Parse(uri)
	if err == nil && u.Scheme != "" {
		return getRunner(u, fs)
	}

	// maybe alias?
	real, ok := aliases[uri]
	if !ok {
		// not alias
		err = errors.New("unsupported repo: " + uri)
		return
	}

	u, err = url.Parse(real)
	if err != nil {
		return
	}

	return getRunner(u, fs)
}

func getRunner(u *url.URL, fs string) (ret Runner, err error) {
	s := strings.ToLower(u.Scheme)
	x, ok := availableRunner[s]
	if !ok {
		err = fmt.Errorf("unsupported runner type: %s", s)
		return
	}

	return x(u, fs)
}
