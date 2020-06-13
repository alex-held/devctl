package shared

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func BootstrapLogger(lvl zerolog.Level) {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: zerolog.TimeFormatUnix,
	})
	zerolog.SetGlobalLevel(lvl)
}
