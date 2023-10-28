package model

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type DateSlot struct {
	Id int `json:"id" db:"id" `
	Date string `json:"date" db:"date"`
}

type DateRangeBody struct {
	From string `json:"from"`
	To string `json:"to"`
	Model
}
func (m DateRangeBody) Validate() (error, map[string]string){
	return m.Model.ValidationRules(&m, validation.Field(&m.From, validation.By(func(value interface{}) error {
	
		strDate ,_ := value.(string)
		from, err := time.Parse(time.DateOnly, strDate)
		if err != nil {
			return fmt.Errorf("invalid date format")
		}
		location, err := time.LoadLocation("Asia/Manila")
		if err != nil {
			return fmt.Errorf("invalid date format")
		}
		//use philippines time
		today := time.Now().In(location)
		from = from.In(location)
		//remove time, compare only date
		today = today.Truncate(24 * time.Hour)
		from = from.Truncate(24 * time.Hour)
		if(from.Before(today)){
			return fmt.Errorf("field shoud be greater than or equal date today")
		}
		return nil
	})), validation.Field(&m.To, validation.By(func(value interface{}) error {
	
		strDate ,_ := value.(string)
		to, err := time.Parse(time.DateOnly, strDate)
		if err != nil {
			return fmt.Errorf("invalid date format")
		}
		location, err := time.LoadLocation("Asia/Manila")
		if err != nil {
			return fmt.Errorf("invalid date format")
		}
		//use philippines time
		today := time.Now().In(location)
		to = to.In(location)
		//remove time, compare only date
		today = today.Truncate(24 * time.Hour)
		to = to.Truncate(24 * time.Hour)
		
		
		if(to.Before(today)){
			return fmt.Errorf("field should be greater than or equal date today")
		}
		from, _ := time.Parse(time.DateOnly, m.From)
		if(to.Before(from)){
			return fmt.Errorf("field should be greater than 'from' field")
		}
		return nil
	})) )
}
