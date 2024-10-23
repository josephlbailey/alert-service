package config

type Config struct {
	Environment string `mapstructure:"environment"`
	Port        string `mapstructure:"port"`

	DB    DBConfig    `mapstructure:"db"`
	Users []BasicUser `mapstructure:"users"`
}

type DBConfig struct {
	Username          string `mapstructure:"username"`
	Password          string `mapstructure:"password"`
	MigrationUsername string `mapstructure:"migration_username"`
	MigrationPassword string `mapstructure:"migration_password"`
	Host              string `mapstructure:"host"`
	Port              string `mapstructure:"port"`
	Database          string `mapstructure:"database"`
	SslMode           string `mapstructure:"ssl_mode"`

	Url string
}

type BasicUser struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}
