package model

import (
	"database/sql/driver"
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Package struct {
	Id          int     `db:"id" json:"id"`
	Description string  `db:"description" json:"description"`
	Price       float64 `db:"price" json:"price"`
	Model
}

func (p Package) Validate()(error, map[string]string) {
	return p.Model.ValidationRules(&p, 
	validation.Field(&p.Description, validation.Required.Error("Description is required."), validation.Length(1, 0).Error("Description is required.")), 
	validation.Field(&p.Price, validation.Required.Error("Price is required."), validation.Min(float64(1)).Error("Price is required.")), )
}


type PackageJSON struct {
	Package
}

func (instance *PackageJSON) Scan(value interface{}) error {
	val, valid := value.([]byte)
	if valid {
		unmarshalErr := json.Unmarshal(val, instance)
		if unmarshalErr != nil {
			*instance = PackageJSON{}
		}
	} else {
		*instance = PackageJSON{}
	}
	return nil

}
func (copy PackageJSON) Value(value interface{}) (driver.Value, error) {
	return copy, nil
}