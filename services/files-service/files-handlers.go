package filesservice

import (
	"panda/apigateway/helpers"
	"panda/apigateway/services/files-service/models"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type FilesHandlers struct {
	filesService IFilesService
}

type IFilesHandlers interface {
	GetFileLinksByParentUid() echo.HandlerFunc
	CreateFileLink() echo.HandlerFunc
	UpdateFileLink() echo.HandlerFunc
	DeleteFileLink() echo.HandlerFunc
	SetMiniImageUrlToNode() echo.HandlerFunc
}

// NewFilesHandlers Files handlers constructor
func NewFilesHandlers(filesSvc IFilesService) IFilesHandlers {
	return &FilesHandlers{filesService: filesSvc}
}

// GetFileLinksByParentUid godoc
// @Summary Get file links by parent UID
// @Description Returns file links associated with the given parent UID.
// @Tags Files
// @Produce json
// @Security BearerAuth
// @Param parentUid path string true "Parent UID"
// @Success 200 {array} models.FileLink
// @Failure 500 "Internal server error"
// @Router /v1/files/links/{parentUid} [get]
func (h *FilesHandlers) GetFileLinksByParentUid() echo.HandlerFunc {

	return func(c echo.Context) error {

		parentUid := c.Param("parentUid")
		links, err := h.filesService.GetFileLinksByParentUid(parentUid)
		if err != nil {
			log.Error().Err(err).Msg("Error getting file links by parent UID")
			return echo.ErrInternalServerError
		}

		return c.JSON(200, links)
	}
}

// CreateFileLink godoc
// @Summary Create file link
// @Description Creates a new file link under the given parent UID.
// @Tags Files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param parentUid path string true "Parent UID"
// @Param body body models.FileLink true "File link"
// @Success 201 {object} models.FileLink
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /v1/files/link/{parentUid} [post]
func (h *FilesHandlers) CreateFileLink() echo.HandlerFunc {

	return func(c echo.Context) error {

		parentUid := c.Param("parentUid")
		fileLink := models.FileLink{}
		if err := c.Bind(&fileLink); err != nil {
			log.Error().Err(err).Msg("Error binding file link")
			return helpers.BadRequest(err.Error())
		}

		result, err := h.filesService.CreateFileLink(parentUid, fileLink)

		if err != nil {
			log.Error().Err(err).Msg("Error creating file link")
			return echo.ErrInternalServerError
		}

		return c.JSON(201, result)
	}
}

// UpdateFileLink godoc
// @Summary Update file link
// @Description Updates an existing file link.
// @Tags Files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uid path string true "File link UID"
// @Param body body models.FileLink true "File link"
// @Success 200 {object} models.FileLink
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /v1/files/link/{uid} [put]
func (h *FilesHandlers) UpdateFileLink() echo.HandlerFunc {

	return func(c echo.Context) error {

		uid := c.Param("uid")
		fileLink := models.FileLink{}
		if err := c.Bind(&fileLink); err != nil {
			log.Error().Err(err).Msg("Error binding file link")
			return helpers.BadRequest(err.Error())
		}

		fileLink.UID = uid
		result, err := h.filesService.UpdateFileLink(fileLink)

		if err != nil {
			log.Error().Err(err).Msg("Error updating file link")
			return echo.ErrInternalServerError
		}

		return c.JSON(200, result)
	}
}

// DeleteFileLink godoc
// @Summary Delete file link
// @Description Deletes a file link by its UID.
// @Tags Files
// @Security BearerAuth
// @Param uid path string true "File link UID"
// @Success 204 "No content"
// @Failure 500 "Internal server error"
// @Router /v1/files/link/{uid} [delete]
func (h *FilesHandlers) DeleteFileLink() echo.HandlerFunc {

	return func(c echo.Context) error {

		uid := c.Param("uid")
		err := h.filesService.DeleteFileLink(uid)
		if err != nil {
			log.Error().Err(err).Msg("Error deleting file link")
			return echo.ErrInternalServerError
		}

		return c.NoContent(204)
	}
}

// SetMiniImageUrlToNode godoc
// @Summary Set mini image URL to node
// @Description Stores a mini image URL list on a node identified by UID.
// @Tags Files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uid path string true "Node UID"
// @Param nodeLabel query string true "Neo4j node label"
// @Param body body models.MiniImageLinks true "Mini image link payload"
// @Success 200 {object} models.MiniImageLinks
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /v1/files/node/{uid}/mini-image-url [post]
func (h *FilesHandlers) SetMiniImageUrlToNode() echo.HandlerFunc {

	return func(c echo.Context) error {

		uid := c.Param("uid")
		nodeLabel := c.QueryParam("nodeLabel")

		if nodeLabel == "" {
			msg := "Node label is required"
			log.Error().Msg(msg)
			return helpers.BadRequest(msg)
		}

		link := models.MiniImageLinks{}
		if err := c.Bind(&link); err != nil {
			msg := "Error binding file link"
			log.Error().Err(err).Msg(msg)
			return helpers.BadRequest(msg)
		}

		link.UID = uid

		err := h.filesService.SetMiniImageUrlToNode(link.UID, link.Url, nodeLabel)

		if err != nil {
			log.Error().Err(err).Msg("Error setting mini image URL to node")
			return echo.ErrInternalServerError
		}

		return c.JSON(200, link)
	}
}
