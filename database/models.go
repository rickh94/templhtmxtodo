package database

import (
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/nrednav/cuid2"
	"gorm.io/gorm"
)

// Base contains common columns for all tables.
type BaseModel struct {
	ID        string `gorm:"type:string;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// BeforeCreate will set a CUID rather than numeric ID.
func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	b.ID = cuid2.Generate()
	return nil
}

type User struct {
	BaseModel
	Name        string
	Email       string       `gorm:"not null;unique;index"`
	Credentials []Credential `gorm:"constraint:OnDelete:CASCADE"`
	Todos       []Todo       `gorm:"constraint:OnDelete:CASCADE"`
}

type Credential struct {
	BaseModel
	CredentialID      []byte                            `json:"credentialId" gorm:"not null;unique;index"`
	PublicKey         []byte                            `json:"publicKey" gorm:"not null"`
	Transport         []protocol.AuthenticatorTransport `json:"transports" gorm:"not null;serializer:json" `
	AttenestationType string                            `json:"attestationType"`
	Flags             webauthn.CredentialFlags          `json:"flags" gorm:"serializer:json"`
	Authenticator     webauthn.Authenticator            `json:"authenticator" gorm:"serializer:json"`
	UserID            string                            `json:"userId"`
	User              User
}

type Todo struct {
	BaseModel
	Text      string `gorm:"not null"`
	Completed bool   `gorm:"default:false"`
	UserID    string `json:"userId"`
	User      User
}
