package models

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type UserTokens struct {
	Tokens Tokens `json:"tokens"`
}
