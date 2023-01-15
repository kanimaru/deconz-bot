package view

type (
	GroupsData struct {
		Groups map[string]string
	}

	LightsData struct {
		GroupName            string
		GroupId              string
		Lights               map[string]string
		On                   bool
		ColorAvailable       bool
		BrightnessAvailable  bool
		TemperatureAvailable bool
	}

	LightData struct {
		GroupName            string
		Id                   string
		Name                 string
		On                   bool
		ColorAvailable       bool
		BrightnessAvailable  bool
		TemperatureAvailable bool
	}
)
