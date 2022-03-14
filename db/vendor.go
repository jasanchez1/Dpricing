package db

import (
	"database/sql"

	"github.com/jasanchez1/Dpricing/models"
)

func (db Database) GetAllVendors() (*models.VendorList, error) {
	list := &models.VendorList{}

	rows, err := db.Conn.Query("SELECT * FROM vendors ORDER BY ID DESC")
	if err != nil {
		return list, err
	}

	for rows.Next() {
		var vendor models.Vendor
		err := rows.Scan(&vendor.ID, &vendor.Name, &vendor.Description, &vendor.CreatedAt)
		if err != nil {
			return list, err
		}
		list.Vendors = append(list.Vendors, vendor)
	}
	return list, nil
}

func (db Database) AddVendor(vendor *models.Vendor) error {
	var id int
	var createdAt string
	query := `INSERT INTO vendors (name, description) VALUES ($1, $2) RETURNING id, created_at`
	err := db.Conn.QueryRow(query, vendor.Name, vendor.Description).Scan(&id, &createdAt)
	if err != nil {
		return err
	}

	vendor.ID = id
	vendor.CreatedAt = createdAt
	return nil
}

func (db Database) GetVendorById(vendorId int) (models.Vendor, error) {
	vendor := models.Vendor{}

	query := `SELECT * FROM vendors WHERE id = $1;`
	row := db.Conn.QueryRow(query, vendorId)
	switch err := row.Scan(&vendor.ID, &vendor.Name, &vendor.Description, &vendor.CreatedAt); err {
	case sql.ErrNoRows:
		return vendor, ErrNoMatch
	default:
		return vendor, err
	}
}

func (db Database) DeleteVendor(vendorId int) error {
	query := `DELETE FROM vendors WHERE id = $1;`
	_, err := db.Conn.Exec(query, vendorId)
	switch err {
	case sql.ErrNoRows:
		return ErrNoMatch
	default:
		return err
	}
}

func (db Database) UpdateVendor(vendorId int, vendorData models.Vendor) (models.Vendor, error) {
	vendor := models.Vendor{}
	query := `UPDATE vendors SET name=$1, description=$2 WHERE id=$3 RETURNING id, name, description, created_at;`
	err := db.Conn.QueryRow(query, vendorData.Name, vendorData.Description, vendorId).Scan(&vendor.ID, &vendor.Name, &vendor.Description, &vendor.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return vendor, ErrNoMatch
		}
		return vendor, err
	}

	return vendor, nil
}
