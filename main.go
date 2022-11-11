package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu/sensu-plugin-sdk/sensu"

	"github.com/sardinasystems/sensu-go-ipmi-power-check/ipmimon"
)

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-go-ipmi-power-check",
			Short:    "plugin to check power supply",
			Keyspace: "sensu.io/plugins/sensu-go-ipmi-power-check/config",
		},
	}

	options = []sensu.ConfigOption{}
)

func main() {
	useStdin := false
	fi, err := os.Stdin.Stat()
	if err != nil {
		fmt.Printf("Error check stdin: %v\n", err)
	}
	// Check the Mode bitmask for Named Pipe to indicate stdin is connected
	if fi.Mode()&os.ModeNamedPipe != 0 {
		useStdin = true
	}

	check := sensu.NewGoCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, useStdin)
	check.Execute()
}

func checkArgs(event *corev2.Event) (int, error) {
	return sensu.CheckStateOK, nil
}

func executeCheck(event *corev2.Event) (int, error) {
	ctx := context.TODO()

	report, err := ipmimon.GetReport(ctx)
	if err != nil {
		return 0, err
	}

	puReport := report.Type(ipmimon.TypePowerUnit)
	psReport := report.Type(ipmimon.TypePowerSupply)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Type", "State", "Reading", "Units", "Event"})

	state := sensu.CheckStateOK

	checkItem := func(it *ipmimon.Item) {
		itSt := mapState(it)
		if itSt > state {
			state = itSt
		}

		table.Append([]string{strconv.Itoa(it.ID), it.Name, it.Type, it.State, it.Reading, it.Units, it.Event})
	}

	for _, it := range puReport {
		checkItem(&it)
	}
	for _, it := range psReport {
		checkItem(&it)
	}

	switch state {
	case sensu.CheckStateOK:
		fmt.Println("OK")

	case sensu.CheckStateWarning:
		fmt.Println("WARNING")

	case sensu.CheckStateCritical:
		fmt.Println("CRITICAL")

	default:
		fmt.Println("UNKNOWN")
	}

	fmt.Println()
	table.Render()

	return state, nil
}

func mapState(it *ipmimon.Item) int {
	switch it.State {
	case ipmimon.StateNominal:
		return sensu.CheckStateOK

	case ipmimon.StateWarning:
		return sensu.CheckStateWarning

	case ipmimon.StateCritical:
		return sensu.CheckStateCritical

	// skip N/A
	default:
		return sensu.CheckStateOK
	}
}
