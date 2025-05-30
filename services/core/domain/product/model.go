package product

import "time"

type Product struct {
	ProductID int64
	Name      string
	Price     float64
	// add more fields here
}

type ProductView struct {
	ID        string
	UserID    string
	ProductID int64
	ViewAt    time.Time
}
