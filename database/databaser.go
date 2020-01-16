package database

type IDatabase interface {
	ConnectionArgs() string
	Driver() string
	Prefix() string
}
