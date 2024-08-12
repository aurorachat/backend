package config

import "github.com/spf13/viper"

func GetListenHost() string {
	return viper.GetString("host")
}

func GetPostgresDSN() string {
	return viper.GetString("postgres_dsn")
}

func GetAuthSecretKey() string {
	return viper.GetString("auth_secret_key")
}

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigFile("config.yaml")
	viper.SetDefault("host", "0.0.0.0:8080")
	viper.SetDefault("postgres_dsn", "host=localhost user=aurorachat password=123456 dbname=aurorachat port=5432 sslmode=disable TimeZone=Asia/Tbilisi")
	viper.SetDefault("auth_secret_key", "CHANGEITORELSE")
	err := viper.ReadInConfig()
	if err != nil {
		err = viper.WriteConfig()
		if err != nil {
			panic(err)
		}
	}
}
