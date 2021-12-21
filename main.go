package main

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"log"
	"strconv"
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

	printErrors(t.Validate())
}

func printErrors(err error) {
	if err == nil {
		return
	}

	var validationErrors validation.Errors
	if !errors.As(err, &validationErrors) {
		return
	}

	printValidationErrors(validationErrors, "")
}

func printValidationErrors(errs validation.Errors, prefix string) {
	for key, err := range errs {
		prefixAndKey := joinPrefixAndKey(prefix, key)

		switch t := err.(type) {
		case validation.Errors:
			printValidationErrors(t, prefixAndKey)
		case validation.Error:
			log.Println(prefixAndKey, t.Error(), t.Code(), t.Params())
		}
	}
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
