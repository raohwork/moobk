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
	"io"
	"net/url"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type btrfs struct {
	program
	backupPath string
}

func (b *btrfs) DriverName() (ret string)      { return "btrfs" }
func (b *btrfs) Repo() (ret string)            { return b.backupPath }
func (b *btrfs) SetRepo(path string) (ok bool) { b.backupPath = path; return true }

// just use ls, so lazy
func (b *btrfs) Snapshots() (ret []Snapshot, err error) {
	maybe, err := filepath.Glob(b.backupPath + "/*-[0-9]*")
	if err != nil {
		err = fmt.Errorf("cannot list snapshots in %s: %w", b.backupPath, err)
		return
	}

	ret = make([]Snapshot, 0, len(maybe))
	for _, p := range maybe {
		s, ok := ParseSnapshot(filepath.Base(p))
		if ok {
			ret = append(ret, s)
		}
	}

	return
}

const btrfsNotFS = "not a btrfs filesystem"
const btrfsNotFound = "no such file or directory"

// use btrfs subvolume show to see if it is a btrfs subvolume
func (b *btrfs) Test(path string) (yes bool, err error) {
	_, err = b.basicRun("sub", "show", path)
	if err == nil {
		return true, nil
	}

	var e *exec.ExitError
	if errors.As(err, &e) {
		// check error message
		str := strings.ToLower(err.Error())
		if strings.Contains(str, btrfsNotFS) {
			return false, nil
		}
		if strings.Contains(str, btrfsNotFound) {
			return false, nil
		}
	}

	return
}

// just btrfs subvolume snapshot, ez
func (b *btrfs) Create(path, name string) (ret Snapshot, err error) {
	ok, err := b.Test(path)
	if !ok {
		err = ErrUnsupportedFS(path)
		return
	}

	if name == "" {
		name = filepath.Base(path)
	}
	ret.Name = name
	ret.CreatedAt = time.Now()
	dest := b.backupPath + "/" + ret.RealName()
	_, err = b.basicRun("sub", "snap", "-r", path, dest)
	return
}

// btrfs subvolume delete, nothing to worry about
func (b *btrfs) Delete(s Snapshot) (err error) {
	dest := b.backupPath + "/" + s.RealName()
	_, err = b.basicRun("sub", "del", dest)
	return
}

// btrfs send, use incremental send if possible
func (b *btrfs) Send(base, s Snapshot, w io.Writer) (err error) {
	dest := b.backupPath + "/" + s.RealName()
	if base.EqualTo(s) {
		// first time
		cmd, r, e := b.forSend("send", dest)
		if e != nil {
			return e
		}
		return sendHelper(cmd, w, r)
	}

	from := b.backupPath + "/" + base.RealName()
	cmd, r, err := b.forSend("send", "-c", from, dest)
	if err != nil {
		return
	}
	return sendHelper(cmd, w, r)
}

// btrfs receive, stay cool
func (b *btrfs) Recv(s Snapshot, r io.Reader) (err error) {
	cmd := b.forRecv(r, "receive", b.backupPath)

	return cmd.Run()
}

func init() {
	addCOW("btrfs", func(opts url.Values) (ret COW) {
		bin := opts.Get("bin")
		if bin == "" {
			bin = "btrfs"
		}

		ret = &btrfs{
			program:    program{prog: bin},
			backupPath: "",
		}
		return
	})
}
