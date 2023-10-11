package model

import (
	"fmt"
	"lift-fitness-gym/app/pkg/acl"
	"slices"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Role struct {
	Id          int      `json:"id" db:"id"`
	Name        string   `json:"name" db:"name"`
	Permissions []string `json:"permissions" db:"permissions"`
	Model
}

func (m Role) Validate() (error, map[string]string) {
	return m.Model.ValidationRules(&m,
		validation.Field(&m.Name, validation.Required.Error("Name is required."), validation.Length(1, 100).Error("Name should be at least 1 to 100.")),
		validation.Field(&m.Permissions, validation.Each(validation.Required.Error("Permission is required"), validation.By(func(value interface{}) error {
			permission, isStr := value.(string)
			if !isStr {
				return fmt.Errorf("Permission is required.")
			}
			if(!slices.Contains(acl.Permissions, permission)){
				return fmt.Errorf("Permission is required.")
			}
			return nil
		}))),
	)
}
