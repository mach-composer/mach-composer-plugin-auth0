package internal

type Auth0Config struct {
	Domain        string `mapstructure:"domain"`
	ClientID      string `mapstructure:"client_id"`
	ClientSecret  string `mapstructure:"client_secret"`
}

func (c *Auth0Config) extendConfig(o *Auth0Config) *Auth0Config {
	cfg := &Auth0Config{
		Domain:        o.Domain,
		ClientID:      o.ClientID,
		ClientSecret:  o.ClientSecret,
	}
	if c.Domain != "" {
		cfg.Domain = c.Domain
	}
	if c.ClientID != "" {
		cfg.ClientID = c.ClientID
	}
	if c.ClientSecret != "" {
		cfg.ClientSecret = c.ClientSecret
	}
	return cfg
}
