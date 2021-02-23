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
	"net/url"
	"reflect"
	"testing"
)

func TestBuildSSHOptFull(t *testing.T) {
	val, _ := url.Parse("ssh://user@example.com:1234/path?ssh_4&ssh_A&ssh_F=/path/config&ssh_i=/path/priv.key&ssh_o=opt1&ssh_o=opt2")
	actual := buildSSHOpt(val)
	base := []string{"-l", "user", "-p", "1234"}
	l := len(base) + 1 + 10 // -4 -A -F /path... -i /path... -o opt1 -o opt2

	if x := len(actual); x != l {
		t.Log("expected length: ", l)
		t.Log("actual: ", actual)
		t.Fatal("unexpected length")
	}

	if !reflect.DeepEqual(base, actual[:4]) {
		t.Log("expected begin: -l user -p 1234")
		t.Log("actual: ", actual)
		t.Fatal("unexpected length")
	}
	if actual[len(actual)-1] != "example.com" {
		t.Log("expected end: example.com")
		t.Log("actual: ", actual)
		t.Fatal("unexpected length")
	}

	f1 := func(s string) {
		for _, v := range actual {
			if v == s {
				return
			}
		}

		t.Log("expected ", s, "exists")
		t.Log("actual: ", actual)
		t.Fatal("unexpected length")
	}
	f2 := func(s1, s2 string) {
		l := len(actual)
		for idx, v := range actual {
			if v == s1 {
				if idx >= l-1 {
					t.Log("expected ", s1, " ", s2, "exists")
					t.Log("actual: ", actual)
					t.Fatal("unexpected length")
				}
				if actual[idx+1] == s2 {
					return
				}
			}
		}

		t.Log("expected ", s1, " ", s2, "exists")
		t.Log("actual: ", actual)
		t.Fatal("unexpected length")
	}

	f1("-4")
	f1("-A")
	f2("-o", "opt1")
	f2("-o", "opt2")
	f2("-i", "/path/priv.key")
	f2("-F", "/path/config")
}
