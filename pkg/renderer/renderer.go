package renderer

// Renderer - An instance of a Renderer much take a file and render it's output.
type Renderer interface {
	Render(filename string) (string, error)
}
