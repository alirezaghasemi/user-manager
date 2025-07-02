package response

type CreatedUserResponse struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	Family string `json:"family"`
	Email  string `json:"email"`
	Age    int    `json:"age"`
}

type UpdatedUserResponse struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	Family string `json:"family"`
	Email  string `json:"email"`
	Age    int    `json:"age"`
}

type DeletedUserResponse struct {
}
