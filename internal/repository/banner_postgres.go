package repository

import (
	"banner-system/internal/models"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"strings"
)

type BannerPostgres struct {
	db *sqlx.DB
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
	query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM %s WHERE $1 = feature_id AND $2 && tag_ids)", bannersTable)
	exists := false
	err := b.db.QueryRow(query, input.FeatureId, input.TagIds).Scan(&exists)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, errors.New("Баннер с такой фичей и тегом уже существует")
	}

	var id int
	query = fmt.Sprintf("INSERT INTO %s (title, text, url, tag_ids, feature_id, is_active, created_at, updated_at) VALUES"+
		"($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id", bannersTable)
	row := b.db.QueryRow(query, input.Content.Title, input.Content.Text, input.Content.URL, input.TagIds, input.FeatureId, input.IsActive, input.CreatedAt, input.UpdatedAt)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (b *BannerPostgres) UpdateBanner(input models.BannerInfo) error {
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

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d",
		bannersTable, strings.Join(updates, ", "), len(args))
	_, err := b.db.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (b *BannerPostgres) DeleteBanner(id int) (models.CachKey, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 RETURNING feature_id, tag_ids", bannersTable)
	row := b.db.QueryRow(query, id)

	var deletedKey models.CachKey
	deletedKey.TagIds = &models.IntArray{}
	var tags pq.Int64Array
	err := row.Scan(&deletedKey.FeatureId, &tags)
	if err != nil {
		return models.CachKey{}, err
	}
	for _, i := range tags {
		*deletedKey.TagIds = append(*deletedKey.TagIds, int(i))
	}
	return deletedKey, nil
}

func NewBannerPostgres(db *sqlx.DB) *BannerPostgres {
	return &BannerPostgres{db: db}
}
