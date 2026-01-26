package services

import (
	"log"
	"oauth2-provider/models"
	"oauth2-provider/storage"
	"oauth2-provider/utils"
)

type ClientService struct {
	store storage.Storage
}

func NewClientService(store storage.Storage) *ClientService {
	return &ClientService{store: store}
}

func (s *ClientService) RegisterClient(req *models.ClientRegistration) (*models.Client, error) {
	// Log the incoming request
	log.Printf("Registering new client with RedirectURIs: %v", req.RedirectURIs)

	client := &models.Client{
		ClientID:     utils.GenerateRandomString(24),
		Secret:       utils.GenerateRandomString(32),
		RedirectURIs: req.RedirectURIs,
		GrantTypes:   []string{"authorization_code"},
	}

	// Log the client data before storing
	log.Printf("Client data before storing: ClientID=%s, RedirectURIs=%v", client.ClientID, client.RedirectURIs)

	err := s.store.StoreClient(client)
	if err != nil {
		log.Printf("Error storing client: %v", err)
		return nil, err
	}

	log.Printf("Successfully registered client with ID: %s", client.ClientID)
	return client, nil
}

func (s *ClientService) GetClient(clientID string) *models.Client {
	return s.store.GetClient(clientID)
}
