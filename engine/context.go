package engine

type GlobalContext struct {
	deviceInfo DeviceInfo
	screenInfo ScreenInfo
	domain     string
	host       string
	port       string
}

func (c *GlobalContext) GetDeviceInfo() DeviceInfo {
	return c.deviceInfo
}

func (c *GlobalContext) GetScreenInfo() ScreenInfo {
	return c.screenInfo
}

func (c *GlobalContext) GetDomain() string {
	return c.domain
}

func (c *GlobalContext) GetHost() string {
	return c.host
}

func (c *GlobalContext) GetPort() string {
	return c.port
}

func (c *GlobalContext) FromMap(m map[string]interface{}) {

}

func newGlobalContext() *GlobalContext {
	return &GlobalContext{}
}
