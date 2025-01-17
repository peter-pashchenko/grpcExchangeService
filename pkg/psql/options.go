package psql

type psqlOptions struct {
	host     string
	port     string
	database string

	user string
	pass string

	migrations string
}

type OptionFunc func(*psqlOptions)

func WithHost(host string) OptionFunc {
	return func(o *psqlOptions) {
		o.host = host
	}
}
func WithPort(port string) OptionFunc {
	return func(o *psqlOptions) {
		o.port = port
	}
}
func WithDatabase(database string) OptionFunc {
	return func(o *psqlOptions) {
		o.database = database
	}
}
func WithUser(user string) OptionFunc {
	return func(o *psqlOptions) {
		o.user = user
	}
}
func WithPass(pass string) OptionFunc {
	return func(o *psqlOptions) {
		o.pass = pass
	}
}
func WithMigrations(migrations string) OptionFunc {
	return func(o *psqlOptions) {
		o.migrations = migrations
	}
}
