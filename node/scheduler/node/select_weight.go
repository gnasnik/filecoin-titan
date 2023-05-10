package node

import (
	"math/rand"
	"sync"
	"time"

	"github.com/Filecoin-Titan/titan/node/modules/dtypes"
)

type weightManager struct {
	config dtypes.GetSchedulerConfigFunc

	// Each node assigned a select weight, when pulling resources, randomly select n select weight, and select the node holding these select weight.
	cSelectWeightLock          sync.RWMutex
	cSelectWeightRand          *rand.Rand
	cSelectWeightMax           int            // Candidate select weight , Distribute from 1
	cDistributedSelectWeight   map[int]string // Already allocated candidate select weights
	cUndistributedSelectWeight map[int]string // Undistributed candidate select weights

	eSelectWeightLock          sync.RWMutex
	eSelectWeightRand          *rand.Rand
	eSelectWeightMax           int            // Edge select weight , Distribute from 1
	eDistributedSelectWeight   map[int]string // Already allocated edge select weights
	eUndistributedSelectWeight map[int]string // Undistributed edge select weights
}

func newWeightManager(config dtypes.GetSchedulerConfigFunc) *weightManager {
	pullSelectSeed := time.Now().UnixNano()

	manager := &weightManager{
		cSelectWeightRand:          rand.New(rand.NewSource(pullSelectSeed)),
		eSelectWeightRand:          rand.New(rand.NewSource(pullSelectSeed)),
		cDistributedSelectWeight:   make(map[int]string),
		cUndistributedSelectWeight: make(map[int]string),
		eDistributedSelectWeight:   make(map[int]string),
		eUndistributedSelectWeight: make(map[int]string),
		config:                     config,
	}

	return manager
}

// distributeCandidateSelectWeight assigns undistributed select weight to node and returns the assigned weights
func (c *weightManager) distributeCandidateSelectWeight(nodeID string, n int) []int {
	c.cSelectWeightLock.Lock()
	defer c.cSelectWeightLock.Unlock()

	out := make([]int, 0)

	for i := 0; i < n; i++ {
		out = append(out, c.getCandidateSelectWeights())
	}

	for _, w := range out {
		// delete from Undistributed map
		delete(c.cUndistributedSelectWeight, w)
		// add to Distributed map
		c.cDistributedSelectWeight[w] = nodeID
	}

	return out
}

// distributeEdgeSelectWeight assigns undistributed select weight to node and returns the assigned weights
func (c *weightManager) distributeEdgeSelectWeight(nodeID string, n int) []int {
	c.eSelectWeightLock.Lock()
	defer c.eSelectWeightLock.Unlock()

	out := make([]int, 0)

	for i := 0; i < n; i++ {
		out = append(out, c.getEdgeSelectWeights())
	}

	for _, w := range out {
		// delete from Undistributed map
		delete(c.eUndistributedSelectWeight, w)
		// add to Distributed map
		c.eDistributedSelectWeight[w] = nodeID
	}

	return out
}

func (c *weightManager) getCandidateSelectWeights() int {
	if len(c.cUndistributedSelectWeight) > 0 {
		for w := range c.cUndistributedSelectWeight {
			return w
		}
	}

	c.cSelectWeightMax++
	return c.cSelectWeightMax
}

func (c *weightManager) getEdgeSelectWeights() int {
	if len(c.eUndistributedSelectWeight) > 0 {
		for w := range c.eUndistributedSelectWeight {
			return w
		}
	}

	c.eSelectWeightMax++
	return c.eSelectWeightMax
}

// repayCandidateSelectWeight repay the selection weight to cUndistributedSelectWeight
func (c *weightManager) repayCandidateSelectWeight(weights []int) {
	c.cSelectWeightLock.Lock()
	defer c.cSelectWeightLock.Unlock()

	for _, w := range weights {
		delete(c.cDistributedSelectWeight, w)
		c.cUndistributedSelectWeight[w] = ""
	}
}

// repayEdgeSelectWeight repay the selection weight to eUndistributedSelectWeight
func (c *weightManager) repayEdgeSelectWeight(weights []int) {
	c.eSelectWeightLock.Lock()
	defer c.eSelectWeightLock.Unlock()

	for _, w := range weights {
		delete(c.eDistributedSelectWeight, w)
		c.eUndistributedSelectWeight[w] = ""
	}
}

func (c *weightManager) getCandidateSelectWeightRandom() (string, int) {
	c.cSelectWeightLock.Lock()
	defer c.cSelectWeightLock.Unlock()

	w := c.cSelectWeightRand.Intn(c.cSelectWeightMax) + 1
	return c.cDistributedSelectWeight[w], w
}

func (c *weightManager) getEdgeSelectWeightRandom() (string, int) {
	c.eSelectWeightLock.Lock()
	defer c.eSelectWeightLock.Unlock()

	w := c.eSelectWeightRand.Intn(c.eSelectWeightMax) + 1
	return c.eDistributedSelectWeight[w], w
}

func (c *weightManager) cleanSelectWeights() {
	c.cSelectWeightLock.Lock()
	defer c.cSelectWeightLock.Unlock()

	c.eSelectWeightLock.Lock()
	defer c.eSelectWeightLock.Unlock()

	c.cDistributedSelectWeight = make(map[int]string)
	c.cUndistributedSelectWeight = make(map[int]string)
	c.eDistributedSelectWeight = make(map[int]string)
	c.eUndistributedSelectWeight = make(map[int]string)

	c.cSelectWeightMax = 0
	c.eSelectWeightMax = 0
}

func (c *weightManager) getWeightScale() map[string]int {
	cfg, err := c.config()
	if err != nil {
		log.Errorf("get config err:%s", err.Error())
		return map[string]int{}
	}

	return cfg.LevelSelectWeight
}

func (c *weightManager) getSelectWeightNum(scoreLevel string) int {
	num, exist := c.getWeightScale()[scoreLevel]
	if exist {
		return num
	}

	return 0
}
