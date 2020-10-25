package config

type Config struct {
	MemberUsername         string
	MemberPassword         string
	GymBaseURL             string
	GymClassesPath         string
	GymLoginPath           string
	GymFormSelector        string
	EmailOAuthClientID     string
	EmailOAuthClientSecret string
	EmailOAuthAccessToken  string
	EmailOAuthRefreshToken string
	EmailToUser            string
}

func (c *Config) GetClassesURL() string {
	return c.GymBaseURL + c.GymClassesPath
}

func (c *Config) GetLoginURL() string {
	return c.GymBaseURL + c.GymLoginPath
}
