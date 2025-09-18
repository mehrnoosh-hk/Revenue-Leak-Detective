// Package db provides SQLC generated codes and comprehensive utility functions and type conversions for working with
// database types, especially for bridging between domain models and the types
// required by SQLC-generated code and PostgreSQL driver (pgx).
//
// This file includes helpers for:
//   - Converting UUIDs between uuid.UUID and pgtype.UUID formats
//   - Handling nullable types with proper error handling
//   - Mapping between Go enums and their nullable database representations
//   - Converting between nullable and non-nullable enum types
//   - Safe conversion of interface{} types to byte slices
//
// All conversion functions provide proper error handling and type safety,
// making them essential for ensuring correctness when persisting and
// retrieving data from the database.
//
// Key Features:
//   - Type-safe enum conversions with error handling
//   - Bidirectional conversion between nullable and non-nullable types
//   - Comprehensive UUID conversion utilities
//   - Safe interface{} to []byte conversion
//   - Extensive documentation and usage examples
package repository

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	db "rdl-api/internal/db/sqlc"
	"rdl-api/internal/domain/models"
)

// convertUUIDToPgtypeUUID converts a uuid.UUID to a pgtype.UUID,
// which is the type expected by pgx and SQLC for UUID columns.
//
// Parameters:
//   - id: The uuid.UUID value to convert.
//
// Returns:
//   - pgtype.UUID: The corresponding pgtype.UUID with Valid set to true.
func convertUUIDToPgtypeUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

// convertNullableUUIDToPgtypeUUID converts a pointer to uuid.UUID to a pgtype.UUID.
// If the input is nil, returns a pgtype.UUID with Valid set to false.
//
// Parameters:
//   - id: Pointer to uuid.UUID (may be nil).
//
// Returns:
//   - pgtype.UUID: The corresponding pgtype.UUID, with Valid set appropriately.
func convertNullableUUIDToPgtypeUUID(id *uuid.UUID) pgtype.UUID { //nolint: unused
	if id == nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: *id, Valid: true}
}

// convertPgtypeUUIDToUUID converts a pgtype.UUID to a uuid.UUID.
// If the pgtype.UUID is not valid, returns uuid.Nil.
//
// Parameters:
//   - pgUUID: The pgtype.UUID to convert.
//
// Returns:
//   - uuid.UUID: The corresponding uuid.UUID, or uuid.Nil if invalid.
func convertPgtypeUUIDToUUID(pgUUID pgtype.UUID) uuid.UUID {
	if !pgUUID.Valid {
		return uuid.Nil
	}
	return pgUUID.Bytes
}

// convertInterfaceToBytes safely converts an input of type any (interface{})
// to a byte slice ([]byte). This is useful for serializing data fields
// that may be stored as JSON or binary in the database.
//
// Supported input types:
//   - string: Returns the UTF-8 bytes of the string.
//   - []byte: Returns the byte slice as-is.
//
// Returns an error if the input is not a string or []byte.
//
// Parameters:
//   - data: The value to convert.
//
// Returns:
//   - []byte: The resulting byte slice.
//   - error: An error if the conversion is not possible.
func convertInterfaceToBytes(data any) ([]byte, error) {
	switch v := data.(type) {
	case string:
		return []byte(v), nil
	case []byte:
		return v, nil
	default:
		return nil, errors.New("ConvertInterfaceToBytes: unsupported data type, expected string or []byte")
	}
}

// convertableEnums is a type constraint that matches any pointer to a supported enum type.
// It is used to enable type-safe handling of all supported enum pointer types for conversion
// to their corresponding nullable equivalents in the ConvertEnumsToNullableEnum function.
type convertableEnums interface {
	*db.EventTypeEnum |
		*db.EventStatusEnum |
		*db.ActionTypeEnum |
		*db.ActionStatusEnum |
		*db.ActionResultEnum |
		*db.LeakTypeEnum |
		*db.PaymentTypeEnum |
		*db.PaymentStatusEnum
}

// convertedEnums is a type constraint that matches any of the nullable enum types
// used in the database layer. This enables type-safe conversion from pointer enum types
// to their corresponding nullable enum representations, ensuring compile-time safety
// when working with generic enum conversion functions.
//
// The supported nullable enum types are:
//   - NullEventTypeEnum
//   - NullEventStatusEnum
//   - NullActionTypeEnum
//   - NullActionStatusEnum
//   - NullActionResultEnum
//   - NullLeakTypeEnum
//   - NullPaymentTypeEnum
//   - NullPaymentStatusEnum
type convertedEnums interface {
	db.NullEventTypeEnum |
		db.NullEventStatusEnum |
		db.NullActionTypeEnum |
		db.NullActionStatusEnum |
		db.NullActionResultEnum |
		db.NullLeakTypeEnum |
		db.NullPaymentTypeEnum |
		db.NullPaymentStatusEnum
}

