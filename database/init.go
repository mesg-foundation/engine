package database

import "github.com/mesg-foundation/core/database/json"

var Db Database

func init() {
	Db = &json.Database{}
}
