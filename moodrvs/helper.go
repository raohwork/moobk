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
)

// helper
type program struct {
	prog string
}

func (p *program) exec(args ...string) (ret *exec.Cmd) {
	return exec.Command(p.prog, args...)
}

func (p *program) basicRun(args ...string) (data []byte, err error) {
	return p.exec(args...).CombinedOutput()
}

func (p *program) forRecv(r io.Reader, args ...string) (ret *exec.Cmd) {
	cmd := p.exec(args...)
	cmd.Stdin = r
	return cmd
}

func (p *program) forSend(args ...string) (ret *exec.Cmd, r io.Reader, err error) {
	ret = p.exec(args...)
	r, err = ret.StdoutPipe()
	return
}
