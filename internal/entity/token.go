package entity

// JWT represents JSON Web Tokens used for authentication.
type JWT struct {
	AccessToken        string `json:"access_token"`  // Access token for authentication.
	RefreshToken       string `json:"refresh_token"` // Refresh token for obtaining a new access token.
	AccessTokenMaxAge  int    `json:"-"`             // Maximum age of the access token in seconds (not serialized).
	RefreshTokenMaxAge int    `json:"-"`             // Maximum age of the refresh token in seconds (not serialized).
	Domain             string `json:"-"`             // Domain to which the tokens are issued (not serialized).
}
