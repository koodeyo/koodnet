package models

import (
	"net/http"

	"gorm.io/gorm"
)

var Errors = map[error]struct {
	Status  int
	Code    string
	Message string
}{
	gorm.ErrRecordNotFound: {
		Status:  http.StatusNotFound,
		Code:    "ERR_NOT_FOUND",
		Message: "The requested resource could not be found.",
	},
	gorm.ErrInvalidTransaction: {
		Status:  http.StatusInternalServerError,
		Code:    "ERR_INVALID_TRANSACTION",
		Message: "The database transaction is invalid. Please try again.",
	},
	gorm.ErrNotImplemented: {
		Status:  http.StatusNotImplemented,
		Code:    "ERR_NOT_IMPLEMENTED",
		Message: "The requested operation is not yet implemented.",
	},
	gorm.ErrMissingWhereClause: {
		Status:  http.StatusBadRequest,
		Code:    "ERR_MISSING_WHERE",
		Message: "A WHERE clause is required for this query.",
	},
	gorm.ErrUnsupportedRelation: {
		Status:  http.StatusBadRequest,
		Code:    "ERR_UNSUPPORTED_RELATION",
		Message: "This database relation is not supported.",
	},
	gorm.ErrPrimaryKeyRequired: {
		Status:  http.StatusBadRequest,
		Code:    "ERR_PRIMARY_KEY_REQUIRED",
		Message: "A primary key is required for this operation.",
	},
	gorm.ErrModelValueRequired: {
		Status:  http.StatusBadRequest,
		Code:    "ERR_MODEL_VALUE_REQUIRED",
		Message: "A model value is required for this operation.",
	},
	gorm.ErrModelAccessibleFieldsRequired: {
		Status:  http.StatusBadRequest,
		Code:    "ERR_MODEL_ACCESSIBLE_FIELDS_REQUIRED",
		Message: "Accessible fields must be defined for the model.",
	},
	gorm.ErrSubQueryRequired: {
		Status:  http.StatusBadRequest,
		Code:    "ERR_SUBQUERY_REQUIRED",
		Message: "A subquery is required for this operation.",
	},
	gorm.ErrInvalidData: {
		Status:  http.StatusBadRequest,
		Code:    "ERR_INVALID_DATA",
		Message: "The provided data is invalid or unsupported.",
	},
	gorm.ErrUnsupportedDriver: {
		Status:  http.StatusInternalServerError,
		Code:    "ERR_UNSUPPORTED_DRIVER",
		Message: "The database driver is unsupported.",
	},
	gorm.ErrRegistered: {
		Status:  http.StatusConflict,
		Code:    "ERR_ALREADY_REGISTERED",
		Message: "The requested resource is already registered.",
	},
	gorm.ErrInvalidField: {
		Status:  http.StatusBadRequest,
		Code:    "ERR_INVALID_FIELD",
		Message: "An invalid field was provided.",
	},
	gorm.ErrEmptySlice: {
		Status:  http.StatusBadRequest,
		Code:    "ERR_EMPTY_SLICE",
		Message: "An empty slice was provided where data was expected.",
	},
	gorm.ErrDryRunModeUnsupported: {
		Status:  http.StatusNotImplemented,
		Code:    "ERR_DRY_RUN_UNSUPPORTED",
		Message: "Dry-run mode is not supported for this operation.",
	},
	gorm.ErrInvalidDB: {
		Status:  http.StatusInternalServerError,
		Code:    "ERR_INVALID_DB",
		Message: "The database is invalid or corrupted.",
	},
	gorm.ErrInvalidValue: {
		Status:  http.StatusBadRequest,
		Code:    "ERR_INVALID_VALUE",
		Message: "The provided value is invalid. It must be a pointer to a struct or slice.",
	},
	gorm.ErrInvalidValueOfLength: {
		Status:  http.StatusBadRequest,
		Code:    "ERR_INVALID_VALUE_LENGTH",
		Message: "The provided values' length does not match the expected length.",
	},
	gorm.ErrPreloadNotAllowed: {
		Status:  http.StatusBadRequest,
		Code:    "ERR_PRELOAD_NOT_ALLOWED",
		Message: "Preloading is not allowed when count is used.",
	},
	gorm.ErrDuplicatedKey: {
		Status:  http.StatusConflict,
		Code:    "ERR_DUPLICATED_KEY",
		Message: "A resource with this identifier already exists.",
	},
	gorm.ErrForeignKeyViolated: {
		Status:  http.StatusBadRequest,
		Code:    "ERR_FOREIGN_KEY_VIOLATED",
		Message: "The operation violates foreign key constraints.",
	},
	gorm.ErrCheckConstraintViolated: {
		Status:  http.StatusBadRequest,
		Code:    "ERR_CHECK_CONSTRAINT_VIOLATED",
		Message: "The operation violates a database check constraint.",
	},
}

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

func NewValidationError(message string) ValidationError {
	return ValidationError{Message: message}
}
