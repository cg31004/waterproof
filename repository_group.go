package errortool

import (
	"log"
	"sync"
)

// newGroupRepository: registration error code, and check error code is unique.
func newGroupRepository() iGroupRepository {
	return &groupRepository{}
}

type iGroupRepository interface {
	Add(code int)
	Get(code int) int
}

type groupRepository struct {
	m sync.Map
}

func (c *groupRepository) Add(code int) {
	if _, ok := c.m.LoadOrStore(code, code); ok {
		log.Panicf("group error code duplicate definition, code: %d", code)
	}
}

func (c *groupRepository) Get(code int) int {
	val, ok := c.m.Load(code)
	if !ok {
		log.Panicf("error group code not exists, code: %s", code)
	}

	return val.(int)
}
