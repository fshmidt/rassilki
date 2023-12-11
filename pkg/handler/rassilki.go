package handler

import (
	"github.com/fshmidt/rassilki"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

const (
	time_changes = iota
	recreate_messages_table
	add_clients_to_messages_table
	no_changes
	deleting_tags
)

// @Summary Get reviews for all rassilkas
// @Tags rassilki
// @Accept json
// @Produce json
// @Success 200 {array} rassilki.RassilkaReview
// @Failure 500 {object} errorResponse
// @Router /rassilki [get]
func (h *Handler) getRassilkiReview(c *gin.Context) {

	ids, err := h.services.Rassilka.GetAll()
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	review, err := h.services.Messages.GetRassilkiReview(ids)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, review)

}

// @Summary Get review for a specific rassilka by ID
// @Tags rassilki
// @Accept json
// @Produce json
// @Param id path int true "Rassilka ID"
// @Success 200 {object} rassilki.RassilkaReview
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /rassilki/{id} [get]
func (h *Handler) getRassilkaReviewById(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	review, err := h.services.Messages.GetRassilkaReviewById(id)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, review)

}

// @Summary Create a new rassilka
// @Tags rassilki
// @Accept json
// @Produce json
// @Param input body rassilki.Rassilka true "Rassilka object to create"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /rassilki [post]
func (h *Handler) createRassilka(c *gin.Context) {
	var input rassilki.Rassilka

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if input.EndTime.Before(time.Now()) {
		NewErrorResponse(c, http.StatusBadRequest, "EndTime is in the past")
		return
	}

	id, err := h.services.Rassilka.Create(input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary Delete a rassilka by ID
// @Tags rassilki
// @Accept json
// @Produce json
// @Param id path int true "Rassilka ID"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /rassilki/{id} [delete]
func (h *Handler) deleteRassilka(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	err = h.services.Rassilka.Delete(id)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

// @Summary Update a rassilka by ID
// @Tags rassilki
// @Accept json
// @Produce json
// @Param id path int true "Rassilka ID"
// @Param input body rassilki.UpdateRassilka true "UpdateRassilka object"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /rassilki/{id} [put]
func (h *Handler) updateRassilka(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	input := rassilki.UpdateRassilka{
		Supplemented: new(bool),
		Recreated:    new(bool),
	}
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = input.Validate(); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	original, err := h.services.Rassilka.GetById(id)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	updateLevel := findUpdateLevel(original, input)
	if updateLevel == no_changes {
		NewErrorResponse(c, http.StatusBadRequest, "no changes in updating struct")
		return
	} else if updateLevel == deleting_tags {
		NewErrorResponse(c, http.StatusBadRequest, "you can only add tags to created rassilka")
		return
	}
	err = h.services.Rassilka.Update(input, id)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	switch updateLevel {
	case recreate_messages_table:
		err = h.services.Messages.DropTable(id)
		if err != nil {
			NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		clientsIds, err2 := h.services.Messages.GetClientsList(id)
		if err2 != nil {
			NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		err = h.services.Messages.CreateTable(id, clientsIds)
		if err != nil {
			NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	case add_clients_to_messages_table:
		_, difference := notContain(*input.Filter, original.Filter)
		*input.Filter = difference
		err = h.services.Messages.RenewTable(input, id)
		if err != nil {
			NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	default:
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
