package pkg

import "fmt"

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Logger   LoggerConfig   `mapstructure:"logger"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Driver    string `mapstructure:"driver"`
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	Database  string `mapstructure:"database"`
	Charset   string `mapstructure:"charset"`
	ParseTime bool   `mapstructure:"parse_time"`
	Loc       string `mapstructure:"loc"`
}

func (c *DatabaseConfig) GetDatabaseUrl() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.Charset,
		c.ParseTime,
		c.Loc,
	)
}

type LoggerConfig struct {
	LogDir      string `mapstructure:"log_dir"`
	MinLevel    string `mapstructure:"min_level"`
	MaxSizeMB   int    `mapstructure:"max_size_mb"`
	MaxBackups  int    `mapstructure:"max_backups"`
	MaxAgeDays  int    `mapstructure:"max_age_days"`
	Compress    bool   `mapstructure:"compress"`
}