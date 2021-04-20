package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/pflag"
)

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
	for i := range outputFormats {
		if outputFormats[i] == outputFormat {
			return nil
		}
	}

	return fmt.Errorf("output format %s is not valid", outputFormat)
}

func writeStructuredOutput(w io.Writer, v interface{}) error {
	switch outputFormat {
	case outputFormatXML:
		return xml.NewEncoder(w).Encode(v)
	case outputFormatJSON:
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		return encoder.Encode(v)
	default:
		panic("invalid format: " + outputFormat)
	}
}