// convertEnumsToNullableEnum converts a pointer to an enum type to its corresponding
// nullable enum type as defined in the db package. This function supports all
// nullable enum types used in the database layer.
//
// If the input is nil, returns the corresponding nullable enum with Valid=false.
// If the input type is not recognized, returns the zero value of the result type and an error.
//
// Type Parameters:
//   - T: A pointer to a supported enum type (see convertableEnums).
//   - R: The corresponding nullable enum type (see convertedEnums).
//
// Parameters:
//   - enumPtr: Pointer to any supported enum type (may be nil).
//
// Returns:
//   - R: The corresponding nullable enum type.
//   - error: An error if the input type is not supported.
//
// Usage example:
//
//	result, err := convertEnumsToNullableEnum[*EventTypeEnum, NullEventTypeEnum](eventTypePtr)
//	if err != nil {
//		// handle error
//	}
//	nullEventType := result
func convertEnumsToNullableEnum[T convertableEnums, R convertedEnums](enumPtr T) (R, error) {
	switch v := any(enumPtr).(type) {
	case *db.EventTypeEnum:
		if v == nil {
			//nolint:errcheck // type assertion is safe due to generic constraints
			return any(db.NullEventTypeEnum{Valid: false}).(R), nil
		}
		//nolint:errcheck // type assertion is safe due to generic constraints
		return any(db.NullEventTypeEnum{EventTypeEnum: *v, Valid: true}).(R), nil
	case *db.EventStatusEnum:
		if v == nil {
			//nolint:errcheck // type assertion is safe due to generic constraints
			return any(db.NullEventStatusEnum{Valid: false}).(R), nil
		}
		//nolint:errcheck // type assertion is safe due to generic constraints
		return any(db.NullEventStatusEnum{EventStatusEnum: *v, Valid: true}).(R), nil
	case *db.ActionTypeEnum:
		if v == nil {
			//nolint:errcheck // type assertion is safe due to generic constraints
			return any(db.NullActionTypeEnum{Valid: false}).(R), nil
		}
		//nolint:errcheck // type assertion is safe due to generic constraints
		return any(db.NullActionTypeEnum{ActionTypeEnum: *v, Valid: true}).(R), nil
	case *db.ActionStatusEnum:
		if v == nil {
			//nolint:errcheck // type assertion is safe due to generic constraints
			return any(db.NullActionStatusEnum{Valid: false}).(R), nil
		}
		//nolint:errcheck // type assertion is safe due to generic constraints
		return any(db.NullActionStatusEnum{ActionStatusEnum: *v, Valid: true}).(R), nil
	case *db.ActionResultEnum:
		if v == nil {
			//nolint:errcheck // type assertion is safe due to generic constraints
			return any(db.NullActionResultEnum{Valid: false}).(R), nil
		}
		//nolint:errcheck // type assertion is safe due to generic constraints
		return any(db.NullActionResultEnum{ActionResultEnum: *v, Valid: true}).(R), nil
	case *db.LeakTypeEnum:
		if v == nil {
			//nolint:errcheck // type assertion is safe due to generic constraints
			return any(db.NullLeakTypeEnum{Valid: false}).(R), nil
		}
		//nolint:errcheck // type assertion is safe due to generic constraints
		return any(db.NullLeakTypeEnum{LeakTypeEnum: *v, Valid: true}).(R), nil
	case *db.PaymentTypeEnum:
		if v == nil {
			//nolint:errcheck // type assertion is safe due to generic constraints
			return any(db.NullPaymentTypeEnum{Valid: false}).(R), nil
		}
		//nolint:errcheck // type assertion is safe due to generic constraints
		return any(db.NullPaymentTypeEnum{PaymentTypeEnum: *v, Valid: true}).(R), nil
	case *db.PaymentStatusEnum:
		if v == nil {
			//nolint:errcheck // type assertion is safe due to generic constraints
			return any(db.NullPaymentStatusEnum{Valid: false}).(R), nil
		}
		//nolint:errcheck // type assertion is safe due to generic constraints
		return any(db.NullPaymentStatusEnum{PaymentStatusEnum: *v, Valid: true}).(R), nil
	default:
		var zero R
		return zero, fmt.Errorf("unsupported enum type: %T", enumPtr)
	}
}

