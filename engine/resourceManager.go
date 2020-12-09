package engine

type ResourceManager struct {
	//availableResources map[string]ResourceDefinition
	reader             resourceReader
}

//func (m *ResourceManager) registerResource(definition ResourceDefinition) {
//	if !IsValidResourcePath(definition.path) {
//		return
//	}
//
//	m.availableResources[definition.path] = definition
//}

func (m *ResourceManager) GetFile(path string) ReadableResource {
	//if res, ok := m.availableResources[path]; ok {
		return &File{
			name: GetResourceName(path),
			path: path,
		}
	//}
	//return nil
}

func newResourceManager() *ResourceManager {
	return &ResourceManager{
		//availableResources: make(map[string]ResourceDefinition),
		reader:             newResourceReader(),
	}
}

type resourceReader interface {
	readResource(path string) ([]byte, error)
	readResourceAsync(path string, callback ReadResourceCallback)
}
