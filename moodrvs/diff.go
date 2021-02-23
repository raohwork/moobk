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

type SnapshotDiff struct {
	Missing []Snapshot
	Base    *Snapshot
}

// Diff finds snapshots exists in l but not in r
func Diff(l, r []Snapshot) (ret []SnapshotDiff) {
	if len(l) < 1 {
		return
	}
	// first, group by name
	lGroup := map[string][]Snapshot{}
	rGroup := map[string][]Snapshot{}

	for _, s := range l {
		lGroup[s.Name] = append(lGroup[s.Name], s)
	}
	for _, s := range r {
		rGroup[s.Name] = append(rGroup[s.Name], s)
	}

	for k, s := range lGroup {
		ret = append(ret, diff(s, rGroup[k]))
	}
	return
}

func diff(l, r []Snapshot) (ret SnapshotDiff) {
	SnapshotSliceAsc(l).Sort()
	size := len(r)
	if size < 1 {
		return SnapshotDiff{
			Missing: l,
		}
	}

	SnapshotSliceAsc(r).Sort()

	for idx, s := range l {
		if size < 1 {
			ret.Missing = l[idx:]
			return
		}

		if s.CreatedAt.Before(r[0].CreatedAt) {
			continue
		}

		for size > 0 && s.CreatedAt.After(r[0].CreatedAt) {
			r = r[1:]
			size--
		}
		if size < 1 {
			ret.Missing = l[idx:]
			return
		}
		if s.CreatedAt.Equal(r[0].CreatedAt) {
			r = r[1:]
			size--
			ret.Base = &l[idx]
			continue
		}
	}

	// remote has more than local, nothing to do
	return
}
