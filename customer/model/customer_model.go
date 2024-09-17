package customermodel

import (
	merchantmodel "sample/merchant/model"
	"time"
)

// Person model
type Customer struct {
	ID                           uint                     `gorm:"primaryKey;autoIncrement" json:"id"`
	Title                        string                   `gorm:"size:10" json:"title"`
	FullName                     string                   `gorm:"size:100;not null" json:"full_name"`
	LastName                     string                   `gorm:"size:100;not null" json:"last_name"`
	OwnerGender                  string                   `gorm:"size:10" json:"owner_gender"`
	DateOfBirth                  time.Time                `gorm:"not null" json:"date_of_birth"`
	PlaceOfBirth                 string                   `gorm:"size:100" json:"place_of_birth"`
	Job                          string                   `gorm:"size:50" json:"job"`
	TaxpayerIdentificationNumber string                   `gorm:"size:20;unique" json:"taxpayer_identification_number"`
	Addresses                    []Address                `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE;" json:"address"`
	Identifications              []Identification         `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE;" json:"identification"`
	Contacts                     []Contact                `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE;" json:"contact"`
	Merchant                     []merchantmodel.Merchant `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE;" json:"merchant"`
}


// Address model
type Address struct {
	ID           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	CustomerID     int    `gorm:"index;not null" json:"customer_id"`
	Address      string `gorm:"size:255" json:"address"`
	Region       string `gorm:"size:50" json:"region"`
	Province     string `gorm:"size:50" json:"province"`
	Municipality string `gorm:"size:50" json:"municipality"`
	Barangays    string `gorm:"size:50" json:"barangays"`
	PostalCode   string `gorm:"size:10" json:"postal_code"`
}


// Identification model
type Identification struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CustomerID     int       `gorm:"index;not null" json:"customer_id"`
	IDType       string    `gorm:"size:50" json:"id_type"`
	IDNumber     string    `gorm:"size:50;unique" json:"id_number"`
	IDExpiryDate time.Time `json:"id_expiry_date"`
}


// Contact model
type Contact struct {
	ID                    uint   `gorm:"primaryKey;autoIncrement" json:"contact_id"`
	CustomerID              int    `gorm:"index;not null" json:"customer_id"`
	OwnerPhoneNumber      string `gorm:"size:15" json:"owner_phone_number"`
	OwnerOtherPhoneNumber string `gorm:"size:15" json:"owner_other_phone_number"`
	Email                 string `gorm:"size:100;unique" json:"email"`
}

