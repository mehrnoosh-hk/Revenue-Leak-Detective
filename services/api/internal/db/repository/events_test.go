package repository

import (
	"encoding/json"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"

	db "rdl-api/internal/db/sqlc"
	"rdl-api/internal/domain/models"
)

// Test helper functions
func createTestEvent() models.Event {
	data := json.RawMessage(`{"amount": 100.50, "currency": "USD"}`)

	return models.Event{
		ID:         uuid.New(),
		TenantID:   uuid.New(),
		ProviderID: uuid.New(),
		EventType:  models.EventTypeEnumPaymentFailed,
		EventID:    "evt_test_123",
		Status:     models.EventStatusEnumPending,
		Data:       &data,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func createTestCreateEventParams() models.CreateEventParams {
	return models.CreateEventParams{
		TenantID:   uuid.New(),
		ProviderID: uuid.New(),
		EventType:  models.EventTypeEnumPaymentFailed,
		EventID:    "evt_test_123",
		Status:     models.EventStatusEnumPending,
		Data:       `{"amount": 100.50, "currency": "USD"}`,
	}
}

func createTestUpdateEventParams() models.UpdateEventParams {
	eventType := models.EventTypeEnumPaymentSucceeded
	status := models.EventStatusEnumProcessed
	data := json.RawMessage(`{"amount": 200.75, "currency": "EUR"}`)
	return models.UpdateEventParams{
		ID:        uuid.New(),
		EventType: &eventType,
		Status:    &status,
		Data:      &data,
	}
}

func createTestDBEvent(t *testing.T) db.Event { //nolint: unused
	now := time.Now()
	data, err := json.Marshal(map[string]interface{}{"amount": 100.50, "currency": "USD"})
	if err != nil {
		t.Fatalf("json.Marshal failed in createTestDBEvent: %v", err)
	}

	return db.Event{
		ID:         convertUUIDToPgtypeUUID(uuid.New()),
		TenantID:   convertUUIDToPgtypeUUID(uuid.New()),
		ProviderID: convertUUIDToPgtypeUUID(uuid.New()),
		EventType:  db.EventTypeEnumPaymentFailed,
		EventID:    "evt_test_123",
		Status:     db.EventStatusEnumPending,
		Data:       data,
		CreatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
	}
}

// createTestDBEventForBenchmark creates a test DB event for benchmarks
// Uses panic for benchmarks since we can't use t.Fatal in benchmark context
func createTestDBEventForBenchmark() db.Event {
	now := time.Now()
	data, err := json.Marshal(map[string]interface{}{"amount": 100.50, "currency": "USD"})
	if err != nil {
		panic("json.Marshal failed in createTestDBEventForBenchmark: " + err.Error())
	}

	return db.Event{
		ID:         convertUUIDToPgtypeUUID(uuid.New()),
		TenantID:   convertUUIDToPgtypeUUID(uuid.New()),
		ProviderID: convertUUIDToPgtypeUUID(uuid.New()),
		EventType:  db.EventTypeEnumPaymentFailed,
		EventID:    "evt_test_123",
		Status:     db.EventStatusEnumPending,
		Data:       data,
		CreatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
	}
}

func createTestLogger() *slog.Logger {
	// Use a no-op handler for testing to avoid nil writer issues
	return slog.New(slog.NewTextHandler(&noOpWriter{}, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

// noOpWriter is a no-op writer for testing
type noOpWriter struct{}

func (w *noOpWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

// TestConversionFunctions tests the conversion helper functions
func TestToEventDomain(t *testing.T) {
	tests := []struct {
		name     string
		input    func(t *testing.T) db.Event
		expected models.Event
	}{
		{
			name: "successful conversion",
			input: func(t *testing.T) db.Event {
				now := time.Now()
				data, err := json.Marshal(map[string]interface{}{"amount": 100.50, "currency": "USD"})
				if err != nil {
					t.Fatalf("json.Marshal failed in TestToEventDomain successful conversion: %v", err)
				}
				return db.Event{
					ID:         convertUUIDToPgtypeUUID(uuid.New()),
					TenantID:   convertUUIDToPgtypeUUID(uuid.New()),
					ProviderID: convertUUIDToPgtypeUUID(uuid.New()),
					EventType:  db.EventTypeEnumPaymentFailed,
					EventID:    "evt_test_123",
					Status:     db.EventStatusEnumPending,
					Data:       data,
					CreatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
					UpdatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
				}
			},
		},
		{
			name: "conversion with nil timestamps",
			input: func(t *testing.T) db.Event {
				data, err := json.Marshal(map[string]interface{}{"amount": 100.50})
				if err != nil {
					t.Fatalf("json.Marshal failed in TestToEventDomain nil timestamps: %v", err)
				}
				return db.Event{
					ID:         convertUUIDToPgtypeUUID(uuid.New()),
					TenantID:   convertUUIDToPgtypeUUID(uuid.New()),
					ProviderID: convertUUIDToPgtypeUUID(uuid.New()),
					EventType:  db.EventTypeEnumPaymentSucceeded,
					EventID:    "evt_test_456",
					Status:     db.EventStatusEnumProcessed,
					Data:       data,
					CreatedAt:  pgtype.Timestamptz{Valid: false},
					UpdatedAt:  pgtype.Timestamptz{Valid: false},
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputEvent := tt.input(t)
			result := toEventDomain(inputEvent)

			// We can't directly compare the structs because of the UUID fields
			// So we compare individual fields
			assert.Equal(t, inputEvent.EventID, result.EventID)
			assert.Equal(t, models.EventTypeEnum(inputEvent.EventType), result.EventType)
			assert.Equal(t, models.EventStatusEnum(inputEvent.Status), result.Status)

			// Compare timestamps
			if inputEvent.CreatedAt.Valid {
				assert.Equal(t, inputEvent.CreatedAt.Time, result.CreatedAt)
			} else {
				assert.Equal(t, time.Time{}, result.CreatedAt)
			}

			if inputEvent.UpdatedAt.Valid {
				assert.Equal(t, inputEvent.UpdatedAt.Time, result.UpdatedAt)
			} else {
				assert.Equal(t, time.Time{}, result.UpdatedAt)
			}
		})
	}
}

func TestToCreateEventDBParams(t *testing.T) {
	tests := []struct {
		name           string
		input          models.CreateEventParams
		expectedError  error
		validateResult func(t *testing.T, result db.CreateEventParams)
	}{
		{
			name: "successful conversion with string data",
			input: models.CreateEventParams{
				TenantID:   uuid.New(),
				ProviderID: uuid.New(),
				EventType:  models.EventTypeEnumPaymentFailed,
				EventID:    "evt_test_123",
				Status:     models.EventStatusEnumPending,
				Data:       `{"amount": 100.50, "currency": "USD"}`,
			},
			expectedError: nil,
			validateResult: func(t *testing.T, result db.CreateEventParams) {
				assert.Equal(t, "evt_test_123", result.EventID)
				assert.Equal(t, db.EventTypeEnumPaymentFailed, result.EventType)
				assert.Equal(t, db.EventStatusEnumPending, result.Status)
				assert.NotNil(t, result.Data)
			},
		},
		{
			name: "conversion with byte data",
			input: models.CreateEventParams{
				TenantID:   uuid.New(),
				ProviderID: uuid.New(),
				EventType:  models.EventTypeEnumPaymentSucceeded,
				EventID:    "evt_test_456",
				Status:     models.EventStatusEnumProcessed,
				Data:       []byte(`{"payment": {"amount": 200.75, "currency": "EUR"}}`),
			},
			expectedError: nil,
			validateResult: func(t *testing.T, result db.CreateEventParams) {
				assert.Equal(t, "evt_test_456", result.EventID)
				assert.Equal(t, db.EventTypeEnumPaymentSucceeded, result.EventType)
				assert.Equal(t, db.EventStatusEnumProcessed, result.Status)
				assert.NotNil(t, result.Data)
			},
		},
		{
			name: "conversion with nil data should fail",
			input: models.CreateEventParams{
				TenantID:   uuid.New(),
				ProviderID: uuid.New(),
				EventType:  models.EventTypeEnumPaymentRefunded,
				EventID:    "evt_test_789",
				Status:     models.EventStatusEnumFailed,
				Data:       nil,
			},
			expectedError: errors.New("ConvertInterfaceToBytes: unsupported data type, expected string or []byte"),
			validateResult: func(t *testing.T, _ db.CreateEventParams) {
				// This should not be called since we expect an error, it should fail if it is called
				t.Fail()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := toCreateEventDBParams(tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				tt.validateResult(t, result)
			}
		})
	}
}

func TestToUpdateEventDBParams(t *testing.T) {
	tests := []struct {
		name           string
		input          models.UpdateEventParams
		expectedError  error
		validateResult func(t *testing.T, result db.UpdateEventParams)
	}{
		{
			name: "successful conversion with all fields",
			input: func() models.UpdateEventParams {
				eventType := models.EventTypeEnumPaymentSucceeded
				status := models.EventStatusEnumProcessed
				data := json.RawMessage(`{"amount": 200.75, "currency": "EUR"}`)

				return models.UpdateEventParams{
					ID: uuid.New(),

					EventType: &eventType,
					Status:    &status,
					Data:      &data,
				}
			}(),
			expectedError: nil,
			validateResult: func(t *testing.T, result db.UpdateEventParams) {
				assert.NotNil(t, result.ID)
				assert.NotNil(t, result.EventType)
				assert.NotNil(t, result.Status)
				assert.NotNil(t, result.Data)
			},
		},
		{
			name: "conversion with partial fields",
			input: func() models.UpdateEventParams {
				status := models.EventStatusEnumFailed
				data := json.RawMessage(`{"error": "payment failed"}`)
				return models.UpdateEventParams{
					ID:     uuid.New(),
					Status: &status,
					Data:   &data,
				}
			}(),
			expectedError: nil,
			validateResult: func(t *testing.T, result db.UpdateEventParams) {
				assert.NotNil(t, result.ID)
				assert.False(t, result.EventType.Valid)
				assert.True(t, result.Status.Valid)
				assert.NotNil(t, result.Data)
			},
		},
		{
			name: "conversion with nil data",
			input: func() models.UpdateEventParams {
				status := models.EventStatusEnumProcessed
				return models.UpdateEventParams{
					ID:     uuid.New(),
					Status: &status,
					Data:   nil,
				}
			}(),
			expectedError: nil,
			validateResult: func(t *testing.T, result db.UpdateEventParams) {
				assert.NotNil(t, result.ID)
				assert.Nil(t, result.Data)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := toUpdateEventDBParams(tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				tt.validateResult(t, result)
			}
		})
	}
}

// TestEdgeCases tests various edge cases and error conditions
func TestEdgeCases(t *testing.T) {
	t.Run("CreateEvent with invalid data type", func(t *testing.T) {
		// Create a struct that can't be converted to bytes
		invalidData := make(chan int)

		params := models.CreateEventParams{
			TenantID:   uuid.New(),
			ProviderID: uuid.New(),
			EventType:  models.EventTypeEnumPaymentFailed,
			EventID:    "evt_test_123",
			Status:     models.EventStatusEnumPending,
			Data:       invalidData,
		}

		_, err := toCreateEventDBParams(params)
		assert.Error(t, err)
		assert.Equal(t, errors.New("ConvertInterfaceToBytes: unsupported data type, expected string or []byte"), err)
	})

	t.Run("CreateEvent with empty event ID", func(t *testing.T) {
		params := models.CreateEventParams{
			TenantID:   uuid.New(),
			ProviderID: uuid.New(),
			EventType:  models.EventTypeEnumPaymentFailed,
			EventID:    "", // Empty event ID
			Status:     models.EventStatusEnumPending,
			Data:       `{"amount": 100.50}`,
		}

		result, err := toCreateEventDBParams(params)
		assert.NoError(t, err)
		assert.Equal(t, "", result.EventID)
	})
}

// TestDataValidation tests data validation scenarios
func TestDataValidation(t *testing.T) {
	t.Run("valid event types", func(t *testing.T) {
		validTypes := []models.EventTypeEnum{
			models.EventTypeEnumPaymentFailed,
			models.EventTypeEnumPaymentSucceeded,
			models.EventTypeEnumPaymentRefunded,
			models.EventTypeEnumPaymentUpdated,
		}

		for _, eventType := range validTypes {
			params := models.CreateEventParams{
				TenantID:   uuid.New(),
				ProviderID: uuid.New(),
				EventType:  eventType,
				EventID:    "evt_test",
				Status:     models.EventStatusEnumPending,
				Data:       `{"test": "data"}`,
			}

			result, err := toCreateEventDBParams(params)
			assert.NoError(t, err)
			assert.Equal(t, db.EventTypeEnum(eventType), result.EventType)
		}
	})

	t.Run("valid event statuses", func(t *testing.T) {
		validStatuses := []models.EventStatusEnum{
			models.EventStatusEnumPending,
			models.EventStatusEnumProcessed,
			models.EventStatusEnumFailed,
		}

		for _, status := range validStatuses {
			params := models.CreateEventParams{
				TenantID:   uuid.New(),
				ProviderID: uuid.New(),
				EventType:  models.EventTypeEnumPaymentFailed,
				EventID:    "evt_test",
				Status:     status,
				Data:       `{"test": "data"}`,
			}

			result, err := toCreateEventDBParams(params)
			assert.NoError(t, err)
			assert.Equal(t, db.EventStatusEnum(status), result.Status)
		}
	})
}

// TestUUIDConversion tests UUID conversion scenarios
func TestUUIDConversion(t *testing.T) {
	t.Run("valid UUID conversion", func(t *testing.T) {
		testUUID := uuid.New()
		pgtypeUUID := convertUUIDToPgtypeUUID(testUUID)

		// Verify the conversion preserves the UUID
		assert.True(t, pgtypeUUID.Valid)
		assert.Equal(t, testUUID, uuid.UUID(pgtypeUUID.Bytes))
	})

	t.Run("nil UUID handling", func(t *testing.T) {
		var testUUID *uuid.UUID
		pgtypeUUID := convertNullableUUIDToPgtypeUUID(testUUID)

		// Verify nil UUID is handled correctly
		assert.False(t, pgtypeUUID.Valid)
	})
}

// Benchmark tests for performance
func BenchmarkToEventDomain(b *testing.B) {
	dbEvent := createTestDBEventForBenchmark()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = toEventDomain(dbEvent)
	}
}

func BenchmarkToCreateEventDBParams(b *testing.B) {
	params := createTestCreateEventParams()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := toCreateEventDBParams(params)
		if err != nil {
			b.Fatalf("toCreateEventDBParams failed: %v", err)
		}
	}
}

func BenchmarkToUpdateEventDBParams(b *testing.B) {
	params := createTestUpdateEventParams()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := toUpdateEventDBParams(params)
		if err != nil {
			b.Fatalf("toUpdateEventDBParams failed: %v", err)
		}
	}
}

// TestPaginationModels tests the pagination models and helper functions
func TestPaginationModels(t *testing.T) {
	t.Run("PaginationParams validation", func(t *testing.T) {
		// Test valid pagination params
		validParams := models.PaginationParams{
			Limit:  10,
			Offset: 0,
		}
		assert.Equal(t, int32(10), validParams.Limit)
		assert.Equal(t, int32(0), validParams.Offset)

		// Test edge cases
		edgeParams := models.PaginationParams{
			Limit:  1,
			Offset: 1000,
		}
		assert.Equal(t, int32(1), edgeParams.Limit)
		assert.Equal(t, int32(1000), edgeParams.Offset)
	})

	t.Run("PaginatedResponse creation", func(t *testing.T) {
		// Test with sample events
		events := []models.Event{
			createTestEvent(),
			createTestEvent(),
		}

		response := models.NewPaginatedResponse(events, 100, 10, 0)

		assert.Equal(t, events, response.Items)
		assert.Equal(t, int64(100), response.TotalCount)
		assert.Equal(t, int32(10), response.Limit)
		assert.Equal(t, int32(0), response.Offset)
		assert.True(t, response.HasNext)      // 0 + 10 < 100
		assert.False(t, response.HasPrevious) // offset = 0
	})

	t.Run("PaginatedResponse edge cases", func(t *testing.T) {
		// Test first page
		firstPage := models.NewPaginatedResponse([]models.Event{}, 50, 10, 0)
		assert.True(t, firstPage.HasNext)
		assert.False(t, firstPage.HasPrevious)

		// Test middle page
		middlePage := models.NewPaginatedResponse([]models.Event{}, 50, 10, 20)
		assert.True(t, middlePage.HasNext)
		assert.True(t, middlePage.HasPrevious)

		// Test last page
		lastPage := models.NewPaginatedResponse([]models.Event{}, 50, 10, 40)
		assert.False(t, lastPage.HasNext) // 40 + 10 >= 50
		assert.True(t, lastPage.HasPrevious)

		// Test exact fit
		exactFit := models.NewPaginatedResponse([]models.Event{}, 30, 10, 20)
		assert.False(t, exactFit.HasNext) // 20 + 10 >= 30
		assert.True(t, exactFit.HasPrevious)
	})

	t.Run("PaginatedResponse with empty results", func(t *testing.T) {
		emptyResponse := models.NewPaginatedResponse([]models.Event{}, 0, 10, 0)
		assert.Empty(t, emptyResponse.Items)
		assert.Equal(t, int64(0), emptyResponse.TotalCount)
		assert.False(t, emptyResponse.HasNext)
		assert.False(t, emptyResponse.HasPrevious)
	})
}

// TestPaginationHelperFunctions tests pagination-related helper functions
func TestPaginationHelperFunctions(t *testing.T) {
	t.Run("NewPaginatedResponse with different types", func(t *testing.T) {
		// Test with Event type
		events := []models.Event{createTestEvent()}
		eventResponse := models.NewPaginatedResponse(events, 1, 10, 0)
		assert.Len(t, eventResponse.Items, 1)
		assert.IsType(t, models.Event{}, eventResponse.Items[0])

		// Test with string type (generic test)
		strings := []string{"test1", "test2"}
		stringResponse := models.NewPaginatedResponse(strings, 2, 10, 0)
		assert.Len(t, stringResponse.Items, 2)
		assert.Equal(t, "test1", stringResponse.Items[0])
	})

	t.Run("Pagination boundary calculations", func(t *testing.T) {
		// Test various pagination scenarios
		testCases := []struct {
			name         string
			totalCount   int64
			limit        int32
			offset       int32
			expectedNext bool
			expectedPrev bool
		}{
			{"first_page", 100, 10, 0, true, false},
			{"middle_page", 100, 10, 50, true, true},
			{"last_page", 100, 10, 90, false, true},
			{"exact_fit", 30, 10, 20, false, true},
			{"single_item", 1, 10, 0, false, false},
			{"empty_result", 0, 10, 0, false, false},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				response := models.NewPaginatedResponse([]models.Event{}, tc.totalCount, tc.limit, tc.offset)
				assert.Equal(t, tc.expectedNext, response.HasNext, "HasNext mismatch for %s", tc.name)
				assert.Equal(t, tc.expectedPrev, response.HasPrevious, "HasPrevious mismatch for %s", tc.name)
			})
		}
	})
}
