package validation

import (
	"errors"
	"testing"
	"time"

	v1resources "github.com/metaclips/LetsTalk/backend/resources/api/v1"
	"github.com/stretchr/testify/assert"
)

func TestValidation(t *testing.T) {

	for _, val := range []struct {
		data v1resources.UserRegistration
		name string
		err  error
	}{
		{
			data: v1resources.UserRegistration{
				Email:      "michael@gmail.com",
				LastName:   "Chike",
				DOB:        time.Now(),
				Department: "Engineering",
				Password:   "test",
			},
			name: "empty value",
			err:  errors.New("firstName is required"),
		},
		{
			data: v1resources.UserRegistration{
				Email:      "mmmmm",
				FirstName:  "Michael",
				LastName:   "Chike",
				DOB:        time.Now(),
				Department: "Engineering",
				Password:   "test",
			},
			name: "invalid email",
			err:  errors.New("email should be a valid email"),
		},
		{
			data: v1resources.UserRegistration{
				Email:      "michael@gmail.com",
				FirstName:  "Michael",
				LastName:   "Chike",
				DOB:        time.Now(),
				Department: "Engineering",
				Password:   "test",
			},
			name: "valid email",
		},
	} {
		err := ValidateInputs(val.data)
		assert.Equal(t, val.err, err, val.name)
	}
}
