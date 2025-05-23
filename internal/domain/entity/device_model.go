package entity

import (
	"time"

	"github.com/google/uuid"
)

type Device struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name       string    `json:"name" gorm:"not null;size:100"`
	Type       string    `json:"type" gorm:"not null;size:50"` // "raspberry_pi", "mini_pc"
	MacAddress string    `json:"mac_address" gorm:"uniqueIndex;not null;size:17"`
	IPAddress  string    `json:"ip_address" gorm:"not null;size:15"`
	Location   string    `json:"location" gorm:"not null;size:200"`
	IsOnline   bool      `json:"is_online" gorm:"default:false"`
	LastSeen   time.Time `json:"last_seen" gorm:"autoUpdateTime"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type DeviceRequest struct {
	Name       string `json:"name" validate:"required"`
	Type       string `json:"type" validate:"required,oneof=raspberry_pi mini_pc"`
	MacAddress string `json:"mac_address" validate:"required"`
	IPAddress  string `json:"ip_address" validate:"required,ip"`
	Location   string `json:"location" validate:"required"`
}

type DeviceResponse struct {
	*Device
	MonitoringData *MonitoringData `json:"monitoring_data,omitempty"`
}
