package config

type Opts struct {
	Verbose     bool `short:"v" long:"verbose" description:"Show verbose information"`
	Debug       bool `short:"d" long:"debug" description:"Show debug information"`
	ShowVersion bool `long:"version" description:"Show version information and exit"`
}
