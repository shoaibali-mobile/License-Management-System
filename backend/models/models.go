package models

import (
	"time"
)

// User represents the authentication user
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Role         string    `gorm:"not null;default:'customer'" json:"role"` // 'admin' or 'customer'
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Customer     *Customer `gorm:"foreignKey:UserID" json:"customer,omitempty"`
	APIKey       string    `gorm:"index" json:"-"`
}

// Customer represents customer profile information
type Customer struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	UserID       uint       `gorm:"uniqueIndex;not null" json:"user_id"`
	Name         string     `gorm:"not null" json:"name"`
	Phone        string     `json:"phone"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	User         User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Subscriptions []Subscription `gorm:"foreignKey:CustomerID" json:"subscriptions,omitempty"`
}

// SubscriptionPack represents available subscription plans
type SubscriptionPack struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	Name          string     `gorm:"not null" json:"name"`
	Description   string     `json:"description"`
	SKU           string     `gorm:"uniqueIndex;not null" json:"sku"`
	Price         float64    `gorm:"not null" json:"price"`
	ValidityMonths int       `gorm:"not null;check:validity_months >= 1 AND validity_months <= 12" json:"validity_months"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	Subscriptions []Subscription `gorm:"foreignKey:PackID" json:"subscriptions,omitempty"`
}

// Subscription represents customer subscription assignments
type Subscription struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	CustomerID    uint       `gorm:"not null;index" json:"customer_id"`
	PackID        uint       `gorm:"not null;index" json:"pack_id"`
	Status        string     `gorm:"not null;default:'requested'" json:"status"` // requested, approved, active, inactive, expired
	RequestedAt   time.Time  `json:"requested_at"`
	ApprovedAt    *time.Time `json:"approved_at,omitempty"`
	AssignedAt    *time.Time `json:"assigned_at,omitempty"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
	DeactivatedAt *time.Time `json:"deactivated_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Customer      Customer   `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Pack          SubscriptionPack `gorm:"foreignKey:PackID" json:"pack,omitempty"`
}







