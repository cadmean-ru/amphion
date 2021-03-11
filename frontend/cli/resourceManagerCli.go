package cli

type ResourceManagerCLI interface {
	ReadFile(path string) ([]byte, error)
}
