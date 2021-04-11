package model

import (
	"pandita/util"
)

var mlog *util.MLogger

func init() {
	mlog, _ = util.InitLog("model", "console")
}
