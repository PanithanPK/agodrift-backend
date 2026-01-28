package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"agodrift/internal/model"
)

var ErrNotEnoughRooms = errors.New("not enough rooms available")

type BookingRepository interface {
	Create(userID int, hotelID int, checkIn time.Time, checkOut time.Time, adults int, children int, rooms int) (model.Booking, error)
	ListByUserID(userID int) ([]model.Booking, error)
}

type mysqlBookingRepo struct {
	db *sql.DB
}

func NewMySQLBookingRepo(db *sql.DB) *mysqlBookingRepo {
	return &mysqlBookingRepo{db: db}
}

func (r *mysqlBookingRepo) Create(userID int, hotelID int, checkIn time.Time, checkOut time.Time, adults int, children int, rooms int) (model.Booking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return model.Booking{}, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// Lock hotel row to prevent overselling rooms
	var priceCents int
	var roomsAvailable int
	err = tx.QueryRowContext(ctx, "SELECT price_cents, rooms_available FROM hotels WHERE id = ? FOR UPDATE", hotelID).Scan(&priceCents, &roomsAvailable)
	if err != nil {
		return model.Booking{}, err
	}
	if roomsAvailable < rooms {
		return model.Booking{}, ErrNotEnoughRooms
	}

	// nights = difference in days
	nights := int(checkOut.Sub(checkIn).Hours() / 24)
	if nights < 1 {
		nights = 1
	}
	total := priceCents * nights * rooms

	res, err := tx.ExecContext(ctx, "INSERT INTO bookings (user_id, hotel_id, check_in, check_out, adults, children, rooms, total_price_cents, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", userID, hotelID, checkIn, checkOut, adults, children, rooms, total, "pending")
	if err != nil {
		return model.Booking{}, err
	}
	id64, _ := res.LastInsertId()

	result, err := tx.ExecContext(ctx, "UPDATE hotels SET rooms_available = rooms_available - ? WHERE id = ? AND rooms_available >= ?", rooms, hotelID, rooms)
	if err != nil {
		return model.Booking{}, err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return model.Booking{}, ErrNotEnoughRooms
	}

	if err := tx.Commit(); err != nil {
		return model.Booking{}, err
	}

	return model.Booking{
		ID:              int(id64),
		UserID:          userID,
		HotelID:         hotelID,
		CheckIn:         checkIn,
		CheckOut:        checkOut,
		Adults:          adults,
		Children:        children,
		Rooms:           rooms,
		TotalPriceCents: total,
		Status:          "pending",
	}, nil
}

func (r *mysqlBookingRepo) ListByUserID(userID int) ([]model.Booking, error) {
	rows, err := r.db.Query("SELECT id, user_id, hotel_id, check_in, check_out, adults, children, rooms, total_price_cents, status, created_at FROM bookings WHERE user_id = ? ORDER BY created_at DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]model.Booking, 0)
	for rows.Next() {
		var b model.Booking
		if err := rows.Scan(&b.ID, &b.UserID, &b.HotelID, &b.CheckIn, &b.CheckOut, &b.Adults, &b.Children, &b.Rooms, &b.TotalPriceCents, &b.Status, &b.CreatedAt); err != nil {
			continue
		}
		out = append(out, b)
	}
	return out, nil
}
