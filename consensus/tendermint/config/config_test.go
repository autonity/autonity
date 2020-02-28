package config

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg == nil {
		t.Fatal("got a nil default config")
	}

	if cfg.ProposerPolicy != RoundRobin {
		t.Fatal("default config is not RoundRobin")
	}
}

func TestDefaultConfigSetPolicy(t *testing.T) {
	cfg := DefaultConfig()
	if cfg == nil {
		t.Fatal("got a nil default config")
	}
	cfg.SetProposerPolicy(Sticky)

	if cfg.ProposerPolicy != Sticky {
		t.Fatal("default config is not changed")
	}
}

func TestDefaultConfigGetPolicy(t *testing.T) {
	cfg := DefaultConfig()
	if cfg == nil {
		t.Fatal("got a nil default config")
	}

	if cfg.GetProposerPolicy() != RoundRobin {
		t.Fatal("default config is not RoundRobin")
	}
}
