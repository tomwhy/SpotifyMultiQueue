package secrets

type LiteralSecret struct {
	secret string
}

func NewLiteralSecret(secret string) LiteralSecret {
	return LiteralSecret{
		secret: secret,
	}
}

func (s LiteralSecret) Read() ([]byte, error) {
	return []byte(s.secret), nil
}