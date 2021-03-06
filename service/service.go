// Package service ...
//
// This layer will act as the business process handler.
// Any process will handled here. This layer will decide, which repository layer will use.
// And have responsibility to provide data to serve into delivery.
// Process the data doing calculation or anything will done here.
//
// Service layer will accept any input from Delivery layer,
// that already sanitized, then process the input could be storing into DB ,
// or Fetching from DB ,etc.
//
// This Service layer will depends to Repository Layer
package service

import (
	"context"
	"pandita/model"
	"pandita/util"
)

var mlog *util.MLogger

func init() {
	mlog, _ = util.InitLog("service", "console")
}

type UserService interface {
	GetUserByID(ctx context.Context, uid uint64) (user *model.User, err error)
}
