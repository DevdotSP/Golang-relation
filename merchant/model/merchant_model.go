package merchantmodel

import (
	"time"
)

// Identification model
type Merchant struct {
	ID              uint              `gorm:"primaryKey;autoIncrement" json:"id"`
	CustomerID     int   `gorm:"index;not null" json:"customer_id"`
	Name        string    `gorm:"size:50" json:"name"`
	Product        []Product        `gorm:"foreignKey:MerchantID;constraint:OnDelete:CASCADE;" json:"product"`
	AddressMerchant []AddressMerchant `gorm:"foreignKey:MerchantID;constraint:OnDelete:CASCADE;" json:"address_merchant"`

	ContactMerchant []ContactMerchant `gorm:"foreignKey:MerchantID;constraint:OnDelete:CASCADE;" json:"contact_merchant"`
}

// Person model
type Product struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	MerchantID  int       `gorm:"index;not null" json:"merchant_id"`
	Name        string    `gorm:"size:20" json:"name"`
	Quantity    int       `gorm:"size:100;not null" json:"quantity"`
	DeliverDate time.Time `gorm:"not null" json:"date_of_delivery"`
}

// Address model
type AddressMerchant struct {
	ID           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	MerchantID   int    `gorm:"index;not null" json:"merchant_id"`
	Address      string `gorm:"size:255" json:"address"`
	Region       string `gorm:"size:50" json:"region"`
	Province     string `gorm:"size:50" json:"province"`
	Municipality string `gorm:"size:50" json:"municipality"`
	Barangays    string `gorm:"size:50" json:"barangays"`
	PostalCode   string `gorm:"size:10" json:"postal_code"`
}

// Contact model
type ContactMerchant struct {
	ID                  uint   `gorm:"primaryKey;autoIncrement" json:"merchant_contact_id"`
	MerchantID          int    `gorm:"index;not null" json:"merchant_id"`
	MerchantPhoneNumber string `gorm:"size:20" json:"merchant_phone_number"`
	MerchantEmail       string `gorm:"size:100;unique" json:"merchant_email"`
}
