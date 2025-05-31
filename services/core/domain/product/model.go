package product

import "time"

type Product struct {
	ProductID  int64
	Name       string
	Price      float64
	CategoryID string
	// add more fields here
}

type UserViewHistory struct {
	ID        string
	UserID    string
	ProductID int64
	ViewAt    time.Time
}

type CategoryViewHistory struct {
	ID         string
	CategoryID string
	TotalView  int64
	LastViewAt time.Time
}

type SummaryProductView struct {
	ProductID int64
	ViewCount *int64
	ViewAt    *time.Time
}
