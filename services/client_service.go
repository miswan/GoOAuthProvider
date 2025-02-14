package services

import (
	"github.com/lib/pq"
	"oauth2-provider/models"
	"oauth2-provider/storage"
	"log"
)

type ClientService struct {
	store *storage.PostgresStorage
}

func NewClientService(store *storage.PostgresStorage) *ClientService {
	return &ClientService{store: store}
}

func (s *ClientService) RegisterClient(req *models.ClientRegistration) (*models.Client, error) {
	// Log the incoming request
	log.Printf("Registering new client with RedirectURIs: %v", req.RedirectURIs)

	// Convert []string to pq.StringArray explicitly
	redirectURIs := make(pq.StringArray, len(req.RedirectURIs))
	copy(redirectURIs, req.RedirectURIs)

	client := &models.Client{
		RedirectURIs: redirectURIs,
		GrantTypes:   pq.StringArray{"authorization_code"},
	}

	// Log the client data before storing
	log.Printf("Client data before storing: RedirectURIs=%v, GrantTypes=%v", client.RedirectURIs, client.GrantTypes)

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