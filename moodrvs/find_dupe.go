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
func FindDupe(l, r []Snapshot) (dupe, orphanL, orphanR [][]Snapshot) {
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
		d, ol, or := findDupe(s, rGroup[k])
		if len(d) > 0 {
			dupe = append(dupe, d)
		}
		if len(ol) > 0 {
			orphanL = append(orphanL, ol)
		}
		if len(or) > 0 {
			orphanR = append(orphanR, or)
		}
	}

	return
}

func findDupe(ls, rs []Snapshot) (dupe, orphanL, orphanR []Snapshot) {
	if len(rs) < 1 {
		return
	}

	ol := make([]Snapshot, 0, len(ls))
	or := make([]Snapshot, 0, len(rs))

	SnapshotSliceAsc(ls).Sort()
	SnapshotSliceAsc(rs).Sort()
	for _, l := range ls {
		lts := l.CreatedAt.Unix()
		for len(rs) > 0 {
			rts := rs[0].CreatedAt.Unix()
			if lts == rts {
				dupe = append(dupe, l)
				rs = rs[1:]
				break
			}

			if lts < rts {
				ol = append(ol, l)
				break
			}

			or = append(or, rs[0])
			rs = rs[1:]
		}
	}

	// remove latest
	if l := len(dupe); l > 0 {
		dupe = dupe[:l-1]
		orphanL = ol
		orphanR = or
	}
	return
}