// convertNullablePgtypeUUIDToUUID converts a pgtype.UUID to a pointer to uuid.UUID.
// If the pgtype.UUID is not valid, returns nil.
//
// Parameters:
//   - pgUUID: The pgtype.UUID to convert.
//
// Returns:
//   - *uuid.UUID: Pointer to the corresponding uuid.UUID, or nil if invalid.
func convertNullablePgtypeUUIDToUUID(pgUUID pgtype.UUID) *uuid.UUID { //nolint: unused
	if !pgUUID.Valid {
		return nil
	}
	// Convert [16]byte to uuid.UUID
	uuidValue := uuid.UUID(pgUUID.Bytes)
	return &uuidValue
}

// convertEnumToNullableEnum converts a non-nullable enum to its nullable equivalent.
// This is useful for converting between domain models and database models.
//
// Supported enum types and their conversions:
//   - EventTypeEnum      -> NullEventTypeEnum
//   - EventStatusEnum    -> NullEventStatusEnum
//   - ActionTypeEnum     -> NullActionTypeEnum
//   - ActionStatusEnum   -> NullActionStatusEnum
//   - ActionResultEnum   -> NullActionResultEnum
//   - LeakTypeEnum       -> NullLeakTypeEnum
//   - PaymentTypeEnum    -> NullPaymentTypeEnum
//   - PaymentStatusEnum  -> NullPaymentStatusEnum
//
// Parameters:
//   - enumValue: The enum value to convert.
//
// Returns:
//   - any: The corresponding nullable enum type.
//   - error: An error if the input type is not supported.
func convertEnumToNullableEnum(enumValue any) (any, error) { //nolint: unused
	switch v := enumValue.(type) {
	case db.EventTypeEnum:
		return db.NullEventTypeEnum{EventTypeEnum: v, Valid: true}, nil
	case db.EventStatusEnum:
		return db.NullEventStatusEnum{EventStatusEnum: v, Valid: true}, nil
	case db.ActionTypeEnum:
		return db.NullActionTypeEnum{ActionTypeEnum: v, Valid: true}, nil
	case db.ActionStatusEnum:
		return db.NullActionStatusEnum{ActionStatusEnum: v, Valid: true}, nil
	case db.ActionResultEnum:
		return db.NullActionResultEnum{ActionResultEnum: v, Valid: true}, nil
	case db.LeakTypeEnum:
		return db.NullLeakTypeEnum{LeakTypeEnum: v, Valid: true}, nil
	case db.PaymentTypeEnum:
		return db.NullPaymentTypeEnum{PaymentTypeEnum: v, Valid: true}, nil
	case db.PaymentStatusEnum:
		return db.NullPaymentStatusEnum{PaymentStatusEnum: v, Valid: true}, nil
	default:
		return nil, errors.New("ConvertEnumToNullableEnum: unsupported enum type")
	}
}

// Action enum conversion functions

// convertActionTypeEnumToDB converts domain ActionTypeEnum to database ActionTypeEnum.
func convertActionTypeEnumToDB(enum models.ActionTypeEnum) db.ActionTypeEnum {
	return db.ActionTypeEnum(enum)
}

// convertActionTypeEnumFromDB converts database ActionTypeEnum to domain ActionTypeEnum.
func convertActionTypeEnumFromDB(enum db.ActionTypeEnum) models.ActionTypeEnum {
	return models.ActionTypeEnum(enum)
}

// convertActionStatusEnumToDB converts domain ActionStatusEnum to database ActionStatusEnum.
func convertActionStatusEnumToDB(enum models.ActionStatusEnum) db.ActionStatusEnum {
	return db.ActionStatusEnum(enum)
}

// convertActionStatusEnumFromDB converts database ActionStatusEnum to domain ActionStatusEnum.
func convertActionStatusEnumFromDB(enum db.ActionStatusEnum) models.ActionStatusEnum {
	return models.ActionStatusEnum(enum)
}

// convertActionResultEnumToDB converts domain ActionResultEnum to database ActionResultEnum.
func convertActionResultEnumToDB(enum models.ActionResultEnum) db.ActionResultEnum {
	return db.ActionResultEnum(enum)
}

// convertActionResultEnumFromDB converts database ActionResultEnum to domain ActionResultEnum.
func convertActionResultEnumFromDB(enum db.ActionResultEnum) models.ActionResultEnum {
	return models.ActionResultEnum(enum)
}
