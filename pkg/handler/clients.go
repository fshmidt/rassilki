package handler

import (
	"github.com/fshmidt/rassilki"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary Get a client by ID
// @Tags clients
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {object} rassilki.Client
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /clients/{id} [get]
func (h *Handler) getClient(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	client, err := h.services.Client.Get(id)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, client)
}

// @Summary Create a new client
// @Tags clients
// @Accept json
// @Produce json
// @Param input body rassilki.Client true "Client object to create"
// @Success 200 {integer} integer 1
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /clients [post]
func (h *Handler) createClient(c *gin.Context) {
	var input rassilki.Client

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Client.Create(input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary Delete a client by ID
// @Tags clients
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /clients/{id} [delete]
func (h *Handler) deleteClient(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	err = h.services.Client.Delete(id)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

// @Tags clients
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param input body rassilki.UpdateClient true "Updated fields for the client"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /clients/{id} [put]
func (h *Handler) updateClient(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input rassilki.UpdateClient
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = input.Validate(); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Client.Update(input, id)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
