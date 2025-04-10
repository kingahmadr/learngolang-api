package schema

type UserRequest struct {
	Name     string `json:"name" `
	Email    string `json:"email" `
	Password string `json:"password" `
}

type UserResponse struct {
	ID        uint    `json:"ID"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
	DeletedAt *string `json:"deletedAt,omitempty"`
}
