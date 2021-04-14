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

type CoreSignal struct {
	Prefix string `toml:"prefix"`
	Name   string `toml:"name"`
}

type Signal struct {
	OverrideCoreName   string `toml:"override_core_name"`
	OverrideRoomPrefix string `toml:"override_room_prefix"`
	OverridePanelName  string `toml:"override_panel_name"`
	OverrideDeviceName string `toml:"override_device_name"`
	RoomLevelSignal    int    `toml:"room_level_signal"`
	SystemType         string `toml:"system_type"`
	DeviceType         string `toml:"device_type"`
	PanelSignalModif   string `toml:"panel_signal_modif"`
	CoreSignalModif    string `toml:"core_signal_modif"`
	CoreSuffix         string `toml:"core_suffix"`
	PanelSuffix        string `toml:"panel_suffix"`
}

type Config struct {
	SheetConfig SheetConfig
	CoreSignal  CoreSignal
	Signals     []Signal
}
