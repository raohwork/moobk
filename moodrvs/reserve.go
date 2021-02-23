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
	"strconv"
	"strings"
	"time"
)

// Reserve removes some items from data if rules applies
func Reserve(data []Snapshot, cnt uint, t string) (ret []Snapshot) {
	return reserveTime(reserveCount(data, cnt), t)
}

func reserveTime(data []Snapshot, t string) (ret []Snapshot) {
	l := len(t)
	if t == "" || l < 2 {
		return data
	}

	now := time.Now()
	unit := strings.ToLower(t[l-1:])
	istr := t[:l-1]
	i, err := strconv.Atoi(istr)
	if err != nil {
		// format error
		return data
	}
	i = -i
	switch unit {
	case "h":
		now = now.Add(time.Duration(i) * time.Hour)
	case "d":
		now = now.AddDate(0, 0, int(i))
	case "w":
		now = now.AddDate(0, 0, int(i)*7)
	case "m":
		now = now.AddDate(0, int(i), 0)
	default:
		// format error
		return data
	}

	ret = make([]Snapshot, 0, len(data))
	for _, s := range data {
		if s.CreatedAt.Before(now) {
			ret = append(ret, s)
		}
	}

	return ret
}

func reserveCount(data []Snapshot, cnt uint) (ret []Snapshot) {
	if cnt == 0 {
		return data
	}
	if l := uint(len(data)); l < cnt {
		return data
	}
	return data[cnt:]
}
