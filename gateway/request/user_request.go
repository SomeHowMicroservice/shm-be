package request

import "time"

type UpdateProfileRequest struct {
	FirstName *string    `json:"first_name" binding:"omitempty"`
	LastName  *string    `json:"last_name" binding:"omitempty"`
	Gender    *string    `json:"gender" binding:"omitempty,oneof=male female"`
	DOB       *time.Time `json:"dob" binding:"omitempty"`
}
