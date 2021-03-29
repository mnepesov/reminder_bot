package cache

import (
	"errors"
	"reminder_bot/internal/models"
	"sync"
)

type userState struct {
	navigation string
	info       *models.User
}

type MemoryCache struct {
	cache map[int]*userState
	sync.RWMutex
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		cache: make(map[int]*userState),
	}
}
func (m *MemoryCache) GetUserInfo(userId int) (*models.User, error) {
	m.RLock()
	defer m.RUnlock()
	state, ok := m.cache[userId]
	if !ok {
		return &models.User{}, errors.New("not found")
	}
	return state.info, nil
}
func (m *MemoryCache) SetUserInfo(userId int, user *models.User) error {

	navigation, err := m.GetNavigation(userId)
	if err != nil {
		navigation = "start"
	}

	m.Lock()
	defer m.Unlock()
	m.cache[userId] = &userState{
		navigation: navigation,
		info:       user,
	}
	return nil
}

func (m *MemoryCache) SetNavigation(userId int, navigation string) error {
	m.Lock()
	defer m.Unlock()
	m.cache[userId] = &userState{
		navigation: navigation,
	}
	return nil
}
func (m *MemoryCache) GetNavigation(userId int) (string, error) {
	m.RLock()
	defer m.RUnlock()
	state, ok := m.cache[userId]
	if !ok {
		return "", errors.New("not found")
	}
	return state.navigation, nil
}
