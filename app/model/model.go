package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type Model struct {}

func (m *Model) ValidationRules(structPtr interface{}, fields ...*validation.FieldRules) (error, map[string]string) {
	fieldErrs := make(map[string]string)
	validationErrs :=  validation.ValidateStruct(structPtr,  fields...)
	if validationErrs != nil  {
		_, isValidationErrors:= validationErrs.(validation.Errors)
		if isValidationErrors{
			for k, validationErr := range validationErrs.(validation.Errors) {
				fieldErrs[k] = validationErr.Error()
			}
		}
		
	}
	return validationErrs, fieldErrs
}
