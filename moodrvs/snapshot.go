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
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Snapshot represents a snapshot
type Snapshot struct {
	Name      string
	CreatedAt time.Time
}

// RealName combines all together
func (s Snapshot) RealName() (ret string) {
	return fmt.Sprintf("%s-%d", s.Name, s.CreatedAt.Unix())
}

// EqualTo compares if n is identical to s
func (s Snapshot) EqualTo(n Snapshot) (yes bool) {
	return s.Name == n.Name && s.CreatedAt.Unix() == n.CreatedAt.Unix()
}

// PArseSnapshot does reverse procedure of Snapshot.RealName()
func ParseSnapshot(name string) (ret Snapshot, ok bool) {
	arr := strings.Split(name, "-")
	l := len(arr)
	if l < 2 {
		return
	}
	tsStr := arr[l-1]
	ts, err := strconv.ParseInt(tsStr, 10, 64)
	if err != nil {
		return
	}
	ret.CreatedAt = time.Unix(ts, 0)
	ret.Name = filepath.Base(strings.Join(arr[:l-1], "-"))
	ok = true
	return
}

// SnapshotSlice sorts snapshots from new to old
type SnapshotSlice []Snapshot

func (s SnapshotSlice) Len() int           { return len(s) }
func (s SnapshotSlice) Less(i, j int) bool { return s[i].CreatedAt.After(s[j].CreatedAt) }
func (s SnapshotSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SnapshotSlice) Sort()              { sort.Sort(s) }

// SnapshotSliceAsc sorts snapshots from old to new
type SnapshotSliceAsc []Snapshot

func (s SnapshotSliceAsc) Len() int           { return len(s) }
func (s SnapshotSliceAsc) Less(i, j int) bool { return s[i].CreatedAt.Before(s[j].CreatedAt) }
func (s SnapshotSliceAsc) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SnapshotSliceAsc) Sort()              { sort.Sort(s) }
