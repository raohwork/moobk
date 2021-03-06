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
	"bytes"
	"errors"
	"io"
	"net/url"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// see methods of zfsr for description
type zfs struct {
	program
	backupPath string
}

func (b *zfs) DriverName() (ret string) { return "zfs" }
func (b *zfs) Repo() (ret string)       { return b.backupPath }
func (b *zfs) SetRepo(path string) (ok bool) {
	b.backupPath = strings.TrimLeft(path, "/")
	return true
}

func (b *zfs) Snapshots() (ret []Snapshot, err error) {
	buf, err := b.basicRun("list", "-H", "-d", "1", "-t", "snapshot", b.backupPath)
	if err != nil {
		return
	}
	lines := bytes.Split(buf, []byte("\n"))

	ret = make([]Snapshot, 0, len(lines))
	for _, l := range lines {
		arr := bytes.Split(l, []byte("\t"))
		x := bytes.Split(arr[0], []byte("@"))
		if len(x) != 2 {
			// empty
			return
		}
		s, ok := ParseSnapshot(string(x[1]))
		if ok {
			ret = append(ret, s)
		}
	}

	return
}

func (b *zfs) Test(path string) (yes bool, err error) {
	if path != b.backupPath {
		err = errors.New("ZFS does not support putting snapshot under different filesystem")
		return
	}

	buf, err := b.basicRun("get", "-H", "type", path)
	if err != nil {
		var e *exec.ExitError
		if errors.As(err, &e) {
			return false, nil
		}
		return
	}

	arr := bytes.Split(buf, []byte("\t"))
	if len(arr) < 4 {
		err = errors.New("invalid output; " + string(buf))
		return
	}

	return (string(arr[2]) == "filesystem"), nil
}

func (b *zfs) Create(path, name string) (ret Snapshot, err error) {
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
	dest := b.backupPath + "@" + ret.RealName()
	_, err = b.basicRun("snap", dest)
	return
}

func (b *zfs) Delete(s Snapshot) (err error) {
	dest := b.backupPath + "@" + s.RealName()
	_, err = b.basicRun("destroy", dest)
	return
}

func (b *zfs) Send(base, s Snapshot, w io.Writer) (err error) {
	dest := b.backupPath + "@" + s.RealName()
	if base.EqualTo(s) {
		// first time
		cmd, r, e := b.forSend("send", dest)
		if e != nil {
			return e
		}
		return sendHelper(cmd, w, r)
	}

	from := b.backupPath + "@" + base.RealName()
	cmd, r, err := b.forSend("send", "-i", from, dest)
	if err != nil {
		return
	}
	return sendHelper(cmd, w, r)
}

func (b *zfs) Recv(s Snapshot, r io.Reader) (err error) {
	cmd := b.forRecv(r, "recv", "-duF", b.backupPath)

	return cmd.Run()
}

func init() {
	addCOW("zfs", func(opts url.Values) (ret COW) {
		bin := opts.Get("bin")
		if bin == "" {
			bin = "zfs"
		}
		ret = &zfs{
			program:    program{prog: bin},
			backupPath: "",
		}
		return
	})
}
