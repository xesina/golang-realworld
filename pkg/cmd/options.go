package cmd

// Options specifies cmd options
type Options struct {
	Name          string
	Description   string
	Version       string
	Usage         string
	ServerName    string
	ServerVersion string
	ServerID      string
	ServerAddress string
	DatabaseURI   string
	Env           string
	Debug         bool
	LogLevel      string
}

// Name Command line Name
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// Description Command line Description
func Description(d string) Option {
	return func(o *Options) {
		o.Description = d
	}
}

// Version Command line Version
func Version(v string) Option {
	return func(o *Options) {
		o.Version = v
	}
}

// Usage Command line Usage
func Usage(v string) Option {
	return func(o *Options) {
		o.Usage = v
	}
}

// DatabaseURI specifies database uri
func DatabaseURI(u string) Option {
	return func(o *Options) {
		o.DatabaseURI = u
	}
}

// LogLevel specifies log level
func LogLevel(l string) Option {
	return func(o *Options) {
		o.LogLevel = l
	}
}

// Env specifies runtime environment
func Env(e string) Option {
	return func(o *Options) {
		o.Env = e
	}
}

// Debug turns debug on
func Debug(d bool) Option {
	return func(o *Options) {
		o.Debug = d
	}
}
