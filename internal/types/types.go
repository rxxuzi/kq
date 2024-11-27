// Package types internal/types/types.go
package types

type Config struct {
	Name    string   `json:"name"`
	Version string   `json:"version"`
	Source  []string `json:"source"`
	Allow   []string `json:"allow"`
	Entry   string   `json:"entry"`
	Options Options  `json:"options"`
}

type Options struct {
	Window   WindowOptions   `json:"window"`
	Security SecurityOptions `json:"security"`
	Debug    DebugOptions    `json:"debug"`
	Env      EnvOptions      `json:"env"`
	UI       UIOptions       `json:"ui"`
}

type WindowOptions struct {
	Height    int  `json:"height"`
	Width     int  `json:"width"`
	Frameless bool `json:"frameless"`
	Resizable bool `json:"resizable"`
}

type SecurityOptions struct {
	NoScript    bool     `json:"noscript"`
	LocalOnly   bool     `json:"localOnly"`
	AllowOrigin []string `json:"allowOrigin"`
}

type DebugOptions struct {
	Devtools bool `json:"devtools"`
	Console  bool `json:"console"`
}

type EnvOptions struct {
	SingleInstance bool `json:"singleInstance"`
}

type UIOptions struct {
	Theme string  `json:"theme"`
	Title string  `json:"title"`
	Icon  *string `json:"icon"`
}
