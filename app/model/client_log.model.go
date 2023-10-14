package model

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
)

type ClientLog struct {
	Id       int  `json:"id" db:"id"`
	ClientId int  `json:"clientId" db:"client_id"`
	IsMember bool `json:"isMember" db:"is_member"`
	Client ClientJSON `json:"client" db:"client"`
	AmountPaid float64 	`json:"amountPaid" db:"amount_paid"`
	CreatedAt string `json:"createdAt,omitempty" db:"created_at"`
	Model
}

func (m ClientLog) Validate() (error, map[string]string) {
	return m.Model.ValidationRules(&m, 
		validation.Field(&m.ClientId, validation.Required.Error("Client is required"), validation.Min(1).Error("Client is required.")),
		validation.Field(&m.AmountPaid, validation.By(func(value interface{}) error {	
				if m.IsMember{
					return nil
				}
				amountPaid,ok := value.(float64)
				if (!ok){
					return fmt.Errorf("Amount is required.")
				}
				if amountPaid == 0 {
					return fmt.Errorf("Amount is required.")
				}
				return nil
		})),
	)
}
