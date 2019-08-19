// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

type Step uint64

const (
	propose Step = iota
	prevote
	precommit
	precommitDone
)

func (s Step) String() string {
	if s == propose {
		return "propose"
	} else if s == prevote {
		return "prevote"
	} else if s == precommit {
		return "precommit"
	} else if s == precommitDone {
		return "precommitDone"
	} else {
		return "Unknown"
	}
}

func (s Step) Cmp(y Step) int {
	if uint64(s) < uint64(y) {
		return -1
	}
	if uint64(s) > uint64(y) {
		return 1
	}
	return 0
}
