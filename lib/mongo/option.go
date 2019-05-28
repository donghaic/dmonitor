package mongo

type DBOption struct {
	Addr    string
	User    string
	Password  string
	Timeout int
	DBName  string
}
