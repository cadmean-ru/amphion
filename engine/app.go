package engine

import (
	"gopkg.in/yaml.v2"
)

type App struct {
	Id        string   `yaml:"id"`
	Name      string   `yaml:"name"`
	appType   byte     `yaml:"type"`
	Scenes    []string `yaml:"scenes"`
	Resources []string `yaml:"resources"`
}

func NewApp(id, name string, aType byte) *App {
	return &App{
		Id:        id,
		Name:      name,
		appType:   aType,
		Scenes:    make([]string, 0, 1),
		Resources: make([]string, 0, 1),
	}
}

func (p *App) AddScene(scene string) {
	p.Scenes = append(p.Scenes, scene)
}

func (p *App) AddResource(path string) {
	p.Resources = append(p.Resources, path)
}

func DecodeApp(data []byte) (*App, error) {
	proj := App{}
	if err := yaml.Unmarshal(data, &proj); err != nil {
		return nil, err
	}
	return &proj, nil
}

func EncodeApp(app *App) ([]byte, error) {
	var data []byte
	var err error
	if data, err = yaml.Marshal(*app); err != nil {
		return nil, err
	}
	return data, nil
}
