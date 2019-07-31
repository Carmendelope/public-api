package output

import (
	"encoding/json"
	"fmt"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strings"
)

// Output structure containing the definition of how the output of the cli must be provided.
type Output struct {
	format string
	// Label Length specifies the maximum length that will be displayed from the labels
	labelLength int
}

func NewOutput(format string, labelLength int) *Output{
	return &Output{format, labelLength}
}

// asJSON returns true if the output must be shown in raw JSON format
func (o * Output) asJSON() bool{
	return strings.ToLower(o.format) == "json" || strings.ToLower(o.format) == "raw"
}

// asTable returns true if the output must be shown in a table
func (o * Output) asTable() bool {
	if strings.ToLower(o.format) == "text"{
		log.Warn().Msg("output set to text is deprecated, use table instead. This option will be deprecated in 0.5.0")
	}
	return strings.ToLower(o.format) == "table"
}

// PrintResultOrError prints in the given output format/method the result of the operation of the error.
func (o *Output) PrintResultOrError(result interface{}, err error, errMsg string) {
	if err != nil {
		converted := conversions.ToDerror(err)
		if zerolog.GlobalLevel() == zerolog.DebugLevel{
			log.Fatal().Str("trace", conversions.ToDerror(err).DebugReport()).Msg(errMsg)
		}else{
			log.Fatal().Str("err", converted.Error()).Msg(errMsg)
		}
	} else {
		if o.asTable(){
			o.PrintResultAsTable(result)
		}else if o.asJSON(){
			o.PrintResultAsJSON(result)
		}else{
			log.Warn().Str("format", o.format).Msg("Invalid output method, defaulting to JSON")
			o.PrintResultAsJSON(result)
		}
	}
}

// PrintResultAsTable transforms the result into a table format and prints it to stdout.
func (o * Output) PrintResultAsTable(result interface{}) {
	table := AsTable(result, o.labelLength)
	table.Print()
}

// PrintResultAsJSON prints the raw JSON result to stdout.
func (o *Output) PrintResultAsJSON(result interface{}) error {
	res, err := json.MarshalIndent(result, "", "  ")
	if err == nil {
		fmt.Println(string(res))
	}
	return err
}
