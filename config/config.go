package config

type SheetConfig struct {
	SheetName        string `toml:"sheet_name"`
	ColomnRoom       string `toml:"colomn_room"`
	ColomnRoomNumber string `toml:"colomn_room_number"`
	ColomnLight      string `toml:"colomn_light"`
	ColomnLightType  string `toml:"colomn_light_type"`
	ColomnShade      string `toml:"colomn_shade"`
	ColomnShadeType  string `toml:"colomn_shade_type"`
	StartRow         int    `toml:"start_row"`
}

type Prefix struct {
	RoomPrefix string `toml:"room_prefix"`
}

type LightSuffix struct {
	SignalName string `toml:"signal_name"`
	On         string `toml:"on"`
	OnFb       string `toml:"on_fb"`
	Off        string `toml:"off"`
	OffFb      string `toml:"off_fb"`
	Raise      string `toml:"raise"`
	Dim        string `toml:"dim"`
	Level      string `toml:"level"`
	LevelFb    string `toml:"level_fb"`
}

type ShadeSuffix struct {
	SignalName string `toml:"signal_name"`
	Open       string `toml:"open"`
	OpenFb     string `toml:"open_fb"`
	Close      string `toml:"close"`
	CloseFb    string `toml:"close_fb"`
	Stop       string `toml:"stop"`
	StopFb     string `toml:"stop_fb"`
	Level      string `toml:"level"`
	LevelFb    string `toml:"level_fb"`
}

type Config struct {
	SheetConfig   SheetConfig
	Prefix        Prefix
	KeypanelLight LightSuffix
	CoreLight     LightSuffix
	KeypanelShade ShadeSuffix
	CoreShade     ShadeSuffix
}
