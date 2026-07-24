package agent

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/HenryNg101/agents-service/internal/model"
	"github.com/HenryNg101/agents-service/internal/server"
	"github.com/HenryNg101/agents-service/internal/shared/auth"
	"github.com/redis/go-redis/v9"
)

type Service interface {
	FindByHashedKey(ctx context.Context, hashedKey string) (*model.Agent, error)
	PushMetrics(ctx context.Context, messages []MetricMessage) error
	RegisterAgent(ctx context.Context, userID uint, req RegisterAgentRequest) (*RegisterAgentResponse, error)
	RotateAPIKey(ctx context.Context, apiKey string) (string, error)
}

type agentService struct {
	repo          Repository
	kafkaProducer KafkaProducer
	serverRepo    server.Repository
	redisClient   *redis.Client
}

func NewService(r Repository, k KafkaProducer, serverRepo server.Repository, redisClient *redis.Client) Service {
	return &agentService{repo: r, kafkaProducer: k, serverRepo: serverRepo, redisClient: redisClient}
}

func (s *agentService) PushMetrics(ctx context.Context, messages []MetricMessage) error {
	return s.kafkaProducer.PublishMetrics(ctx, messages)
}

func (s *agentService) FindByHashedKey(ctx context.Context, hashedKey string) (*model.Agent, error) {
	return s.repo.FindByKey(ctx, hashedKey)
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

	// Check duplicate instance
	// No need to do this, because there's already constraint in DB to enforce this
	// existing, _ := s.repo.FindByInstance(ctx, server.ID, req.InstanceID)
	// if existing != nil {
	// 	return nil, fmt.Errorf("Agent with instance ID of %s is already registered for this server", req.InstanceID)
	// }

	// 2. Generate API key
	rawKey, hashedKey, err := auth.GenerateAPIKey()
	if err != nil {
		return nil, err
	}

	// 3. Create agent
	agent := &model.Agent{
		ServerID:   server.ID,
		APIKey:     hashedKey,
		InstanceID: req.InstanceID,
		Status:     "active",
	}

	// Create or replace. 
	// The only risk is someone having email and password of the user, and override this. But if that happens, user could just reset account information
	agent, err = s.repo.Upsert(ctx, agent)
	if err != nil {
		return nil, err
	}

	// 4. Write-through to Redis
	redisKey := "agent:auth:" + hashedKey
	err = s.redisClient.Set(ctx, redisKey, server.ID, 0).Err() // no TTL
	if err != nil {
		// don't fail the request → system still consistent via DB fallback
		// but LOG it
		log.Println("[WARN] failed to write agent key to redis:", err)
	}

	return &RegisterAgentResponse{
		ServerID: server.ID,
		AgentID:  agent.ID,
		APIKey:   rawKey,
	}, nil
}

func (s *agentService) RotateAPIKey(ctx context.Context, apiKey string) (string, error) {
	// 1. hash incoming key
	oldHashedKey := auth.HashAPIKey(apiKey)

	// 2. find agent by API key
	agent, err := s.repo.FindByKey(ctx, oldHashedKey)
	if err != nil {
		return "", err
	}

	// 3. generate new key
	newRawKey, newlyHashedKey, err := auth.GenerateAPIKey()
	if err != nil {
		return "", err
	}

	// 4. update DB
	err = s.repo.UpdateAPIKey(ctx, agent.ID, newlyHashedKey)
	if err != nil {
		return "", err
	}

	// 5. update Redis (write-through sync)
	oldRedisKey := "agent:auth:" + oldHashedKey
	newRedisKey := "agent:auth:" + newlyHashedKey

	// delete old
	if err := s.redisClient.Del(ctx, oldRedisKey).Err(); err != nil {
		log.Println("[WARN] failed to delete old key from redis:", err)
	}

	// set new
	if err := s.redisClient.Set(ctx, newRedisKey, agent.ServerID, 0).Err(); err != nil {
		log.Println("[WARN] failed to set new key to redis:", err)
	}

	return newRawKey, nil
}
