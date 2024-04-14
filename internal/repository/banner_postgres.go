package repository

import (
	"banner-system/internal/models"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"strings"
	"time"
)

type BannerPostgres struct {
	db *sqlx.DB
}

func NewBannerPostgres(db *sqlx.DB) *BannerPostgres {
	return &BannerPostgres{db: db}
}

func (b *BannerPostgres) GetUserBanner(tagId, featId int, accessAdmin bool) (models.BannerContent, error) {
	var banner models.BannerContent
	var isActive bool
	query := fmt.Sprintf("SELECT title, text, url, is_active FROM %s WHERE $1 = feature_id AND $2 = ANY(tag_ids)", bannersTable)
	row := b.db.QueryRow(query, featId, tagId)
	if err := row.Scan(&banner.Title, &banner.Text, &banner.URL, &isActive); err != nil {
		return banner, err
	}
	if !accessAdmin && !isActive {
		return banner, errors.New("Пользователь не имеет доступа")
	}
	return banner, nil
}

func (b *BannerPostgres) GetAllBanners(filter models.Filter) ([]models.BannerInfo, error) {
	var query strings.Builder
	var params []interface{}

	query.WriteString(fmt.Sprintf("SELECT * FROM %s ", bannersTable))

	if filter.TagId != 0 || filter.FeatureId != 0 {
		query.WriteString("WHERE ")
	}

	if filter.TagId != 0 {
		query.WriteString(fmt.Sprintf("$%d = ANY(tag_ids) ", len(params)+1))
		params = append(params, filter.TagId)
	}

	if filter.TagId != 0 && filter.FeatureId != 0 {
		query.WriteString("AND ")
	}

	if filter.FeatureId != 0 {
		query.WriteString(fmt.Sprintf("$%d = feature_id ", len(params)+1))
		params = append(params, filter.FeatureId)
	}

	if filter.Limit != 0 {
		query.WriteString(fmt.Sprintf("LIMIT $%d ", len(params)+1))
		params = append(params, filter.Limit)
	}

	if filter.Offset != 0 {
		query.WriteString(fmt.Sprintf("OFFSET $%d ", len(params)+1))
		params = append(params, filter.Offset)
	}

	rows, err := b.db.Query(query.String(), params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var banners []models.BannerInfo
	for rows.Next() {
		var banner models.BannerInfo
		banner.Content = &models.BannerContent{}
		banner.TagIds = &models.IntArray{}
		var tags pq.Int64Array
		if err := rows.Scan(&banner.Id, &banner.Content.Title, &banner.Content.Text, &banner.Content.URL, &tags,
			&banner.FeatureId, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt); err != nil {
			return nil, err
		}
		for _, i := range tags {
			*banner.TagIds = append(*banner.TagIds, int(i))
		}
		banners = append(banners, banner)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return banners, nil
}

func (b *BannerPostgres) CreateBanner(input models.BannerInfo) (int, error) {
	exist, err := b.existBanner(*input.FeatureId, *input.TagIds)
	if err != nil {
		return 0, err
	}
	if exist {
		return 0, errors.New("Баннер с такой фичей и тегом уже существует")
	}

	var id int
	query := fmt.Sprintf("INSERT INTO %s (title, text, url, tag_ids, feature_id, is_active, created_at, updated_at) VALUES"+
		"($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id", bannersTable)
	row := b.db.QueryRow(query, input.Content.Title, input.Content.Text, input.Content.URL, input.TagIds, input.FeatureId, input.IsActive, input.CreatedAt, input.UpdatedAt)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (b *BannerPostgres) UpdateBanner(input models.BannerInfo) error {
	if input.FeatureId != nil && input.TagIds != nil {
		exist, err := b.existBanner(*input.FeatureId, *input.TagIds)
		if err != nil {
			return err
		}
		if exist {
			return errors.New("Баннер с такой фичей и тегом уже существует")
		}
	}
	updates := make([]string, 0)
	args := make([]interface{}, 0)
	if input.Content != nil {
		if input.Content.Title != nil {
			updates = append(updates, fmt.Sprintf("title = $%d", len(args)+1))
			args = append(args, *input.Content.Title)
		}
		if input.Content.Text != nil {
			updates = append(updates, fmt.Sprintf("text = $%d", len(args)+1))
			args = append(args, *input.Content.Text)
		}
		if input.Content.URL != nil {
			updates = append(updates, fmt.Sprintf("url = $%d", len(args)+1))
			args = append(args, *input.Content.URL)
		}
	}
	if input.TagIds != nil {
		updates = append(updates, fmt.Sprintf("tag_ids = $%d", len(args)+1))
		args = append(args, *input.TagIds)
	}
	if input.FeatureId != nil {
		updates = append(updates, fmt.Sprintf("feature_id = $%d", len(args)+1))
		args = append(args, *input.FeatureId)
	}
	if input.IsActive != nil {
		updates = append(updates, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *input.IsActive)
	}
	if len(updates) == 0 {
		return nil
	}
	updates = append(updates, fmt.Sprintf("updated_at = $%d", len(args)+1))
	args = append(args, input.UpdatedAt)
	args = append(args, input.Id)

	tx, err := b.db.Begin()
	if err != nil {
		return err
	}
	query := fmt.Sprintf("INSERT INTO %s (banner_id, title, text, url, tag_ids, feature_id, is_active, created_at, updated_at)"+
		"SELECT id, title, text, url, tag_ids, feature_id, is_active, created_at, updated_at "+
		"FROM %s WHERE id = $1", historyBannersTable, bannersTable)
	_, err = tx.Exec(query, input.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	query = fmt.Sprintf("DELETE FROM %s WHERE banner_id = $1 AND history_id NOT IN (SELECT history_id FROM %s "+
		"WHERE banner_id = $1 ORDER BY updated_at DESC LIMIT 3)", historyBannersTable, historyBannersTable)
	_, err = tx.Exec(query, input.Id)
	query = fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d",
		bannersTable, strings.Join(updates, ", "), len(args))
	_, err = tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (b *BannerPostgres) DeleteBanner(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", bannersTable)
	_, err := b.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (b *BannerPostgres) DeleteBannerByTagFeat(tagId int, featId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE ", bannersTable)
	var arg int
	if tagId != 0 {
		query = query + "$1 = ANY(tag_ids)"
		arg = tagId
	} else {
		query = query + "$1 = feature_id"
		arg = featId
	}
	_, err := b.db.Exec(query, arg)
	if err != nil {
		return err
	}
	return nil
}

func (b *BannerPostgres) GetOldVersionBanner(id int) ([]models.BannerInfo, error) {
	query := fmt.Sprintf("SELECT banner_id, title,  text, url, tag_ids, feature_id, is_active, created_at, updated_at "+
		"FROM %s WHERE banner_id = $1", historyBannersTable)
	rows, err := b.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var banners []models.BannerInfo
	for rows.Next() {
		var banner models.BannerInfo
		banner.Content = &models.BannerContent{}
		banner.TagIds = &models.IntArray{}
		var tags pq.Int64Array
		if err := rows.Scan(&banner.Id, &banner.Content.Title, &banner.Content.Text, &banner.Content.URL, &tags,
			&banner.FeatureId, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt); err != nil {
			return nil, err
		}
		for _, i := range tags {
			*banner.TagIds = append(*banner.TagIds, int(i))
		}
		banners = append(banners, banner)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return banners, nil
}

func (b *BannerPostgres) BackToVersion(id, version int) error {
	query := fmt.Sprintf("WITH banner_history AS ("+
		"SELECT * FROM %s WHERE banner_id = $1 LIMIT 1 OFFSET $2)\n"+
		"UPDATE %s SET title = banner_history.title, text = banner_history.text, url = banner_history.url, tag_ids = banner_history.tag_ids, "+
		"feature_id = banner_history.feature_id, is_active = banner_history.is_active, created_at = banner_history.created_at, "+
		"updated_at = $3 FROM banner_history "+
		"WHERE banners.id = banner_history.banner_id", historyBannersTable, bannersTable)
	_, err := b.db.Exec(query, id, version-1, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (b *BannerPostgres) existBanner(featureId int, tagIds models.IntArray) (bool, error) {
	query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM %s WHERE $1 = feature_id AND $2 && tag_ids)", bannersTable)
	exists := false
	err := b.db.QueryRow(query, featureId, tagIds).Scan(&exists)
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil
	}
	return false, nil
}
