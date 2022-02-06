package main

import (
	"sync"
	"time"
)

type aclStorageItem struct {
	rule   string
	name   string
	expire int64
}

type aclStorage struct {
	mu    sync.Mutex
	count int
	items map[string]aclStorageItem
}

func newAclStorage() *aclStorage {
	items := make(map[string]aclStorageItem)

	a := aclStorage{
		count: 0,
		items: items,
	}

	return &a
}

func (a *aclStorage) add(key string, rule string, name string, ttl int64) {
	a.mu.Lock()
	defer a.mu.Unlock()

	i := aclStorageItem{
		rule:   rule,
		name:   name,
		expire: time.Now().Unix() + ttl,
	}

	a.items[key] = i
	a.count += 1
}

func (a *aclStorage) exist(key string) *aclStorageItem {
	a.mu.Lock()
	defer a.mu.Unlock()

	storageItem, foundExpire := a.items[key]

	if foundExpire && storageItem.expire >= time.Now().Unix() {
		return &storageItem
	}

	return nil
}

func (a *aclStorage) gc() int {
	a.mu.Lock()
	defer a.mu.Unlock()

	for key, registerTime := range a.items {
		if registerTime.expire < time.Now().Unix() {
			delete(a.items, key)
		}
	}
	a.count = len(a.items)
	return a.count
}
