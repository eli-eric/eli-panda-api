package filesservice

import (
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

func (h *FilesHandlers) CreateFileLink() echo.HandlerFunc {

	return func(c echo.Context) error {

		parentUid := c.Param("parentUid")
		fileLink := models.FileLink{}
		if err := c.Bind(&fileLink); err != nil {
			log.Error().Err(err).Msg("Error binding file link")
			return echo.ErrBadRequest
		}

		result, err := h.filesService.CreateFileLink(parentUid, fileLink)

		if err != nil {
			log.Error().Err(err).Msg("Error creating file link")
			return echo.ErrInternalServerError
		}

		return c.JSON(201, result)
	}
}

func (h *FilesHandlers) UpdateFileLink() echo.HandlerFunc {

	return func(c echo.Context) error {

		uid := c.Param("uid")
		fileLink := models.FileLink{}
		if err := c.Bind(&fileLink); err != nil {
			log.Error().Err(err).Msg("Error binding file link")
			return echo.ErrBadRequest
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

func (h *FilesHandlers) SetMiniImageUrlToNode() echo.HandlerFunc {

	return func(c echo.Context) error {

		uid := c.Param("uid")
		nodeLabel := c.QueryParam("nodeLabel")

		if nodeLabel == "" {
			log.Error().Msg("Node label is required")
			return echo.ErrBadRequest
		}

		link := models.MiniImageLinks{}
		if err := c.Bind(&link); err != nil {
			log.Error().Err(err).Msg("Error binding file link")
			return echo.ErrBadRequest
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
