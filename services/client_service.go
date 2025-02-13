package services

import (
	"oauth2-provider/models"
	"oauth2-provider/storage"
)

type ClientService struct {
	store *storage.PostgresStorage
}

func NewClientService(store *storage.PostgresStorage) *ClientService {
	return &ClientService{store: store}
}

func (s *ClientService) RegisterClient(req *models.ClientRegistration) (*models.Client, error) {
	client := &models.Client{
		RedirectURIs: req.RedirectURIs,
		GrantTypes:   []string{"authorization_code"},
	}

	err := s.store.StoreClient(client)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (s *ClientService) GetClient(clientID string) *models.Client {
	return s.store.GetClient(clientID)
}