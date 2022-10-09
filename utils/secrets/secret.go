package secrets

type Secret interface {
	Read() ([]byte, error)
}