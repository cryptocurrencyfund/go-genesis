// Copyright 2016 The go-daylight Authors
// This file is part of the go-daylight library.
//
// The go-daylight library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-daylight library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-daylight library. If not, see <http://www.gnu.org/licenses/>.

package schema

import (
	"time"

	"github.com/AplaProject/go-apla/packages/consts"
	"github.com/AplaProject/go-apla/packages/model"
	"github.com/AplaProject/go-apla/packages/utils"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("daemons")
)

func Migration() {
	oldDbVersion, err := model.Single(`SELECT version FROM migration_history ORDER BY id DESC LIMIT 1`).String()
	if err != nil {
		log.Error("%v", utils.ErrInfo(err))
	}
	if len(*utils.OldVersion) == 0 && consts.VERSION != oldDbVersion {
		*utils.OldVersion = oldDbVersion
	}

	log.Debug("*utils.OldVersion %v", *utils.OldVersion)
	if len(*utils.OldVersion) > 0 {
		err = model.InsertIntoMigration(consts.VERSION, time.Now().Unix())
		if err != nil {
			log.Error("%v", utils.ErrInfo(err))
		}
	}
}
