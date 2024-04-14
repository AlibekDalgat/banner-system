package delivery

import (
	"banner-system/internal/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) GetUserBanner(c *gin.Context) {
	accessAdminStr, _ := c.Get("isAdmin")
	accessAdmin, ok := accessAdminStr.(bool)
	if !ok {
		c.JSON(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	tagIdStr := c.Query("tag_id")
	tagId, err := strconv.Atoi(tagIdStr)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	featureIdStr := c.Query("feature_id")
	featureId, err := strconv.Atoi(featureIdStr)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	useLastRevision := false
	useLastRevisionStr := c.Query("use_last_revision")
	if useLastRevisionStr != "" {
		useLastRevision, err = strconv.ParseBool(useLastRevisionStr)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	var banner models.BannerContent
	val, err := h.cache.Get(c, featureIdStr+":"+tagIdStr).Result()
	if !useLastRevision && err != redis.Nil {
		err = json.Unmarshal([]byte(val), &banner)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		banner, err = h.services.GetUserBanner(tagId, featureId, accessAdmin)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		bannerJSON, err := json.Marshal(banner)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		h.cache.Set(c, featureIdStr+":"+tagIdStr, bannerJSON, 5*time.Minute)
	}
	c.JSON(http.StatusOK, banner)
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
	if *input.IsActive {
		bannerContentJSON, err := json.Marshal(input.Content)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		for _, tagId := range *input.TagIds {
			combineKey := strconv.Itoa(*input.FeatureId) + ":" + strconv.Itoa(tagId)
			h.cache.Set(c, combineKey, string(bannerContentJSON), 5*time.Minute)
		}
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
	err = h.services.UpdateBanner(input)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
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

	err = h.services.DeleteBanner(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusNoContent, map[string]interface{}{
		"description": "Баннер успешно удален",
	})
}

func (h *Handler) DeleteBannerByTagFeat(c *gin.Context) {
	checkAccess, _ := c.Get(isAdmin)
	if !checkAccess.(bool) {
		newErrorResponse(c, http.StatusForbidden, "Пользователь не имеет доступа")
		return
	}
	tagIdStr := c.Query("tag_id")
	var err error
	tagId := 0
	if tagIdStr != "" {
		tagId, err = strconv.Atoi(tagIdStr)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}
	featureIdStr := c.Query("feature_id")
	featureId := 0
	if featureIdStr != "" {
		featureId, err = strconv.Atoi(featureIdStr)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if featureId != 0 && tagId != 0 || featureId == 0 && tagId == 0 {
		newErrorResponse(c, http.StatusBadRequest, "Должен быть только один из параметров")
		return
	}

	err = h.services.DeleteBannerByTagFeat(tagId, featureId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusNoContent, map[string]interface{}{
		"description": "Баннеры успешно удалены",
	})
}

func (h *Handler) GetBannerVersions(c *gin.Context) {
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

	banners, err := h.services.GetOldVersionBanner(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, banners)
}

func (h *Handler) BackToVersion(c *gin.Context) {
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

	var version models.Version
	if err = c.BindJSON(&version); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.BackToVersion(id, version.Number)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"description": "ok",
	})
}
