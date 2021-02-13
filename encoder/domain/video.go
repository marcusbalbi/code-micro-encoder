package domain

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// Video Entidade Principal do Nosso sistema
type Video struct {
	ID         string    `valid:"uuid"`
	ResourceID string    `valid:"notnull"`
	FilePath   string    `valid:"notnull"`
	CreatedAt  time.Time `valid:"-"`
}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

// NewVideo cria um novo Video
func NewVideo() *Video {
	return &Video{}
}

// Validate valida se existe algum erro no objeto
func (v *Video) Validate() error {
	_, err := govalidator.ValidateStruct(v)

	if err != nil {
		return err
	}
	return nil
}
