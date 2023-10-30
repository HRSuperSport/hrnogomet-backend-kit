package logging

import (
	"fmt"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
	"strings"
	"time"
)

// ConfigureShortFileNameInLogMessages includes file name (e.g. 'kafka_repo.go') into every message logged by zerlogo
// By default absolute path is included which results in very long and unreadable log messages. File name only is usually sufficient
func ConfigureShortFileNameInLogMessages() {
	zerolog.CallerMarshalFunc = CallerMarshalFuncWithShortFileName
	log.Logger = log.With().Caller().Logger()
}

// CallerMarshalFuncWithShortFileName Details here: https://github.com/rs/zerolog#add-file-and-line-number-to-log
func CallerMarshalFuncWithShortFileName(pc uintptr, file string, line int) string {
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short
	return file + ":" + strconv.Itoa(line)
}

// ConfigureCommonFieldsInLogMessages configures one letter names for standard fields (timestamp/level/message)
// This is to make logged messages as short and compact as possible due to cost
// efficiency (when sending logs to services like AWS CloudWatch where we pay for volume of pushed data)
func ConfigureCommonFieldsInLogMessages() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.TimestampFieldName = "T"
	zerolog.LevelFieldName = "L"
	zerolog.MessageFieldName = "M"
}

// ConfigureDefaultLoggingSetup should be used in main file to configure common logging setup
func ConfigureDefaultLoggingSetup() {
	ConfigureShortFileNameInLogMessages()
	ConfigureCommonFieldsInLogMessages()
	// ConfigureLoggingMetrics()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.ErrorStackMarshaler = MarshalStack
	log.Logger = log.Output(os.Stdout)
}

func MarshalStack(err error) interface{} {
	ue := eris.Unpack(err)
	out := make([]map[string]string, 0, len(ue.ErrRoot.Stack))
	for _, frame := range ue.ErrRoot.Stack {
		// stop processing for stack not from hrnogomet
		// TODO: do we need this here?
		parsedPath := strings.Split(frame.File, "hrnogomet-api")
		if len(parsedPath) < 2 {
			break
		}
		file := fmt.Sprintf("%s:%d", parsedPath[len(parsedPath)-1], frame.Line)
		out = append(out, map[string]string{
			"source": file,
			"func":   frame.Name,
		})
	}
	return out
}
