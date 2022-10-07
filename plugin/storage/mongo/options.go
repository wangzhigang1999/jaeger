package mongo

import (
	"flag"
	"github.com/spf13/viper"
)

const (
	mongoHost       = "mongo_url"
	mongoPort       = "mongo_port"
	mongoDatabase   = "mongo_database"
	mongoCollection = "mongo_collection"
	mongoUserName   = "mongo_user"
	mongoPassWord   = "mongo_pass"
	output          = "output"
	namespacePrefix = ""
)

type Configuration struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	UserName   string `json:"userName"`
	PassWord   string `json:"passWord"`
	Collection string `json:"collection"`
	Database   string `json:"database"`
	Output     bool   `json:"output"`
}

type Options struct {
	Configuration Configuration
}

// InitFromViper initializes the options struct with values from Viper
func (opt *Options) InitFromViper(v *viper.Viper) {
	opt.Configuration.Host = v.GetString(mongoHost)
	opt.Configuration.Port = v.GetInt(mongoPort)
	opt.Configuration.Database = v.GetString(mongoDatabase)
	opt.Configuration.Collection = v.GetString(mongoCollection)
	opt.Configuration.PassWord = v.GetString(mongoPassWord)
	opt.Configuration.UserName = v.GetString(mongoUserName)
	opt.Configuration.Output = v.GetBool(output)

}

func (opt Options) AddFlags(flagSet *flag.FlagSet) {
	flagSet.String(
		namespacePrefix+mongoHost,
		"localhost",
		"mongodb host",
	)
	flagSet.Int(
		namespacePrefix+mongoPort,
		27017,
		"mongodb host",
	)
	flagSet.String(
		namespacePrefix+mongoCollection,
		"span",
		"which collection you want to store data in",
	)
	flagSet.String(
		namespacePrefix+mongoDatabase,
		"sock-shop",
		"which database you want to store data in",
	)
	flagSet.String(
		namespacePrefix+mongoUserName,
		"root",
		"mongodb username",
	)
	flagSet.String(
		namespacePrefix+mongoPassWord,
		"root",
		"mongodb password",
	)
	flagSet.Bool(
		namespacePrefix+output,
		false,
		"print span to screen in json format",
	)

}
