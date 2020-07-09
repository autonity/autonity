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

package config

type ProposerPolicy uint64

const (
	RoundRobin ProposerPolicy = iota
	WeightedRandomSampling
)

type Config struct {
	BlockPeriod    uint64         `toml:",omitempty" json:"block-period"` // Default minimum difference between two consecutive block's timestamps in second
	ProposerPolicy ProposerPolicy `toml:",omitempty" json:"policy"`       // The policy for proposer selection
}

func (c *Config) String() string {
	return "tendermint"
}

func DefaultConfig() *Config {
	return &Config{
		BlockPeriod:    1,
		ProposerPolicy: WeightedRandomSampling,
	}
}

func RoundRobinConfig() *Config {
	return &Config{
		BlockPeriod:    1,
		ProposerPolicy: RoundRobin,
	}
}
