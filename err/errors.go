package err

import log "github.com/sirupsen/logrus"

func PanicErr(err error) {
	if err != nil {
		log.Debug("Err " + err.Error())
		panic(err)
	}
}

func DebugErr(err error) {
	if err != nil {
		log.Debug("Err " + err.Error())
	}
}
