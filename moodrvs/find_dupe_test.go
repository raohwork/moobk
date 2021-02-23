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
	"testing"
	"time"
)

type findDupeTestCase struct {
	name   string
	l      []Snapshot
	r      []Snapshot
	expect []Snapshot
}

func (s *findDupeTestCase) cmpSlice(a, b []Snapshot) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v := range a {
		if !v.EqualTo(b[k]) {
			return false
		}
	}

	return true
}

func (s *findDupeTestCase) cmp(d []Snapshot) bool {
	return s.cmpSlice(s.expect, d)
}

func (s *findDupeTestCase) Run(t *testing.T) {
	actual := findDupe(s.l, s.r)
	if !s.cmp(actual) {
		t.Logf("expect: %+v", s.expect)
		t.Logf("actual: %+v", actual)
		t.Fatal("unexpected result")
	}
}

func TestFindDupe(t *testing.T) {
	// not testing empty l since it's prevented by FindDupe()

	n := "name"
	tz := time.FixedZone("TEST", 0)
	d := make([]Snapshot, 5)
	for i := 0; i < 5; i++ {
		d[i] = Snapshot{
			Name:      n,
			CreatedAt: time.Date(2021, 1, 1, 0, 0, i, 0, tz),
		}
	}
	cases := []findDupeTestCase{
		{
			name:   "empty_r",
			l:      []Snapshot{d[0]},
			r:      []Snapshot{},
			expect: []Snapshot{},
		},
		{
			name:   "l_newer",
			l:      []Snapshot{d[0], d[1], d[2]},
			r:      []Snapshot{d[0], d[1]},
			expect: []Snapshot{d[0]},
		},
		{
			name:   "r_newer",
			l:      []Snapshot{d[0], d[1]},
			r:      []Snapshot{d[0], d[1], d[2]},
			expect: []Snapshot{d[0]},
		},
		{
			name:   "synced-one",
			l:      []Snapshot{d[0], d[1]},
			r:      []Snapshot{d[1]},
			expect: []Snapshot{},
		},
		{
			name:   "synced",
			l:      []Snapshot{d[0], d[1]},
			r:      []Snapshot{d[0], d[1]},
			expect: []Snapshot{d[0]},
		},
		{
			name:   "synced-only",
			l:      []Snapshot{d[1]},
			r:      []Snapshot{d[1]},
			expect: []Snapshot{},
		},
		{
			name:   "findDupeerent_r_newer",
			l:      []Snapshot{d[0]},
			r:      []Snapshot{d[1]},
			expect: []Snapshot{},
		},
		{
			name:   "findDupeerent_l_newer",
			l:      []Snapshot{d[1]},
			r:      []Snapshot{d[0]},
			expect: []Snapshot{},
		},
		{
			name:   "chaos_zebra_l_newer",
			l:      []Snapshot{d[0], d[2], d[4]},
			r:      []Snapshot{d[1], d[3]},
			expect: []Snapshot{},
		},
		{
			name:   "chaos_zebra_r_newer",
			l:      []Snapshot{d[1], d[3]},
			r:      []Snapshot{d[0], d[2], d[4]},
			expect: []Snapshot{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, c.Run)
	}
}
