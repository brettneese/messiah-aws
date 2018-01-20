package main

import (
	"hbkauth-api/packages/config"

	"github.com/brettneese/messiah-aws"
)

// StatusHandler is the handler for health checks.
type StatusHandler struct {
	Status config.StatusInfo
}

func (handler StatusHandler) Handle(req Messiah.Request) interface{} {
	status := handler.Status

	body := map[string]interface{}{
		"status":  status,
		"request": req,
	}

	res := Messiah.Response{
		StatusCode:   200,
		ResponseData: body,
	}

	return res
}

func main() {
	status := config.GetStatus()

	handler := StatusHandler{
		Status: status,
	}

	Messiah.Start(handler)
}
