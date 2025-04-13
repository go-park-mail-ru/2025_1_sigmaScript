package dto

type UpdateUserRequest struct {
	Username            string `json:"username"`
	Avatar              string `json:"avatar,omitempty"`
	OldPassword         string `json:"old_password"`
	NewPassword         string `json:"new_password"`
	RepeatedNewPassword string `json:"repeated_new_password"`
}
