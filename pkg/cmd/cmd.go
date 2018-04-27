package cmd

import (
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/urfave/cli"
)

// Cmd specifies cli apps
type Cmd interface {
	// The cli app within this cmd
	App() *cli.App
	// Adds options, parses flags and initialise
	// exits on error
	Init(opts ...Option) error
	// Options set within this command
	Options() Options
}

type cmd struct {
	opts Options
	app  *cli.App
}

// Option specifies cli options
type Option func(o *Options)

var (
	// DefaultCmd default cmd
	DefaultCmd = newCmd()

	// DefaultFlags default cmd flags
	DefaultFlags = []cli.Flag{
		cli.StringFlag{
			Name:   "server_name",
			EnvVar: "SERVER_NAME",
			Usage:  "Name of the server",
		},
		cli.StringFlag{
			Name:   "server_version",
			EnvVar: "SERVER_VERSION",
			Usage:  "Version of the server. eg. 1.1.0",
		},
		cli.StringFlag{
			Name:   "server_id",
			EnvVar: "SERVER_ID",
			Usage:  "Id of the server. Auto-generated if not specified",
		},
		cli.StringFlag{
			Name:   "server_address",
			EnvVar: "SERVER_ADDRESS",
			Usage:  "Bind address for the server. default 127.0.0.1:8000",
		},
		cli.StringFlag{
			Name:   "database_uri",
			EnvVar: "DATABASE_URI",
			Usage:  "database connection URI",
		},
		cli.StringFlag{
			Name:   "log-level",
			EnvVar: "LOG_LEVEL",
			Usage:  "Log level",
		},
		cli.BoolFlag{
			Name:   "debug",
			EnvVar: "DEBUG",
			Usage:  "Turn on debug",
		},
		cli.StringFlag{
			Name:   "env",
			EnvVar: "ENV",
			Usage:  "Specifies runtime environment",
		},
	}
)

func init() {
	rand.Seed(time.Now().Unix())
	help := cli.HelpPrinter
	cli.HelpPrinter = func(writer io.Writer, templ string, data interface{}) {
		help(writer, templ, data)
		os.Exit(0)
	}
}

func newCmd(opts ...Option) Cmd {
	options := Options{}

	for _, o := range opts {
		o(&options)
	}

	if len(options.Description) == 0 {
		options.Description = "a cli app"
	}

	cmd := new(cmd)
	cmd.opts = options
	cmd.app = cli.NewApp()
	cmd.app.Name = cmd.opts.Name
	cmd.app.Version = cmd.opts.Version
	cmd.app.Usage = cmd.opts.Description
	cmd.app.Before = cmd.Before
	cmd.app.Flags = DefaultFlags
	cmd.app.Action = func(ctx *cli.Context) error {
		return nil
	}
	if len(options.Version) == 0 {
		cmd.app.HideVersion = true
	}

	return cmd
}

func (c *cmd) App() *cli.App {
	return c.app
}

func (c *cmd) Options() Options {
	return c.opts
}

func (c *cmd) Before(ctx *cli.Context) error {

	if len(ctx.String("server_name")) > 0 {
		c.opts.ServerName = ctx.String("server_name")
	}

	if len(ctx.String("server_version")) > 0 {
		c.opts.ServerVersion = ctx.String("server_version")
	}

	if len(ctx.String("server_id")) > 0 {
		c.opts.ServerID = ctx.String("server_id")
	}

	if len(ctx.String("server_address")) > 0 {
		c.opts.ServerAddress = ctx.String("server_address")
	}

	if len(ctx.String("database_uri")) > 0 {
		c.opts.DatabaseURI = ctx.String("database_uri")
	}

	if len(ctx.String("log_level")) > 0 {
		c.opts.LogLevel = ctx.String("log_level")
	}

	if ctx.Bool("debug") {
		c.opts.Debug = ctx.Bool("debug")
	}

	if len(ctx.String("env")) > 0 {
		c.opts.Env = ctx.String("env")
	}

	return nil
}

func (c *cmd) Init(opts ...Option) error {
	for _, o := range opts {
		o(&c.opts)
	}
	c.app.Name = c.opts.Name
	c.app.Version = c.opts.Version
	c.app.HideVersion = len(c.opts.Version) == 0
	c.app.Usage = c.opts.Description
	return c.app.Run(os.Args)
}

// DefaultOptions default cmd options
func DefaultOptions() Options {
	return DefaultCmd.Options()
}

// App cli app
func App() *cli.App {
	return DefaultCmd.App()
}

// Init init cmd
func Init(opts ...Option) error {
	return DefaultCmd.Init(opts...)
}

// NewCmd initialize new cmd
func NewCmd(opts ...Option) Cmd {
	return newCmd(opts...)
}
