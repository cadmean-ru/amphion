package frontend

type App struct {
	Name          string `yaml:"name"`
	Author        string `yaml:"author"`
	CompanyDomain string `yaml:"companyDomain"`
	PublicUrl     string `yaml:"publicUrl"`
	Frontend      string `yaml:"frontend"`
	Debug         string `yaml:"debug"`
	MainScene     string `yaml:"mainScene"`
}
