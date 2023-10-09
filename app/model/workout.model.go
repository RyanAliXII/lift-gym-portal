package model

import validation "github.com/go-ozzo/ozzo-validation"

type WorkoutCategory struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Model
}

func (m WorkoutCategory) Validate() (error, map[string]string) {
	return m.Model.ValidationRules(&m, validation.Field(&m.Name, validation.Required.Error("Name is required."), validation.Length(1, 100).Error("Name must be atleast 1 to 100 characters long.")))
}

type Workout struct {
	Id int `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	ImagePath string `json:"imagePath" db:"image_path"`
	Model
}

func (m Workout) Validate() (error, map[string]string){
	return m.Model.ValidationRules(&m, 
		validation.Field(&m.Name, validation.Required.Error("Name is required."), validation.Length(1,100).Error(" Name should be atleast 1 to 100.")),
		validation.Field(&m.Description, validation.Required.Error("Description is required."), validation.Length(1,0).Error("Description should be atleast 1.")),
	)
}