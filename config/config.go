package config

// T is the basic config struct
type T struct {
	// Root for ftp to serve files from
	RootDir string
	// Path to sqlite file used for storing all mappings
	SqlitePath string
	// Address to listen on
	TCPAddr string
}

// Default returns a config (T) struct populated with default values
func Default() *T {
	config := &T{}

	config.RootDir = "/tmp/user/ezftp/files"
	config.SqlitePath = "/user/local/var/ez-ftp/sqlite.db"
	config.TCPAddr = "9999"

	return config
}
