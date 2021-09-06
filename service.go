package api

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type Manager struct {
	echo *echo.Echo
}

type EndpointHandler struct {
	Type     string
	Endpoint string
	Handler  func(context echo.Context, user *User) error
}

func GET(route string, handler func(context echo.Context, user *User) error) EndpointHandler {
	return EndpointHandler{
		Type:     "GET",
		Endpoint: route,
		Handler:  handler,
	}
}
func POST(route string, handler func(context echo.Context, user *User) error) EndpointHandler {
	return EndpointHandler{
		Type:     "POST",
		Endpoint: route,
		Handler:  handler,
	}
}
func PUT(route string, handler func(context echo.Context, user *User) error) EndpointHandler {
	return EndpointHandler{
		Type:     "PUT",
		Endpoint: route,
		Handler:  handler,
	}
}
func DELETE(route string, handler func(context echo.Context, user *User) error) EndpointHandler {
	return EndpointHandler{
		Type:     "DELETE",
		Endpoint: route,
		Handler:  handler,
	}
}

func NewApi() *Manager {
	echoI := echo.New()
	echoI.HideBanner = true
	echoI.HidePort = true

	return &Manager{
		echo: echoI,
	}
}

func IsFromGateway(c echo.Context) bool {
	return c.Request().Header.Get("API_GATEWAY") != ""
}

func (m *Manager) handler(apiHandler EndpointHandler, requireAuth bool) {
	log.Debugf("Registering route %s %s ", apiHandler.Type, apiHandler.Endpoint)
	var fn func(string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route
	switch apiHandler.Type {
	case "GET":
		fn = m.echo.GET
	case "POST":
		fn = m.echo.POST
	case "PUT":
		fn = m.echo.PUT
	case "DELETE":
		fn = m.echo.DELETE
	default:
		log.Errorf("Invalid method: %s", apiHandler.Type)
		return
	}

	fn(apiHandler.Endpoint, func(context echo.Context) error {
		user := &User{}
		err := json.Unmarshal([]byte(context.Request().Header.Get("AuthUser")), user)
		if err != nil || user.Id == uuid.Nil {
			if requireAuth {
				return context.NoContent(204)
			}
			return apiHandler.Handler(context, nil)
		}
		return apiHandler.Handler(context, user)
	})
}
func (m *Manager) PublicHandler(apiHandler EndpointHandler) {
	m.handler(apiHandler, false)
}

func (m *Manager) PrivateHandler(apiHandler EndpointHandler) {
	m.handler(apiHandler, true)
}

func (m *Manager) Listen() error {
	return m.ListenOnPort(":80")
}
func (m *Manager) ListenOnPort(port string) error {
	return m.echo.Start(port)
}
