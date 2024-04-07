package other

type Config struct {
	Analyse Analyse `json:"analyse"`
}

type Analyse struct {
	ClientID     string `json:"client_id"     required:"true"`
	ClientSecret string `json:"client_secret" required:"true"`
	GetTokenUrl  string `json:"get_token_url" required:"true"`
}
