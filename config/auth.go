package config

// AuthConfiguration holds required values to use with authorization service (while setting JWT middleware for example).
type AuthConfiguration struct {
	URI string
}

// GetAuthServiceURL provides URL used to call Authorization service.
func (c *AuthConfiguration) GetAuthServiceURL() string {
	return c.URI
}
