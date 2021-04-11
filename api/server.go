package main

import (
	"os"
	conf "pandita/conf"
	"pandita/util"

	"github.com/juju/errors"
)

func prepareServer(pandita *conf.ViperConfig, cancel <-chan os.Signal) bool {

	log, err := util.InitLog("server", pandita.GetString("loglevel"))
	if err != nil {
		log.Infow("InitLog", "err", errors.Details(err))
		return false
	}

	return true
}
