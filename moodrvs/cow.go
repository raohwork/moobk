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
	"io"
	"log"
	"strings"
)

// COW defines what actions the driver should support
//
// A driver is some code to do predefined jobs on something supports snapshot
// management. Driver SHOULD promise integrity and correctness (mostly by related
// fs or tool it use).
//
// Although it is designed to be used with CoWfs, the interface also applicable to
// some *hipster* usage, like
//
//     - tar some folder and manage by local filesystem tools (rm/mv/ls)
//     - diff some documents and send/receive by email (!!)
//     - dump database content and save to Amazon S3
//
// PRs to add hipster driver are ALWAYS welcome IF THE CODE IS COMPITABLE WITH GPL
// AND NEEDED RESOURCES ARE DESCRIBED (like special hardware/helper programs ...),
// but will NEVER be included in document nor official binary (locked by go build
// tags).
//
// See also Runner interface if you're intresting in writing driver! PR is welcome!
type COW interface {
	// driver name
	DriverName() string
	// get path where snapshots stores
	Repo() string
	// change to different repo, not all driver supports it.
	SetRepo(path string) (ok bool)
	// retrieve recognized snapshots
	Snapshots() (ret []Snapshot, err error)
	// test if specified path is compatible with this driver
	Test(path string) (yes bool, err error)
	// create a snapshot
	Create(path, name string) (ret Snapshot, err error)
	// delete a snapshot
	Delete(s Snapshot) (err error)
	// incremental send a snapshot to somewhere. it will block until
	// successfully sent or any error
	Send(base, s Snapshot, w io.Writer) (err error)
	// receive a snapshot from somewhere. it will block until successfully
	// received or any error
	Recv(s Snapshot, r io.Reader) (err error)
}

var availableCOW = map[string]func() COW{}

func addCOW(n string, f func() COW) {
	n = strings.ToLower(n)
	_, ok := availableCOW[n]
	if ok {
		log.Fatalf("fs type %s has been registered", n)
	}
	availableCOW[n] = f
}

// GetCOW retrieves a driver by its name
func GetCOW(fs string) (ret COW, ok bool) {
	x, ok := availableCOW[strings.ToLower(fs)]
	if !ok {
		return
	}

	return x(), ok
}

// ErrUnsupportedFS indicates path/repo is not supported by the driver
type ErrUnsupportedFS string

func (e ErrUnsupportedFS) Error() string {
	return "moodrvs: unsupported filesystem: " + string(e)
}
