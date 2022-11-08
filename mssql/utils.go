package mssql

import (
	"fmt"

	"github.com/betr-io/terraform-provider-mssql/mssql/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rs/zerolog"
)

func getLoginID(host string, port string, data *schema.ResourceData) string {
	loginName := data.Get(loginNameProp).(string)
	return fmt.Sprintf("sqlserver://%s:%s/%s", host, port, loginName)
}

func getUserID(host string, port string, data *schema.ResourceData) string {
	database := data.Get(databaseProp).(string)
	username := data.Get(usernameProp).(string)
	return fmt.Sprintf("sqlserver://%s:%s/%s/%s", host, port, database, username)
}

func loggerFromMeta(meta interface{}, resource, function string) zerolog.Logger {
	return meta.(model.Provider).ResourceLogger(resource, function)
}
