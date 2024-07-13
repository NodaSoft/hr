package main

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
)

// TaskProcessingConfig contains properties for configuring creating, processing and processing of processed tasks.
type TaskProcessingConfig struct {
	// How long does it take to generate tasks.
	GeneratingTasksDuration int `validate:"required,gt=0"`

	// The maximum time that a worker processing tasks can work.
	// In order for the workers not to freeze in case of an error.
	MaxProcessingWorkerDuration int `validate:"required,gt=0"`

	// The maximum time that a worker handling processed tasks can work.
	// In order for the workers not to freeze in case of an error.
	MaxHandleProcessedDuration int `validate:"required,gt=0"`

	// Size of the channel that is used to store unprocessed tasks.
	UnprocessedTasksChannelBufferSize int `validate:"required,gte=0"`

	// Size of the channel that is used to store processed tasks.
	ProcessedTasksChannelBufferSize int `validate:"required,gte=0"`

	// The number of workers that are used to creating new tasks.
	FillingTasksChannelWorkerCount int `validate:"required,gt=0"`

	// The number of workers that are used to processing tasks.
	ProcessingTasksChannelWorkerCount int `validate:"required,gt=0"`

	// Frequency of output of statistics on task processing.
	PrintingTasksPeriod int `validate:"required,gt=0"`

	// Whether to display detailed information on each of the tasks.
	// If set to true, each processed string will be output in the format `+v`.
	IsPrintTasksDetailed bool

	// The level of the logging.
	LogLevel int8 `validate:"gte=-1,lte=5"`
}

// ErrInvalidFieldsValues indicate that at least one of the following fields has invalid value.
var ErrInvalidFieldsValues = errors.New("invalid values of few fields")

// ValidateTaskProcessingConfig validates TaskProcessingConfig fields.
func ValidateTaskProcessingConfig(cfg TaskProcessingConfig) error {
	validate := validator.New()

	err := validate.Struct(cfg)
	if err != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			fmt.Println(err)
			return invalidValidationError
		}

		var fieldErrors []validator.FieldError
		for _, fieldErr := range err.(validator.ValidationErrors) {
			fmt.Printf("Invalid field %s. Error: %s\n", fieldErr.Field(), fieldErr.Error())

			fieldErrors = append(fieldErrors, fieldErr)
		}

		if len(fieldErrors) > 0 {
			return ErrInvalidFieldsValues
		}
	}

	return nil
}
