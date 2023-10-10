package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
)

type WorkoutCategory struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Workouts Workouts `json:"workouts" db:"workouts"`
	Model
}

func (m WorkoutCategory) Validate() (error, map[string]string) {
	err, fields := m.Model.ValidationRules(&m, 
		validation.Field(&m.Name, validation.Required.Error("Name is required."), 
		validation.Length(1, 100).Error("Name must be atleast 1 to 100 characters long.")),
	)
	
	if (len(m.Workouts) == 0){	
		fields["workouts"] = "Workouts are required."
		return fmt.Errorf("Validation err"), fields  
	}
	for _, workout := range m.Workouts {
		err, _ := m.Model.ValidationRules(&workout, validation.Field(&workout.Id, validation.Min(1).Error("Workouts is required.")) )
		if(err != nil ){
			fields["workouts"] = "Workouts are required."
			break;
		}
	}
	return err, fields
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

type WorkoutJSON struct {
	Workout
}

type Workouts []Workout

func (instance *Workouts) Scan(value interface{}) error {
	val, valid := value.([]byte)
	workouts := make(Workouts, 0)
	if valid {
		unmarshalErr := json.Unmarshal(val, instance)
		if unmarshalErr != nil {
			*instance = workouts
		}
	} else {
		*instance = workouts
	}
	return nil

}
func (copy Workouts) Value(value interface{}) (driver.Value, error) {
	return copy, nil
}

func (instance *WorkoutJSON) Scan(value interface{}) error {
	val, valid := value.([]byte)
	fmt.Println(valid)
	if valid {
		unmarshalErr := json.Unmarshal(val, instance)
		if unmarshalErr != nil {
			*instance = WorkoutJSON{}
		}
	} else {
		*instance = WorkoutJSON{}
	}
	return nil

}
func (copy WorkoutJSON) Value(value interface{}) (driver.Value, error) {
	return copy, nil
}