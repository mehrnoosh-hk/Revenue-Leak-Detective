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
package db

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// ConvertUUIDToPgtypeUUID converts a uuid.UUID to a pgtype.UUID,
// which is the type expected by pgx and SQLC for UUID columns.
//
// Parameters:
//   - id: The uuid.UUID value to convert.
//
// Returns:
//   - pgtype.UUID: The corresponding pgtype.UUID with Valid set to true.
func ConvertUUIDToPgtypeUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

// ConvertNullableUUIDToPgtypeUUID converts a pointer to uuid.UUID to a pgtype.UUID.
// If the input is nil, returns a pgtype.UUID with Valid set to false.
//
// Parameters:
//   - id: Pointer to uuid.UUID (may be nil).
//
// Returns:
//   - pgtype.UUID: The corresponding pgtype.UUID, with Valid set appropriately.
func ConvertNullableUUIDToPgtypeUUID(id *uuid.UUID) pgtype.UUID {
	if id == nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: *id, Valid: true}
}

// ConvertPgtypeUUIDToUUID converts a pgtype.UUID to a uuid.UUID.
// If the pgtype.UUID is not valid, returns uuid.Nil.
//
// Parameters:
//   - pgUUID: The pgtype.UUID to convert.
//
// Returns:
//   - uuid.UUID: The corresponding uuid.UUID, or uuid.Nil if invalid.
func ConvertPgtypeUUIDToUUID(pgUUID pgtype.UUID) uuid.UUID {
	if !pgUUID.Valid {
		return uuid.Nil
	}
	return pgUUID.Bytes
}

// ConvertInterfaceToBytes safely converts an input of type any (interface{})
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
func ConvertInterfaceToBytes(data any) ([]byte, error) {
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
	*EventTypeEnum |
		*EventStatusEnum |
		*ActionTypeEnum |
		*ActionStatusEnum |
		*ActionResultEnum |
		*LeakTypeEnum |
		*PaymentTypeEnum |
		*PaymentStatusEnum
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
	NullEventTypeEnum |
		NullEventStatusEnum |
		NullActionTypeEnum |
		NullActionStatusEnum |
		NullActionResultEnum |
		NullLeakTypeEnum |
		NullPaymentTypeEnum |
		NullPaymentStatusEnum
}

// ConvertEnumsToNullableEnum converts a pointer to an enum type to its corresponding
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
//	result, err := ConvertEnumsToNullableEnum[*EventTypeEnum, NullEventTypeEnum](eventTypePtr)
//	if err != nil {
//		// handle error
//	}
//	nullEventType := result
func ConvertEnumsToNullableEnum[T convertableEnums, R convertedEnums](enumPtr T) R {
	switch v := any(enumPtr).(type) {
	case *EventTypeEnum:
		if v == nil {
			//nolint:errcheck // type assertion is safe due to generic constraints
			return any(NullEventTypeEnum{Valid: false}).(R)
		}
		//nolint:errcheck // type assertion is safe due to generic constraints
		return any(NullEventTypeEnum{EventTypeEnum: *v, Valid: true}).(R)
	case *EventStatusEnum:
		if v == nil {
			//nolint:errcheck // type assertion is safe due to generic constraints
			return any(NullEventStatusEnum{Valid: false}).(R)
		}
		//nolint:errcheck // type assertion is safe due to generic constraints
		return any(NullEventStatusEnum{EventStatusEnum: *v, Valid: true}).(R)
	case *ActionTypeEnum:
		if v == nil {
			//nolint:errcheck // type assertion is safe due to generic constraints
			return any(NullActionTypeEnum{Valid: false}).(R)
		}
		//nolint:errcheck // type assertion is safe due to generic constraints
		return any(NullActionTypeEnum{ActionTypeEnum: *v, Valid: true}).(R)
	case *ActionStatusEnum:
		if v == nil {
			//nolint:errcheck // type assertion is safe due to generic constraints
			return any(NullActionStatusEnum{Valid: false}).(R)
		}
		//nolint:errcheck // type assertion is safe due to generic constraints
		return any(NullActionStatusEnum{ActionStatusEnum: *v, Valid: true}).(R)
	case *ActionResultEnum:
		if v == nil {
			//nolint:errcheck // type assertion is safe due to generic constraints
			return any(NullActionResultEnum{Valid: false}).(R)
		}
		//nolint:errcheck // type assertion is safe due to generic constraints
		return any(NullActionResultEnum{ActionResultEnum: *v, Valid: true}).(R)
	case *LeakTypeEnum:
		if v == nil {
			//nolint:errcheck // type assertion is safe due to generic constraints
			return any(NullLeakTypeEnum{Valid: false}).(R)
		}
		//nolint:errcheck // type assertion is safe due to generic constraints
		return any(NullLeakTypeEnum{LeakTypeEnum: *v, Valid: true}).(R)
	case *PaymentTypeEnum:
		if v == nil {
			//nolint:errcheck // type assertion is safe due to generic constraints
			return any(NullPaymentTypeEnum{Valid: false}).(R)
		}
		//nolint:errcheck // type assertion is safe due to generic constraints
		return any(NullPaymentTypeEnum{PaymentTypeEnum: *v, Valid: true}).(R)
	case *PaymentStatusEnum:
		if v == nil {
			//nolint:errcheck // type assertion is safe due to generic constraints
			return any(NullPaymentStatusEnum{Valid: false}).(R)
		}
		//nolint:errcheck // type assertion is safe due to generic constraints
		return any(NullPaymentStatusEnum{PaymentStatusEnum: *v, Valid: true}).(R)
	default:
		var zero R
		return zero
	}
}

// ConvertNullablePgtypeUUIDToUUID converts a pgtype.UUID to a pointer to uuid.UUID.
// If the pgtype.UUID is not valid, returns nil.
//
// Parameters:
//   - pgUUID: The pgtype.UUID to convert.
//
// Returns:
//   - *uuid.UUID: Pointer to the corresponding uuid.UUID, or nil if invalid.
func ConvertNullablePgtypeUUIDToUUID(pgUUID pgtype.UUID) *uuid.UUID {
	if !pgUUID.Valid {
		return nil
	}
	// Convert [16]byte to uuid.UUID
	uuidValue := uuid.UUID(pgUUID.Bytes)
	return &uuidValue
}

// ConvertEnumToNullableEnum converts a non-nullable enum to its nullable equivalent.
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
func ConvertEnumToNullableEnum(enumValue any) (any, error) {
	switch v := enumValue.(type) {
	case EventTypeEnum:
		return NullEventTypeEnum{EventTypeEnum: v, Valid: true}, nil
	case EventStatusEnum:
		return NullEventStatusEnum{EventStatusEnum: v, Valid: true}, nil
	case ActionTypeEnum:
		return NullActionTypeEnum{ActionTypeEnum: v, Valid: true}, nil
	case ActionStatusEnum:
		return NullActionStatusEnum{ActionStatusEnum: v, Valid: true}, nil
	case ActionResultEnum:
		return NullActionResultEnum{ActionResultEnum: v, Valid: true}, nil
	case LeakTypeEnum:
		return NullLeakTypeEnum{LeakTypeEnum: v, Valid: true}, nil
	case PaymentTypeEnum:
		return NullPaymentTypeEnum{PaymentTypeEnum: v, Valid: true}, nil
	case PaymentStatusEnum:
		return NullPaymentStatusEnum{PaymentStatusEnum: v, Valid: true}, nil
	default:
		return nil, errors.New("ConvertEnumToNullableEnum: unsupported enum type")
	}
}
