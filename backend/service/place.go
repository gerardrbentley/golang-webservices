package service

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
)

type PlaceRecord struct {
	FoodResourceType  string `db:"food_resource_type" json:"-"`
	Agency            string `db:"agency" json:"name"`
	Location          string `db:"location" json:"location"`
	OperationalStatus string `db:"operational_status" json:"-"`
	OperationalNotes  string `db:"operational_notes" json:"notes"`
	WhoTheyServe      string `db:"who_they_serve" json:"-"`
	Address           string `db:"address" json:"address"`
	Latitude          string `db:"latitude" json:"latitude"`
	Longitude         string `db:"longitude" json:"longitude"`
	PhoneNumber       string `db:"phone_number" json:"phone_number"`
	Website           string `db:"website" json:"website"`
	DaysOrHours       string `db:"days_or_hours" json:"days_or_hours"`
	DateUpdated       string `db:"date_updated" json:"-"`
}

type PlaceService interface {
	LookupByName(searchName string) ([]PlaceRecord, error)
}

type DbPool interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

type PgPlaceService struct {
	dbpool DbPool
}

func NewPgPlaceService(dbpool DbPool) *PgPlaceService {
	return &PgPlaceService{dbpool: dbpool}
}

func (s *PgPlaceService) LookupByName(searchName string) ([]PlaceRecord, error) {
	rows, err := s.dbpool.Query(context.Background(),
		`select "food_resource_type",
			"agency",
			"location",
			"operational_status",
			"operational_notes",
			"who_they_serve",
			"address",
			"latitude",
			"longitude",
			"phone_number",
			"website",
			"days_or_hours",
			"date_updated" 
		from place where tsv @@ plainto_tsquery($1);`,
		searchName)
	if err != nil {
		log.Printf("Query failed: %v\n", err)
		return []PlaceRecord{}, errors.New("Database Error")
	}
	records, err := pgx.CollectRows(rows, pgx.RowToStructByName[PlaceRecord])
	if err == pgx.ErrNoRows || len(records) == 0 {
		return []PlaceRecord{}, errors.New("Not Found")
	} else if err != nil {
		log.Printf("Parsing Record failed: %v\n", err)
		return []PlaceRecord{}, errors.New("Database Record Corrupted")
	}
	return records, nil
}

type InMemoryPlaceService struct {
	LookupByNameFunc func(searchName string) ([]PlaceRecord, error)
}

func (s *InMemoryPlaceService) LookupByName(searchName string) ([]PlaceRecord, error) {
	return s.LookupByNameFunc(searchName)
}

type InMemoryDbPool struct {
	QueryFunc func(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func (p InMemoryDbPool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return p.QueryFunc(ctx, sql, args)
}
