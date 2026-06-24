package repository

import (
	"errors"
	"time"

	"github.com/unsia-erp/unsia-core-service/internal/domain"
	"gorm.io/gorm"
)

type OAuthRepository struct {
	db *gorm.DB
}

func NewOAuthRepository(db *gorm.DB) *OAuthRepository {
	return &OAuthRepository{db: db}
}

func (r *OAuthRepository) GetClientByID(clientID string) (*domain.OAuthClient, error) {
	var client domain.OAuthClient
	err := r.db.Where("client_id = ? AND is_active = true", clientID).First(&client).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *OAuthRepository) GetClientByUUID(uuid string) (*domain.OAuthClient, error) {
	var client domain.OAuthClient
	err := r.db.Where("id = ? AND is_active = true", uuid).First(&client).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *OAuthRepository) ValidateRedirectURI(oauthClientID string, uri string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.RedirectURI{}).
		Where("oauth_client_id = ? AND redirect_uri = ? AND is_active = true", oauthClientID, uri).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *OAuthRepository) SaveAuthCode(code *domain.OAuthAuthorizationCode) error {
	return r.db.Create(code).Error
}

func (r *OAuthRepository) GetAuthCode(codeHash string) (*domain.OAuthAuthorizationCode, error) {
	var code domain.OAuthAuthorizationCode
	err := r.db.Where("code_hash = ?", codeHash).First(&code).Error
	if err != nil {
		return nil, err
	}
	return &code, nil
}

func (r *OAuthRepository) MarkAuthCodeUsed(id string) error {
	return r.db.Model(&domain.OAuthAuthorizationCode{}).
		Where("id = ?", id).
		Update("used_at", time.Now()).Error
}

func (r *OAuthRepository) SaveRegistrationRequest(req *domain.ClientRegistrationRequest) error {
	return r.db.Create(req).Error
}

func (r *OAuthRepository) GetRegistrationRequest(id string) (*domain.ClientRegistrationRequest, error) {
	var req domain.ClientRegistrationRequest
	err := r.db.Where("id = ?", id).First(&req).Error
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *OAuthRepository) GetAllRegistrationRequests() ([]domain.ClientRegistrationRequest, error) {
	var reqs []domain.ClientRegistrationRequest
	err := r.db.Order("created_at desc").Find(&reqs).Error
	return reqs, err
}

func (r *OAuthRepository) CreateClient(client *domain.OAuthClient, uris []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(client).Error; err != nil {
			return err
		}
		for _, u := range uris {
			redirect := domain.RedirectURI{
				OAuthClientID: client.ID,
				RedirectURI:   u,
				IsActive:      true,
			}
			if err := tx.Create(&redirect).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *OAuthRepository) UpdateClientStatus(id string, status string) error {
	var update map[string]interface{}
	now := time.Now()
	switch status {
	case "ACTIVE":
		update = map[string]interface{}{"status": status, "approved_at": &now}
	case "SUSPENDED":
		update = map[string]interface{}{"status": status, "suspended_at": &now}
	case "REVOKED":
		update = map[string]interface{}{"status": status, "revoked_at": &now}
	default:
		return errors.New("invalid status")
	}
	return r.db.Model(&domain.OAuthClient{}).Where("id = ?", id).Updates(update).Error
}

func (r *OAuthRepository) GetClientRedirectURIs(clientID string) ([]string, error) {
	var uris []string
	err := r.db.Table("redirect_uris").
		Select("redirect_uri").
		Where("oauth_client_id = ? AND is_active = true", clientID).
		Pluck("redirect_uri", &uris).Error
	return uris, err
}
