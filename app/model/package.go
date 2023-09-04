package model

import validation "github.com/go-ozzo/ozzo-validation"

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


