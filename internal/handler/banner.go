package delivery

import (
	"banner-system/internal/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) GetUserBanner(c *gin.Context) {
	fmt.Println("Get User Banner")
}

func (h *Handler) GetAllBanners(c *gin.Context) {
	checkAccess, _ := c.Get(isAdmin)
	if !checkAccess.(bool) {
		newErrorResponse(c, http.StatusForbidden, "Пользователь не имеет доступа")
		return
	}
	var (
		filter models.Filter
		err    error
	)
	tagIdStr := c.Query("tag_id")
	if tagIdStr != "" {
		filter.TagId, err = strconv.Atoi(tagIdStr)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}
	featureIdStr := c.Query("feature_id")
	if featureIdStr != "" {
		filter.FeatureId, err = strconv.Atoi(featureIdStr)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}
	limitStr := c.Query("limit")
	if limitStr != "" {
		filter.Limit, err = strconv.Atoi(limitStr)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}
	offsetStr := c.Query("offset")
	if offsetStr != "" {
		filter.Offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}
	banners, err := h.services.GetAllBanners(filter)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, banners)
}

func (h *Handler) CreateBanner(c *gin.Context) {
	checkAccess, _ := c.Get(isAdmin)
	if !checkAccess.(bool) {
		newErrorResponse(c, http.StatusForbidden, "Пользователь не имеет доступа")
		return
	}

	var input models.BannerInfo
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.CreateBanner(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	bannerContentJSON, err := json.Marshal(input.Content)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	for _, tagId := range *input.TagIds {
		combineKey := strconv.Itoa(*input.FeatureId) + ":" + strconv.Itoa(tagId)
		h.cache.Set(c, combineKey, string(bannerContentJSON), 5*time.Minute)
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"banner_id": id,
	})
}

func (h *Handler) UpdateBanner(c *gin.Context) {
	checkAccess, _ := c.Get(isAdmin)
	if !checkAccess.(bool) {
		newErrorResponse(c, http.StatusForbidden, "Пользователь не имеет доступа")
		return
	}
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Некорректные данные")
		return
	}

	var input models.BannerInfo
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	input.Id = id
	oldKey, banner, err := h.services.UpdateBanner(input)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	for _, tagId := range *oldKey.TagIds {
		combineKey := strconv.Itoa(*oldKey.FeatureId) + ":" + strconv.Itoa(tagId)
		h.cache.Del(c, combineKey)
	}

	bannerContentJSON, err := json.Marshal(banner.Content)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	for _, tagId := range *banner.TagIds {
		combineKey := strconv.Itoa(*banner.FeatureId) + ":" + strconv.Itoa(tagId)
		h.cache.Set(c, combineKey, string(bannerContentJSON), 5*time.Minute)
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"description": "ok",
	})
}

func (h *Handler) DeleteBanner(c *gin.Context) {
	checkAccess, _ := c.Get(isAdmin)
	if !checkAccess.(bool) {
		newErrorResponse(c, http.StatusForbidden, "Пользователь не имеет доступа")
		return
	}
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Некорректные данные")
		return
	}

	deletedKey, err := h.services.DeleteBanner(id)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	for _, tagId := range *deletedKey.TagIds {
		combineKey := strconv.Itoa(*deletedKey.FeatureId) + ":" + strconv.Itoa(tagId)
		h.cache.Del(c, combineKey)
	}
	c.JSON(http.StatusNoContent, map[string]interface{}{
		"description": "Баннер успешно удален",
	})
}
