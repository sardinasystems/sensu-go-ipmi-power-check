// Package ipmimon provides interface to IPMI monitoing
package ipmimon

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"os/exec"
	"strings"
	"syscall"

	"github.com/gocarina/gocsv"
)

const (
	TypePowerSupply = "Power Supply"
	TypePowerUnit   = "Power Unit"
)

const (
	StateNominal  = "Nominal"
	StateWarning  = "Warning"
	StateCritical = "Critical"
)

// Item represents one line from the report
//
//	ipmimonitoring --comma-separated-output | head -n1
//	ID,Name,Type,State,Reading,Units,Event
type Item struct {
	ID      int    `csv:"ID" json:"id"`
	Name    string `csv:"Name" json:"name"`
	Type    string `csv:"Type" json:"type"`
	State   string `csv:"State" json:"state"`
	Reading string `csv:"Reading" json:"reading"`
	Units   string `csv:"Units" json:"units"`
	Event   string `csv:"Event" json:"event"`
}

// Report represent parsed report
type Report []Item

// ParseCSV parses report csv into Report
func ParseCSV(data []byte) (Report, error) {
	items := make([]Item, 0)

	// skip any messages before report table
	idx := bytes.Index(data, []byte("ID,"))
	if idx > 0 {
		data = data[idx:]
	}

	err := gocsv.UnmarshalBytes(data, &items)
	if err != nil {
		return nil, fmt.Errorf("csv parse error: %w", err)
	}

	return Report(items), nil
}

// GetReport run ipmimonitoring and parse report
func GetReport(ctx context.Context) (Report, error) {
	cmd := exec.CommandContext(ctx, "ipmimonitoring", "--comma-separated-output")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGKILL,
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("ipmimonitoring command failed: %w", err)
	}

	return ParseCSV(out)
}

func (report Report) Filter(testFunc func(*Item) bool) Report {
	items := make([]Item, 0)

	for _, it := range report {
		if testFunc(&it) {
			items = append(items, it)
		}
	}

	return Report(items)
}

func (report Report) Type(typ string) Report {
	return report.Filter(func(it *Item) bool {
		return it.Type == typ
	})
}

func (it *Item) Events() []string {
	r := csv.NewReader(strings.NewReader(strings.ReplaceAll(it.Event, "'", "\"")))
	r.Comma = ' '
	fields, err := r.Read()
	if err != nil {
		panic(err)
	}

	return fields
}
