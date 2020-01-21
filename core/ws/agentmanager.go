package ws

import (
	"sync"

	"github.com/olahol/melody"
)

type agentManager struct {
	sessions map[*melody.Session]*AgentStatus
	agents   map[uint]*melody.Session
	sync.RWMutex
}

func newAgentManager() *agentManager {
	return &agentManager{
		sessions: make(map[*melody.Session]*AgentStatus),
		agents:   make(map[uint]*melody.Session),
	}
}

func (am *agentManager) add(s *melody.Session, agent *AgentStatus) bool {
	am.Lock()
	defer am.Unlock()
	id := agent.ID
	if _, ok := am.agents[id]; ok {
		return false
	}
	am.agents[id] = s
	am.sessions[s] = agent
	return true
}

func (am *agentManager) remove(s *melody.Session) bool {
	am.Lock()
	defer am.Unlock()
	if agent, ok := am.sessions[s]; ok {
		delete(am.agents, agent.ID)
		delete(am.sessions, s)
	}
	return true
}

func (am *agentManager) get(s *melody.Session) *AgentStatus {
	am.RLock()
	defer am.RUnlock()
	agent, _ := am.sessions[s]
	return agent
}

func (am *agentManager) getsession(id uint) *melody.Session {
	am.RLock()
	defer am.RUnlock()
	s, _ := am.agents[id]
	return s
}
