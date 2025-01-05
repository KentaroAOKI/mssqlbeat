// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Input struct {
	Field                            string `config:"field"`
	Enabled                          bool   `config:"enabled"`
	MssqlserverHost                  string `config:"mssqlserver_host"`
	MssqlserverPort                  int    `config:"mssqlserver_port"`
	MssqlserverUserId                string `config:"mssqlserver_userid"`
	MssqlserverPassword              string `config:"mssqlserver_password"`
	MssqlserverDatabase              string `config:"mssqlserver_database"`
	MssqlserverTlsmin                string `config:"mssqlserver_tlsmin"`
	SqlQuery                         string `config:"sql_query"`
	SqlTimeColumn                    string `config:"sql_time_column"`
	SqlTimeInitializeWithCurrentTime bool   `config:"sql_time_initialize_with_current_time"`
	FieldPrefix                      string `config:"field_prefix"`
}

type Config struct {
	Period  time.Duration `config:"period"`
	Inputs  []Input       `config:"inputs"`
	Threads int           `config:"threads"`
}

var DefaultConfig = Config{
	Period: 1 * time.Second,
}
