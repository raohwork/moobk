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
	"os/exec"
	"path/filepath"
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

// use btrfs subvolume show to see if it is a btrfs subvolume
func (b *btrfs) Test(path string) (yes bool, err error) {
	_, err = b.basicRun("sub", "show", path)
	if err == nil {
		return true, nil
	}

	if _, ok := err.(*exec.ExitError); ok {
		return false, nil
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
		if err = cmd.Start(); err != nil {
			return
		}
		_, err = io.Copy(w, r)
		if err != nil {
			return
		}
		return cmd.Wait()
	}

	from := b.backupPath + "/" + base.RealName()
	cmd, r, err := b.forSend("send", "-c", from, dest)
	if err != nil {
		return
	}
	if err = cmd.Start(); err != nil {
		return
	}
	_, err = io.Copy(w, r)
	if err != nil {
		return
	}
	return cmd.Wait()
}

// btrfs receive, stay cool
func (b *btrfs) Recv(s Snapshot, r io.Reader) (err error) {
	cmd := b.forRecv(r, "receive", b.backupPath)

	return cmd.Run()
}

func init() {
	addCOW("btrfs", func() (ret COW) {
		ret = &btrfs{
			program:    program{prog: "btrfs"},
			backupPath: "",
		}
		return
	})
}
