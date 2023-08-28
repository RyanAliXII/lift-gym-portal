package model

type EmailVerification struct {
	Id        int    `json:"id"`
	PublicId  string `json:"publicId" db:"public_id"`
	ClientId  int    `json:"clientId" db:"client_id"`
	ExpiresAt string `json:"expiresAt" db:"expires_at"`
	CreatedAt string `json:"createdAt" db:"created_at"`
}