package config

type Config struct {
	Theme     string
	SwwwCache string

	NixStdout  bool
	NixOut     string
	JSONStdout bool
	JSONOut    string
	SCSSStdout bool
	SCSSOut    string
	CSSStdout  bool
	CSSOut     string

	Activate   bool
}
