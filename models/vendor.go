package models

import (
	"fmt"
	"net/http"
)

type Vendor struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

type VendorList struct {
	Vendors []Vendor `json:"vendors"`
}

func (i *Vendor) Bind(r *http.Request) error {
	if i.Name == "" {
		return fmt.Errorf("name is a required field")
	}
	return nil
}

func (*VendorList) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (*Vendor) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
