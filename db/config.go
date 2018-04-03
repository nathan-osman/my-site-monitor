package db

// Config stores the database connection information.
type Config struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}
