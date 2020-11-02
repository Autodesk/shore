package renderer

// Renderer - An instance of a Renderer much take a file and render it's output.
type Renderer interface {
	Render(projectPath string) (string, error)
}
