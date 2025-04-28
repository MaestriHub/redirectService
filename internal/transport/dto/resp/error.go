package resp

type ErrorDTO struct {
	Error string `json:"error" example:"Описание ошибки"`
}

func NewErrorDTO(s string) *ErrorDTO {
	return &ErrorDTO{Error: s}
}
