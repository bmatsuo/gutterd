package log

type Config struct {
	Path    string   `json:"path"`    // Log output path (&2/&1 for stderr/stdout).
	Accepts []string `json:"accepts"` // Names logs accepted ("gutterd", "http", ...).
}
