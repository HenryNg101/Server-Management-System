package agent

import (
	"context"
	"fmt"
	"strings"

	"github.com/HenryNg101/agent-service/internal/model"
	"github.com/HenryNg101/agent-service/internal/server"
	"github.com/HenryNg101/agent-service/internal/shared/auth"
)

type Service interface {
	AgentExist(ctx context.Context, key string, server *model.Agent) error
	PushMetrics(ctx context.Context, messages []MetricMessage) error
	RegisterAgent(ctx context.Context, userID uint, req RegisterAgentRequest) (*RegisterAgentResponse, error)
	RotateAPIKey(ctx context.Context, apiKey string) (string, error)
}

type agentService struct {
	repo          Repository
	kafkaProducer KafkaProducer
	serverRepo    server.Repository
}

func NewService(r Repository, k KafkaProducer, serverRepo server.Repository) Service {
	return &agentService{repo: r, kafkaProducer: k, serverRepo: serverRepo}
}

func (s *agentService) PushMetrics(ctx context.Context, messages []MetricMessage) error {
	return s.kafkaProducer.PublishMetrics(ctx, messages)
}

func (s *agentService) AgentExist(ctx context.Context, key string, server *model.Agent) error {
	hashedKey := auth.HashAPIKey(key)
	_, err := s.repo.FindByKey(ctx, hashedKey)

	if err != nil {
		return err
	}
	return nil
}

func (s *agentService) RegisterAgent(ctx context.Context, userID uint, req RegisterAgentRequest) (*RegisterAgentResponse, error) {
	// 1. Find if server exist for the user
	server, err := s.serverRepo.FindByNameAndUser(ctx, strings.TrimSpace(req.ServerName))
	if err != nil || server.UserID != userID {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("There's no server %s exist for this user", req.ServerName)
	}

	// 2. Check duplicate instance
	// TODO: Right now, this only prevents registration of duplicated same agents,
	existing, _ := s.repo.FindByInstance(ctx, server.ID, req.InstanceID)
	if existing != nil {
		return nil, fmt.Errorf("Agent with instance ID of %s is already registered for this server", req.InstanceID)
	}

	// 3. Generate API key
	rawKey, hashedKey, err := auth.GenerateAPIKey()
	if err != nil {
		return nil, err
	}

	// 4. Create agent
	agent := &model.Agent{
		ServerID:   server.ID,
		APIKey:     hashedKey,
		InstanceID: req.InstanceID,
		Hostname:   req.Hostname,
		Status:     "active",
	}

	agent, err = s.repo.Create(ctx, agent)
	if err != nil {
		return nil, err
	}

	return &RegisterAgentResponse{
		ServerID: server.ID,
		AgentID:  agent.ID,
		APIKey:   rawKey,
	}, nil
}

func (s *agentService) RotateAPIKey(ctx context.Context, apiKey string) (string, error) {
	// 1. find agent by API key
	hashedKey := auth.HashAPIKey(apiKey)
	agent, err := s.repo.FindByKey(ctx, apiKey)
	if err != nil {
		return "", err
	}

	// 2. generate new key
	rawKey, hashedKey, err := auth.GenerateAPIKey()
	if err != nil {
		return "", err
	}

	// 3. update DB
	err = s.repo.UpdateAPIKey(ctx, agent.ID, hashedKey)
	if err != nil {
		return "", err
	}

	return rawKey, nil
}
