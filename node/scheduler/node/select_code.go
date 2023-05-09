package node

import (
	"math/rand"
	"sync"
	"time"
)

type codeManager struct {
	// Each node assigned a select code, when pulling resources, randomly select n select code, and select the node holding these select code.
	cSelectCodeLock          sync.RWMutex
	cSelectCodeRand          *rand.Rand
	cSelectCodeMax           int            // Candidate select code , Distribute from 1
	cDistributedSelectCode   map[int]string // Already allocated candidate select codes
	cUndistributedSelectCode map[int]string // Undistributed candidate select codes

	eSelectCodeLock          sync.RWMutex
	eSelectCodeRand          *rand.Rand
	eSelectCodeMax           int            // Edge select code , Distribute from 1
	eDistributedSelectCode   map[int]string // Already allocated edge select codes
	eUndistributedSelectCode map[int]string // Undistributed edge select codes
}

func newCodeManager() *codeManager {
	pullSelectSeed := time.Now().UnixNano()

	manager := &codeManager{
		cSelectCodeRand:          rand.New(rand.NewSource(pullSelectSeed)),
		eSelectCodeRand:          rand.New(rand.NewSource(pullSelectSeed)),
		cDistributedSelectCode:   make(map[int]string),
		cUndistributedSelectCode: make(map[int]string),
		eDistributedSelectCode:   make(map[int]string),
		eUndistributedSelectCode: make(map[int]string),
	}

	return manager
}

// distributeCandidateSelectCode assigns undistributed select code to node and returns the assigned codes
func (c *codeManager) distributeCandidateSelectCode(nodeID string, n int) []int {
	c.cSelectCodeLock.Lock()
	defer c.cSelectCodeLock.Unlock()

	selectCodes := make([]int, 0)

	for i := 0; i < n; i++ {
		selectCodes = append(selectCodes, c.getCandidateSelectCodes())
	}

	for _, code := range selectCodes {
		// delete from Undistributed map
		delete(c.cUndistributedSelectCode, code)
		// add to Distributed map
		c.cDistributedSelectCode[code] = nodeID
	}

	return selectCodes
}

// distributeEdgeSelectCode assigns undistributed select code to node and returns the assigned codes
func (c *codeManager) distributeEdgeSelectCode(nodeID string, n int) []int {
	c.eSelectCodeLock.Lock()
	defer c.eSelectCodeLock.Unlock()

	selectCodes := make([]int, 0)

	for i := 0; i < n; i++ {
		selectCodes = append(selectCodes, c.getEdgeSelectCodes())
	}

	for _, code := range selectCodes {
		// delete from Undistributed map
		delete(c.eUndistributedSelectCode, code)
		// add to Distributed map
		c.eDistributedSelectCode[code] = nodeID
	}

	return selectCodes
}

func (c *codeManager) getCandidateSelectCodes() int {
	if len(c.cUndistributedSelectCode) > 0 {
		for code := range c.cUndistributedSelectCode {
			return code
		}
	}

	c.cSelectCodeMax++
	return c.cSelectCodeMax
}

func (c *codeManager) getEdgeSelectCodes() int {
	if len(c.eUndistributedSelectCode) > 0 {
		for code := range c.eUndistributedSelectCode {
			return code
		}
	}

	c.eSelectCodeMax++
	return c.eSelectCodeMax
}

// repayCandidateSelectCode repay the selection code to cUndistributedSelectCode
func (c *codeManager) repayCandidateSelectCode(codes []int) {
	c.cSelectCodeLock.Lock()
	defer c.cSelectCodeLock.Unlock()

	for _, code := range codes {
		delete(c.cDistributedSelectCode, code)
		c.cUndistributedSelectCode[code] = ""
	}
}

// repayEdgeSelectCode repay the selection code to eUndistributedSelectCode
func (c *codeManager) repayEdgeSelectCode(codes []int) {
	c.eSelectCodeLock.Lock()
	defer c.eSelectCodeLock.Unlock()

	for _, code := range codes {
		delete(c.eDistributedSelectCode, code)
		c.eUndistributedSelectCode[code] = ""
	}
}

func (c *codeManager) getCandidateSelectCodeRandom() (string, int) {
	c.cSelectCodeLock.Lock()
	defer c.cSelectCodeLock.Unlock()

	code := c.cSelectCodeRand.Intn(c.cSelectCodeMax) + 1
	return c.cDistributedSelectCode[code], code
}

func (c *codeManager) getEdgeSelectCodeRandom() (string, int) {
	c.eSelectCodeLock.Lock()
	defer c.eSelectCodeLock.Unlock()

	code := c.eSelectCodeRand.Intn(c.eSelectCodeMax) + 1
	return c.eDistributedSelectCode[code], code
}

func (c *codeManager) cleanSelectCodes() {
	c.cSelectCodeLock.Lock()
	defer c.cSelectCodeLock.Unlock()

	c.eSelectCodeLock.Lock()
	defer c.eSelectCodeLock.Unlock()

	c.cDistributedSelectCode = make(map[int]string)
	c.cUndistributedSelectCode = make(map[int]string)
	c.eDistributedSelectCode = make(map[int]string)
	c.eUndistributedSelectCode = make(map[int]string)

	c.cSelectCodeMax = 0
	c.eSelectCodeMax = 0
}
