package clientData

type PC struct {
	UserAgent      string   `json:"userAgent"`
	Platform       string   `json:"platform"`
	Language       string   `json:"language"`
	Languages      []string `json:"languages"`
	CookiesEnabled bool     `json:"cookiesEnabled"`
	ConnectionType string   `json:"connectionType"`
	IsOnline       bool     `json:"isOnline"`
	Cores          int      `json:"cores"`
	Memory         int      `json:"memory"`
	ScreenWidth    int      `json:"screenWidth"`
	ScreenHeight   int      `json:"screenHeight"`
	ColorDepth     int      `json:"colorDepth"`
	PixelRatio     float64  `json:"pixelRatio"`
	ViewportWidth  int      `json:"viewportWidth"`
	ViewportHeight int      `json:"viewportHeight"`
	TimeZone       string   `json:"timeZone"`
	CurrentTime    string   `json:"currentTime"`
}

type IOS struct {
	UserAgent              string   `json:"userAgent"`
	Platform               string   `json:"platform"`
	Language               string   `json:"language"`
	Languages              []string `json:"languages"`
	CookiesEnabled         bool     `json:"cookiesEnabled"`
	ConnectionType         string   `json:"connectionType"`
	IsOnline               bool     `json:"isOnline"`
	Cores                  int      `json:"cores"`
	Memory                 int      `json:"memory"`
	ScreenWidth            int      `json:"screenWidth"`
	ScreenHeight           int      `json:"screenHeight"`
	ColorDepth             int      `json:"colorDepth"`
	PixelRatio             float64  `json:"pixelRatio"`
	ViewportWidth          int      `json:"viewportWidth"`
	ViewportHeight         int      `json:"viewportHeight"`
	TimeZone               string   `json:"timeZone"`
	CurrentTime            string   `json:"currentTime"`
	BatteryLevel           float64  `json:"batteryLevel"`
	BatteryCharging        bool     `json:"batteryCharging"`
	BatteryChargingTime    float64  `json:"batteryChargingTime"`
	BatteryDischargingTime float64  `json:"batteryDischargingTime"`
}
