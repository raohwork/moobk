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
	"encoding/json"
	"errors"
	"io"
	"net/url"
	"os/exec"
	"strings"
)

type sshRunner struct {
	execName   string
	sshExec    string
	sshOpts    []string
	prefix     []string
	name       string
	fs         string
	backupPath string
	cow        COW
}

func (r *sshRunner) BackupPath() (ret string)      { return r.backupPath }
func (r *sshRunner) RunnerName() (ret string)      { return r.name }
func (r *sshRunner) DriverName() string            { return r.cow.DriverName() }
func (r *sshRunner) Repo() string                  { return r.backupPath }
func (r *sshRunner) SetRepo(path string) (ok bool) { return false }

func (r *sshRunner) run(args ...string) (ret *exec.Cmd) {
	x := append([]string{}, r.sshOpts...)

	if len(r.prefix) > 0 {
		x = append(x, r.prefix...)
	}
	x = append(
		x,
		r.execName, "internal",
		"-t", r.fs,
		"-r", "local://"+r.backupPath,
	)
	x = append(x, args...)

	return exec.Command(r.sshExec, x...)
}

func (r *sshRunner) simpleRun(data interface{}, args ...string) (err error) {
	cmd := r.run(args...)
	stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	cmd.Stderr = stderr
	cmd.Stdout = stdout

	err = cmd.Run()
	if err != nil {
		// see if there's output in stderr
		if stderr.Len() > 0 {
			return errors.New(stderr.String())
		}
		return
	}

	if data == nil {
		return
	}
	return json.Unmarshal(stdout.Bytes(), data)
}

func (r *sshRunner) Snapshots() (ret []Snapshot, err error) {
	err = r.simpleRun(&ret, "snapshots")
	return
}

func (r *sshRunner) Test(path string) (yes bool, err error) {
	err = r.simpleRun(&yes, "test", path)
	return
}

func (r *sshRunner) Create(path, name string) (ret Snapshot, err error) {
	err = r.simpleRun(&ret, "create", path, name)
	return
}

func (r *sshRunner) Delete(s Snapshot) (err error) {
	return r.simpleRun(nil, "delete", s.RealName())
}

func (r *sshRunner) Send(base, s Snapshot, w io.Writer) (err error) {
	cmd := r.run("send", base.RealName(), s.RealName())
	cmd.Stdout = w
	return cmd.Run()
}

func (r *sshRunner) Recv(s Snapshot, rd io.Reader) (err error) {
	cmd := r.run("recv", s.RealName())
	cmd.Stdin = rd
	return cmd.Run()
}

func buildSSHOpt(uri *url.URL) (sshOpts []string) {
	if uri.User.Username() != "" {
		sshOpts = append(sshOpts, "-l", uri.User.Username())
	}
	if uri.Port() != "" {
		sshOpts = append(sshOpts, "-p", uri.Port())
	}

	val := uri.Query()
	for k, arr := range val {
		if !strings.HasPrefix(k, "ssh_") {
			continue
		}

		k = strings.Replace(k[3:], "_", "-", -1)
		if len(arr) == 0 || (len(arr) == 1 && arr[0] == "") {
			sshOpts = append(sshOpts, k)
			continue
		}

		for _, v := range arr {
			sshOpts = append(sshOpts, k, v)
		}
	}

	sshOpts = append(sshOpts, uri.Hostname())
	return
}

// use moobk as wrapper on remote machine, and use ssh to execute it.
func newSSHRunner(uri *url.URL, fs string) (ret *sshRunner, err error) {
	cow, ok := GetCOW(fs)
	if !ok {
		err = errors.New("unsupported fs type: " + fs)
		return
	}

	bin := "moobk"
	ssh := "ssh"

	val := uri.Query()
	if s := val.Get("moobk"); s != "" {
		bin = s
	}
	if s := val.Get("ssh"); s != "" {
		ssh = s
	}

	x := &sshRunner{
		execName:   bin,
		sshExec:    ssh,
		sshOpts:    buildSSHOpt(uri),
		name:       "ssh",
		fs:         fs,
		backupPath: uri.Path,
		cow:        cow,
	}
	ret = x
	return
}

func init() {
	addRunner("ssh", func(uri *url.URL, fs string) (ret Runner, err error) {
		return newSSHRunner(uri, fs)
	})
	addRunner("ssh+sudo", func(uri *url.URL, fs string) (ret Runner, err error) {
		x, err := newSSHRunner(uri, fs)
		if err != nil {
			return
		}

		x.prefix = append(x.prefix, "sudo")
		ret = x
		return
	})
}
