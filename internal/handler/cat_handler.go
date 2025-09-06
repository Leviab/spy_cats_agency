package handler

import (
	"errors"
	"fmt"
	"net/http"
	"spy_cats_agency/internal/domain"
	"spy_cats_agency/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CatHandler handles the HTTP requests for cats.
type CatHandler struct {
	catService service.CatService
}

// NewCatHandler creates a new CatHandler.
func NewCatHandler(catService service.CatService) *CatHandler {
	return &CatHandler{catService: catService}
}

// CreateCat handles the creation of a new cat.
// @Summary Create a new spy cat
// @Description Adds a new spy cat to the system. Breed must be valid according to TheCatAPI.
// @Tags cats
// @Accept json
// @Produce json
// @Param cat body CreateCatRequest true "Cat to create"
// @Success 201 {object} domain.Cat
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /cats [post]
func (h *CatHandler) CreateCat(c *gin.Context) {
	var req CreateCatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	cat := &domain.Cat{
		Name:              req.Name,
		YearsOfExperience: req.YearsOfExperience,
		Breed:             req.Breed,
		Salary:            req.Salary,
	}

	if err := h.catService.CreateCat(c.Request.Context(), cat); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, cat)
}

// GetCat handles retrieving a single cat by its ID.
// @Summary Get a spy cat by ID
// @Description Retrieves details of a specific spy cat.
// @Tags cats
// @Produce json
// @Param id path int true "Cat ID"
// @Success 200 {object} domain.Cat
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /cats/{id} [get]
func (h *CatHandler) GetCat(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(errors.New("invalid id format"))
		return
	}

	cat, err := h.catService.GetCat(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, cat)
}

// ListCats handles listing all cats.
// @Summary List all spy cats
// @Description Retrieves a list of all spy cats in the system.
// @Tags cats
// @Produce json
// @Success 200 {array} domain.Cat
// @Failure 500 {object} ErrorResponse
// @Router /cats [get]
func (h *CatHandler) ListCats(c *gin.Context) {
	cats, err := h.catService.ListCats(c.Request.Context())
	if err != nil {
		fmt.Printf("err type: %T, value: %+v\n", err, err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, cats)
}

// UpdateCatSalary handles updating a cat's salary.
// @Summary Update a spy cat's salary
// @Description Updates the salary of a specific spy cat.
// @Tags cats
// @Accept json
// @Produce json
// @Param id path int true "Cat ID"
// @Param salary body UpdateCatSalaryRequest true "New salary"
// @Success 200 {object} domain.Cat
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /cats/{id}/salary [patch]
func (h *CatHandler) UpdateCatSalary(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(errors.New("invalid id format"))
		return
	}

	var req UpdateCatSalaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	cat, err := h.catService.UpdateCatSalary(c.Request.Context(), id, req.Salary)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, cat)
}

// DeleteCat handles deleting a cat.
// @Summary Delete a spy cat
// @Description Removes a spy cat from the system.
// @Tags cats
// @Param id path int true "Cat ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /cats/{id} [delete]
func (h *CatHandler) DeleteCat(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(errors.New("invalid id format"))
		return
	}

	if err := h.catService.DeleteCat(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
