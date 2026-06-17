package agent

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"github.com/HenryNg101/server-management-system/internal/model"
)

type Service interface {
	CreateAgent(ctx context.Context, serverID uint) (*model.Agent, string, error)
	AgentExist(ctx context.Context, key string, server *model.Agent) error
	PushMetrics(ctx context.Context, msg MetricMessage) error
}

type agentService struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &agentService{repo: r}
}

// TODO: Push to Kafka
func (s *agentService) PushMetrics(ctx context.Context, msg MetricMessage) error {
	return nil
}

// TODO: Use this in server's creation
func (s *agentService) CreateAgent(ctx context.Context, serverID uint) (*model.Agent, string, error) {
	rawKey, hash, err := generateAPIKey()
	if err != nil {
		return nil, "", err
	}

	agent := model.Agent{
		ServerID:   serverID,
		APIKeyHash: hash,
	}
	createdAgent, err := s.repo.Create(ctx, &agent)

	return createdAgent, rawKey, err
}

func (s *agentService) AgentExist(ctx context.Context, key string, server *model.Agent) error {
	hashedKey := hashAPIKey(key)
	_, err := s.repo.FindByKey(ctx, hashedKey)

	if err != nil {
		return err
	}
	return nil
}

// Generate a secure random API key
func generateAPIKey() (string, string, error) {
	keyBytes := make([]byte, 32) // 256-bit
	_, err := rand.Read(keyBytes)
	if err != nil {
		return "", "", err
	}

	rawKey := hex.EncodeToString(keyBytes)

	hash := sha256.Sum256([]byte(rawKey))
	hashStr := hex.EncodeToString(hash[:])

	return rawKey, hashStr, nil
}

// Hash incoming key (for comparison)
func hashAPIKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}
