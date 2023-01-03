package servicedefs

import (
	"github.com/sarulabs/dingo/v4"
	database "github.com/yapi-teklif/internal/pkg/database/connection"
)

var ExternalConnectionDefs = []dingo.Def{
	{
		Name: "external-connections",
		Build: func() (database.IConnection, error) {
			return database.Connect(), nil
		},
		Params: dingo.Params{},
	},
	{
		Name: "cache-connections",
		Build: func() (database.ICacheDB, error) {
			return database.NewCacheClient(), nil
		},
		Params: dingo.Params{},
	},
}
