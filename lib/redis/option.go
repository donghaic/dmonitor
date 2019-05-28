package redis

type PoolOption struct {
	Address     string
	Password    string
	DB          int
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
}
