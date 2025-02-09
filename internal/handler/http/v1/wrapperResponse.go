package v1

type WrapperResponse struct {
	Message *string `json:"message" example:",Все пропало"`          // Сообщение об ошибке
	Data    any     `json:"data"`                                    // Запрашиваемые данные при успешном запросе
	Code    uint16  `json:"code"  binding:"required"  example:"200"` // Код запроса
}

func NewErrorWrapper(c uint16, err error) WrapperResponse {
	m := "Ошибка: " + err.Error()
	return WrapperResponse{
		Message: &m,
		Code:    c,
	}
}

func NewDataWrapper[T any](d T) WrapperResponse {
	return WrapperResponse{
		Data: d,
		Code: 200,
	}
}

func NewCDWrapper[T any](c uint16, d T) WrapperResponse {
	return WrapperResponse{
		Data: d,
		Code: c,
	}
}

func New200Wrapper() WrapperResponse {
	return WrapperResponse{
		Code: 200,
	}
}
