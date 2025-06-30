package game

import "jhgambling/protocol"

type GameManager struct {
	GameProviders []protocol.GameProvider
}

func NewGameManager() *GameManager {
	return &GameManager{
		GameProviders: []protocol.GameProvider{},
	}
}

// RegisterProvider adds a new game provider to the manager.
func (gm *GameManager) RegisterProvider(provider protocol.GameProvider) {
	gm.GameProviders = append(gm.GameProviders, provider)
}

// GetProviderByID retrieves a game provider by its unique ID.
func (gm *GameManager) GetProviderByID(id string) protocol.GameProvider {
	for _, provider := range gm.GameProviders {
		if provider.GetID() == id {
			return provider
		}
	}
	return nil
}

// GetAllProviders returns all registered game providers.
func (gm *GameManager) GetAllProviders() []protocol.GameProvider {
	return gm.GameProviders
}

// GetGameInstances retrieves all game instances across all providers.
func (gm *GameManager) GetGameInstances() []protocol.GameInstance {
	var instances []protocol.GameInstance
	for _, provider := range gm.GameProviders {
		instances = append(instances, provider.GetInstances()...)
	}
	return instances
}

// GetInstanceByID retrieves a specific game instance by its provider ID and instance ID.
func (gm *GameManager) GetInstanceByID(providerID, instanceID string) protocol.GameInstance {
	provider := gm.GetProviderByID(providerID)
	if provider == nil {
		return nil
	}
	for _, instance := range provider.GetInstances() {
		if instance.GetID() == instanceID {
			return instance
		}
	}
	return nil
}