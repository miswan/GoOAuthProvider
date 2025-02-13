package handlers

import (
    "github.com/labstack/echo/v4"
    "net/http"
    "oauth2-provider/models"
    "oauth2-provider/services"
)

type ClientHandler struct {
    clientService *services.ClientService
}

func NewClientHandler(clientService *services.ClientService) *ClientHandler {
    return &ClientHandler{clientService: clientService}
}

func (h *ClientHandler) Register(c echo.Context) error {
    req := new(models.ClientRegistration)
    if err := c.Bind(req); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }

    client, err := h.clientService.RegisterClient(req)
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }

    return c.JSON(http.StatusCreated, map[string]interface{}{
        "client_id": client.ID,
        "client_secret": client.Secret,
        "redirect_uris": client.RedirectURIs,
    })
}

func (h *ClientHandler) Get(c echo.Context) error {
    clientID := c.Param("id")
    client := h.clientService.GetClient(clientID)
    if client == nil {
        return echo.NewHTTPError(http.StatusNotFound, "client not found")
    }

    return c.JSON(http.StatusOK, client)
}
