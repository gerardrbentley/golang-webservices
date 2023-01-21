package service

import (
	"context"
	"reflect"
	"testing"
	"unsafe"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type MockRows struct {
	values   [][]any
	visitIdx int
}

func (r *MockRows) Close() {}

func (r *MockRows) Err() error {
	return nil
}

func (r *MockRows) CommandTag() pgconn.CommandTag {
	return pgconn.NewCommandTag("mock")
}

func (r *MockRows) FieldDescriptions() []pgconn.FieldDescription {
	return []pgconn.FieldDescription{}
}

func (r *MockRows) Next() bool {
	return r.visitIdx < len(r.values)
}

type mockScanner struct {
	ptrToStruct any
}

func (r *MockRows) Scan(dest ...any) error {
	v := reflect.ValueOf(dest[0]).Elem().FieldByName("ptrToStruct")
	v = reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem()

	for i, value := range r.values[r.visitIdx] {
		n := v.Elem().Elem().Field(i)
		n.Set(reflect.ValueOf(value))
	}

	r.visitIdx = r.visitIdx + 1

	return nil
}

func (r *MockRows) Values() ([]any, error) {
	return r.values[r.visitIdx], nil
}

func (r *MockRows) RawValues() [][]byte {
	return nil
}

func (r *MockRows) Conn() *pgx.Conn {
	return nil
}

func TestLookupByName(t *testing.T) {
	mockName := "seattle"
	mockPool := InMemoryDbPool{
		QueryFunc: func(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
			pgxrows := MockRows{
				values: [][]any{
					{
						"meal",
						mockName,
					},
					{
						"food bank",
						"Fremont " + mockName,
					},
				},
			}
			return &pgxrows, nil
		},
	}
	s := PgPlaceService{dbpool: mockPool}
	records, _ := s.LookupByName(mockName)
	record := records[0]
	if record.Agency != mockName {
		t.Errorf("Not expected: %v", record.Agency)
	}
}
