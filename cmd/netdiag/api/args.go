package api

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/autonity/autonity/cmd/netdiag/strats"
)

type Argument interface {
	AskUserInput() error
}

type ArgTarget struct {
	Target int
}

func (a *ArgTarget) AskUserInput() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter target peer index, 0 for all:")
	input, _ := reader.ReadString('\n')
	targetIndex, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		fmt.Println("Invalid target index.")
		return err
	}
	a.Target = targetIndex
	return nil
}

type ArgSize struct {
	Size int
}

func (a *ArgSize) AskUserInput() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter size (kB) - max 15000: ")
	input, _ := reader.ReadString('\n')
	size, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		fmt.Println("Invalid size.")
		return err
	}
	a.Size = size * 1000
	return nil
}

type ArgCount struct {
	PacketCount int
}

func (a *ArgCount) AskUserInput() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter number of DevP2P packets: ")
	input, _ := reader.ReadString('\n')
	packetCount, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		fmt.Println("Invalid number.")
		return err
	}
	a.PacketCount = packetCount
	return nil
}

type ArgTargetSize struct {
	ArgTarget
	ArgSize
}

func (a *ArgTargetSize) AskUserInput() error {
	if err := a.ArgTarget.AskUserInput(); err != nil {
		return err
	}
	return a.ArgSize.AskUserInput()
}

type ArgSizeCount struct {
	ArgSize
	ArgCount
}

func (a *ArgSizeCount) AskUserInput() error {
	if err := a.ArgSize.AskUserInput(); err != nil {
		return err
	}
	return a.ArgCount.AskUserInput()
}

type ArgTargetSizeCount struct {
	ArgTarget
	ArgSize
	ArgCount
}

func (a *ArgTargetSizeCount) AskUserInput() error {
	if err := a.ArgTarget.AskUserInput(); err != nil {
		return err
	}
	if err := a.ArgCount.AskUserInput(); err != nil {
		return err
	}
	return a.ArgSize.AskUserInput()
}

type ArgEmpty struct {
}

type ArgStrategy struct {
	Strategy int
}

func (a *ArgStrategy) AskUserInput() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Available Dissemination Strategies: ")
	for i, s := range strats.StrategyRegistry {
		fmt.Printf("[%d] %s\n", i, s.Name)
	}
	fmt.Print("Chose strategy: ")
	input, _ := reader.ReadString('\n')
	index, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || index >= len(strats.StrategyRegistry) || index < 0 {
		fmt.Println("Invalid strategy selected")
		return err
	}
	a.Strategy = index
	return nil
}

type ArgMaxPeers struct {
	MaxPeers int
}

func (a *ArgMaxPeers) AskUserInput() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Number of peers to disseminate, 0 for all: ")
	input, _ := reader.ReadString('\n')
	count, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		fmt.Println("Invalid number")
		return err
	}
	a.MaxPeers = count
	return nil
}

type ArgDisseminate struct {
	ArgStrategy
	ArgMaxPeers
	ArgSize
	ArgOutputFile
}

func (a *ArgDisseminate) AskUserInput() error {
	if err := a.ArgStrategy.AskUserInput(); err != nil {
		return err
	}
	if err := a.ArgMaxPeers.AskUserInput(); err != nil {
		return err
	}
	if err := a.ArgSize.AskUserInput(); err != nil {
		return err
	}
	if err := a.ArgOutputFile.AskUserInput(); err != nil {
		return err
	}
	return nil
}

type ArgWarmUp struct {
	ArgMaxPeers
	ArgSize
}

func (a *ArgWarmUp) AskUserInput() error {
	if err := a.ArgMaxPeers.AskUserInput(); err != nil {
		return err
	}
	if err := a.ArgSize.AskUserInput(); err != nil {
		return err
	}
	return nil
}

type ArgGraphConstruct struct {
	ArgStrategy
	ArgMaxPeers
}

func (a *ArgGraphConstruct) AskUserInput() error {
	if err := a.ArgStrategy.AskUserInput(); err != nil {
		return err
	}
	if err := a.ArgMaxPeers.AskUserInput(); err != nil {
		return err
	}
	return nil
}

type ArgOutputFile struct {
	OutputFile string
}

func (a *ArgOutputFile) AskUserInput() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Output filename: ")
	input, _ := reader.ReadString('\n')
	a.OutputFile = strings.TrimSpace(input)
	return nil
}
