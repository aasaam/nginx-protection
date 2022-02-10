package main

import (
	"sync"
	"time"
)

type challengeStorage struct {
	mu    sync.Mutex
	count int
	items map[string]int64
}

func newChallengeStorage() *challengeStorage {
	items := make(map[string]int64)

	cs := challengeStorage{
		items: items,
		count: 0,
	}

	return &cs
}

func (s *challengeStorage) exist(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	registerTime, foundExpire := s.items[key]

	if foundExpire && registerTime >= time.Now().Unix() {
		return true
	}

	return false
}

func (s *challengeStorage) set(key string, lifeTime int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[key] = time.Now().Unix() + lifeTime
	s.count++
}

func (s *challengeStorage) gc() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	for key, registerTime := range s.items {
		if registerTime < time.Now().Unix() {
			delete(s.items, key)
		}
	}
	s.count = len(s.items)
	return s.count
}
