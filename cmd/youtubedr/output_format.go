package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/pflag"
)

// newTable returns an ASCII table writer for plain output that never wraps
// cell text. align sets the alignment of each column's rows; headers are
// printed as given and centered.
func newTable(w io.Writer, align ...tw.Align) *tablewriter.Table {
	return tablewriter.NewTable(w,
		tablewriter.WithSymbols(tw.NewSymbols(tw.StyleASCII)),
		tablewriter.WithHeaderAutoFormat(tw.Off),
		tablewriter.WithHeaderAutoWrap(tw.WrapNone),
		tablewriter.WithHeaderAlignment(tw.AlignCenter),
		tablewriter.WithRowAutoWrap(tw.WrapNone),
		tablewriter.WithRowAlignmentConfig(tw.CellAlignment{PerColumn: align}),
	)
}

// the selected output Format
var outputFormat string

const (
	outputFormatPlain = "plain"
	outputFormatJSON  = "json"
	outputFormatXML   = "xml"
)

var outputFormats = []string{outputFormatPlain, outputFormatJSON, outputFormatXML}

func addFormatFlag(flagSet *pflag.FlagSet) {
	flagSet.StringVarP(&outputFormat, "format", "f", outputFormatPlain, "The output format ("+strings.Join(outputFormats, "/")+")")
}

func checkOutputFormat() error {
	if slices.Contains(outputFormats, outputFormat) {
		return nil
	}

	return errInvalidFormat(outputFormat)
}

type outputWriter func(w io.Writer) error

func writeOutput(w io.Writer, v any, plainWriter outputWriter) error {
	switch outputFormat {
	case outputFormatPlain:
		return plainWriter(w)
	case outputFormatJSON:
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		return encoder.Encode(v)
	case outputFormatXML:
		return xml.NewEncoder(w).Encode(v)
	default:
		return errInvalidFormat(outputFormat)
	}
}

type errInvalidFormat string

func (err errInvalidFormat) Error() string {
	return fmt.Sprintf("invalid output format: %s", outputFormat)
}
