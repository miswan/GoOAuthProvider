package services

import (
    "oauth2-provider/models"
    "oauth2-provider/storage"
    "oauth2-provider/utils"
)

type ClientService struct {
    store *storage.MemoryStorage
}

func NewClientService(store *storage.MemoryStorage) *ClientService {
    return &ClientService{store: store}
}

func (s *ClientService) RegisterClient(req *models.ClientRegistration) (*models.Client, error) {
    client := &models.Client{
        ID:           utils.GenerateRandomString(24),
        Secret:       utils.GenerateRandomString(32),
        RedirectURIs: req.RedirectURIs,
        GrantTypes:   []string{"authorization_code"},
    }

    s.store.StoreClient(client)
    return client, nil
}

func (s *ClientService) GetClient(clientID string) *models.Client {
    return s.store.GetClient(clientID)
}
