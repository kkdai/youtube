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
	outputVideoIds    = "video-ids"
	outputFormatPlain = "plain"
	outputFormatJSON  = "json"
	outputFormatXML   = "xml"
)

var outputFormats = []string{outputFormatPlain, outputFormatJSON, outputFormatXML, outputVideoIds}

func addFormatFlag(flagSet *pflag.FlagSet) {
	flagSet.StringVarP(&outputFormat, "format", "f", outputFormatPlain, "The output format ("+strings.Join(outputFormats, "/")+")")
}

func checkOutputFormat() error {
	for i := range outputFormats {
		if outputFormats[i] == outputFormat {
			return nil
		}
	}

	return errInvalidFormat(outputFormat)
}

type outputWriter func(w io.Writer)

func writeOutput(w io.Writer, v interface{}, writers map[string]outputWriter) error {
	switch outputFormat {
	case outputVideoIds:
		fallthrough
	case outputFormatPlain:
		writer, ok := writers[outputFormat]
		if !ok {
			return errInvalidFormat(outputFormat)
		}
		writer(w)
		return nil
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
