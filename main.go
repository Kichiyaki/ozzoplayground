package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Task struct {
	ID          string `json:"id"`
	IsDone      bool   `json:"isDone"`
	Description string `json:"description"`
	SubTasks    []Task `json:"subTasks"`
	Something   []int  `json:"something"`
}

func (t Task) Validate() error {
	if err := validation.ValidateStruct(&t,
		validation.Field(&t.Description, validation.Required, validation.Length(5, 100)),
		validation.Field(&t.SubTasks),
		validation.Field(&t.Something, validation.Each(validation.Min(25))),
	); err != nil {
		return err
	}

	return nil
}

func main() {
	t := Task{
		ID:          "id",
		IsDone:      false,
		Description: "123",
		SubTasks: []Task{
			{
				ID: "subtask",
			},
		},
		Something: []int{35, 45, 20, 44, 15},
	}

	err := t.Validate()

	var validationErrors validation.Errors
	if errors.As(err, &validationErrors) {
		errs := newValidationErrors(validationErrors, "")
		b, _ := json.MarshalIndent(errs, "", "   ")
		fmt.Println(string(b))
	}
}

type ValidationError struct {
	Field   string                 `json:"field"`
	Code    string                 `json:"code"`
	Params  map[string]interface{} `json:"params,omitempty"`
	Message string                 `json:"message"`
}

func newValidationError(field string, err validation.Error) ValidationError {
	return ValidationError{
		Field:   field,
		Code:    err.Code(),
		Params:  err.Params(),
		Message: err.Error(),
	}
}

func newValidationErrors(errs validation.Errors, prefix string) []ValidationError {
	var result []ValidationError

	for key, err := range errs {
		field := joinPrefixAndKey(prefix, key)

		switch t := err.(type) {
		case validation.Errors:
			result = append(result, newValidationErrors(t, field)...)
		case validation.Error:
			result = append(result, newValidationError(field, t))
		}
	}

	return result
}

func joinPrefixAndKey(prefix, key string) string {
	if prefix == "" {
		return key
	}

	if _, err := strconv.Atoi(key); err == nil {
		return prefix + "[" + key + "]"
	}

	return prefix + "." + key
}
