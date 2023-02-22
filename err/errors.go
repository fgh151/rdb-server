package err

import log "github.com/sirupsen/logrus"

func PanicErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func DebugErr(err error) {
	if err != nil {
		log.Debug("Err " + err.Error())
	}
}

func WarnErr(err error) {
	if err != nil {
		log.Warn("Err " + err.Error())
	}
}
