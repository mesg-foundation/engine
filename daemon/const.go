package daemon

const (
	name = "core"
)

// Namespace returns the namespace of the MESG Core.
func Namespace() []string {
	return []string{name}
}
