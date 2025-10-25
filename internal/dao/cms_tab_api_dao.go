package dao

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CMSTabAPI represents a mapping between CMS tab and API endpoint
type CMSTabAPI struct {
	ID          string
	TabName     string
	APIPath     string
	APIMethod   string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CMSTabAPIDAO handles database operations for CMS tab-API mappings
type CMSTabAPIDAO interface {
	Create(ctx context.Context, tabAPI *CMSTabAPI) error
	FindByID(ctx context.Context, id string) (*CMSTabAPI, error)
	FindByTab(ctx context.Context, tabName string) ([]*CMSTabAPI, error)
	FindByAPI(ctx context.Context, apiPath, method string) ([]*CMSTabAPI, error)
	ListAll(ctx context.Context) ([]*CMSTabAPI, error)
	Update(ctx context.Context, tabAPI *CMSTabAPI) error
	Delete(ctx context.Context, id string) error
	DeleteByTab(ctx context.Context, tabName string) error
	Exists(ctx context.Context, tabName, apiPath, method string) (bool, error)
}

type cmsTabAPIDAO struct {
	db *sql.DB
}

// NewCMSTabAPIDAO creates a new CMS tab-API DAO
func NewCMSTabAPIDAO(db *sql.DB) CMSTabAPIDAO {
	return &cmsTabAPIDAO{db: db}
}

func (d *cmsTabAPIDAO) Create(ctx context.Context, tabAPI *CMSTabAPI) error {
	if tabAPI.ID == "" {
		tabAPI.ID = uuid.New().String()
	}
	
	now := time.Now()
	tabAPI.CreatedAt = now
	tabAPI.UpdatedAt = now

	query := `
		INSERT INTO cms_tab_apis (id, tab_name, api_path, api_method, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := d.db.ExecContext(ctx, query,
		tabAPI.ID,
		tabAPI.TabName,
		tabAPI.APIPath,
		tabAPI.APIMethod,
		tabAPI.Description,
		tabAPI.CreatedAt,
		tabAPI.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create cms tab-api mapping: %w", err)
	}

	return nil
}

func (d *cmsTabAPIDAO) FindByID(ctx context.Context, id string) (*CMSTabAPI, error) {
	query := `
		SELECT id, tab_name, api_path, api_method, description, created_at, updated_at
		FROM cms_tab_apis
		WHERE id = $1
	`

	tabAPI := &CMSTabAPI{}
	err := d.db.QueryRowContext(ctx, query, id).Scan(
		&tabAPI.ID,
		&tabAPI.TabName,
		&tabAPI.APIPath,
		&tabAPI.APIMethod,
		&tabAPI.Description,
		&tabAPI.CreatedAt,
		&tabAPI.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("cms tab-api mapping not found with id: %s", id)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find cms tab-api mapping: %w", err)
	}

	return tabAPI, nil
}

func (d *cmsTabAPIDAO) FindByTab(ctx context.Context, tabName string) ([]*CMSTabAPI, error) {
	query := `
		SELECT id, tab_name, api_path, api_method, description, created_at, updated_at
		FROM cms_tab_apis
		WHERE tab_name = $1
		ORDER BY api_path, api_method
	`

	rows, err := d.db.QueryContext(ctx, query, tabName)
	if err != nil {
		return nil, fmt.Errorf("failed to find cms tab-apis by tab: %w", err)
	}
	defer rows.Close()

	var tabAPIs []*CMSTabAPI
	for rows.Next() {
		tabAPI := &CMSTabAPI{}
		err := rows.Scan(
			&tabAPI.ID,
			&tabAPI.TabName,
			&tabAPI.APIPath,
			&tabAPI.APIMethod,
			&tabAPI.Description,
			&tabAPI.CreatedAt,
			&tabAPI.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cms tab-api: %w", err)
		}
		tabAPIs = append(tabAPIs, tabAPI)
	}

	return tabAPIs, nil
}

func (d *cmsTabAPIDAO) FindByAPI(ctx context.Context, apiPath, method string) ([]*CMSTabAPI, error) {
	query := `
		SELECT id, tab_name, api_path, api_method, description, created_at, updated_at
		FROM cms_tab_apis
		WHERE api_path = $1 AND api_method = $2
		ORDER BY tab_name
	`

	rows, err := d.db.QueryContext(ctx, query, apiPath, method)
	if err != nil {
		return nil, fmt.Errorf("failed to find cms tab-apis by api: %w", err)
	}
	defer rows.Close()

	var tabAPIs []*CMSTabAPI
	for rows.Next() {
		tabAPI := &CMSTabAPI{}
		err := rows.Scan(
			&tabAPI.ID,
			&tabAPI.TabName,
			&tabAPI.APIPath,
			&tabAPI.APIMethod,
			&tabAPI.Description,
			&tabAPI.CreatedAt,
			&tabAPI.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cms tab-api: %w", err)
		}
		tabAPIs = append(tabAPIs, tabAPI)
	}

	return tabAPIs, nil
}

func (d *cmsTabAPIDAO) ListAll(ctx context.Context) ([]*CMSTabAPI, error) {
	query := `
		SELECT id, tab_name, api_path, api_method, description, created_at, updated_at
		FROM cms_tab_apis
		ORDER BY tab_name, api_path, api_method
	`

	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list cms tab-apis: %w", err)
	}
	defer rows.Close()

	var tabAPIs []*CMSTabAPI
	for rows.Next() {
		tabAPI := &CMSTabAPI{}
		err := rows.Scan(
			&tabAPI.ID,
			&tabAPI.TabName,
			&tabAPI.APIPath,
			&tabAPI.APIMethod,
			&tabAPI.Description,
			&tabAPI.CreatedAt,
			&tabAPI.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cms tab-api: %w", err)
		}
		tabAPIs = append(tabAPIs, tabAPI)
	}

	return tabAPIs, nil
}

func (d *cmsTabAPIDAO) Update(ctx context.Context, tabAPI *CMSTabAPI) error {
	tabAPI.UpdatedAt = time.Now()

	query := `
		UPDATE cms_tab_apis
		SET tab_name = $2, api_path = $3, api_method = $4, description = $5, updated_at = $6
		WHERE id = $1
	`

	result, err := d.db.ExecContext(ctx, query,
		tabAPI.ID,
		tabAPI.TabName,
		tabAPI.APIPath,
		tabAPI.APIMethod,
		tabAPI.Description,
		tabAPI.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update cms tab-api mapping: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("cms tab-api mapping not found with id: %s", tabAPI.ID)
	}

	return nil
}

func (d *cmsTabAPIDAO) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM cms_tab_apis WHERE id = $1`

	result, err := d.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete cms tab-api mapping: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("cms tab-api mapping not found with id: %s", id)
	}

	return nil
}

func (d *cmsTabAPIDAO) DeleteByTab(ctx context.Context, tabName string) error {
	query := `DELETE FROM cms_tab_apis WHERE tab_name = $1`

	_, err := d.db.ExecContext(ctx, query, tabName)
	if err != nil {
		return fmt.Errorf("failed to delete cms tab-apis by tab: %w", err)
	}

	return nil
}

func (d *cmsTabAPIDAO) Exists(ctx context.Context, tabName, apiPath, method string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM cms_tab_apis
			WHERE tab_name = $1 AND api_path = $2 AND api_method = $3
		)
	`

	var exists bool
	err := d.db.QueryRowContext(ctx, query, tabName, apiPath, method).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check cms tab-api existence: %w", err)
	}

	return exists, nil
}

