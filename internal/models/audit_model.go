package models

import (
	"time"
)

// AuditModel contains standard audit fields.
// Embed this into models that can be updated and soft-deleted (e.g., Users, Clubs).
type AuditModel struct {
	ID        int64      `bun:"id,pk,autoincrement" json:"id"`
	CreatedAt time.Time  `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time  `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt *time.Time `bun:"deleted_at,soft_delete,nullzero" json:"-"`
	CreatedBy *int64     `bun:"created_by" json:"created_by,omitempty"`
	UpdatedBy *int64     `bun:"updated_by" json:"updated_by,omitempty"`
}
