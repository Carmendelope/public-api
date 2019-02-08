package cli

import (
	"github.com/rs/zerolog/log"
	"strings"
)

func GetLabels(rawLabels string) map[string]string{
	labels := make(map[string]string, 0)

	split := strings.Split(rawLabels, ";")
	for _, l := range split{
		ls := strings.Split(l, ":")
		if len(ls) != 2{
			log.Fatal().Str("label", l).Msg("malformed label, expecting key:value")
		}
		labels[ls[0]] = ls[1]
	}
	return labels
}
