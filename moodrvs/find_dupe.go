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

// FindDupe finds all snapsnhots exist in both l and r
func FindDupe(l, r []Snapshot) (ret [][]Snapshot) {
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
		x := findDupe(s, rGroup[k])
		if len(x) > 0 {
			ret = append(ret, x)
		}
	}

	return
}

func findDupe(l, r []Snapshot) (ret []Snapshot) {
	if len(r) < 1 {
		return
	}

	SnapshotSliceAsc(r).Sort()
	for _, s := range l {
		for _, x := range r {
			if s.CreatedAt.Before(x.CreatedAt) {
				break
			}

			if s.CreatedAt.After(x.CreatedAt) {
				continue
			}

			ret = append(ret, s)
			break
		}
	}

	// remove latest
	if l := len(ret); l > 0 {
		ret = ret[:l-1]
	}
	return
}
