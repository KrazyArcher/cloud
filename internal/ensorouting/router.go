package ensorouting

import (
    "crypto/sha256"
    "encoding/binary"
    "sync"
)

// NodeState tracks active KV-cache prefix affinity and load per backend node.
type NodeState struct {
    ID          string
    ActiveLoad  int
    PrefixCache map[uint64]bool
}

// RouteDecision represents the Enso prompt-locality routing output (HIP-1345).
type RouteDecision struct {
    NodeID   string
    CacheHit bool
    Score    float64
}

// PromptLocalityRouter routes prompts to nodes holding matching KV prefixes.
type PromptLocalityRouter struct {
    mu    sync.RWMutex
    nodes []*NodeState
}

func NewPromptLocalityRouter(nodeIDs []string) *PromptLocalityRouter {
    nodes := make([]*NodeState, 0, len(nodeIDs))
    for _, id := range nodeIDs {
        nodes = append(nodes, &NodeState{
            ID:          id,
            PrefixCache: make(map[uint64]bool),
        })
    }
    return &PromptLocalityRouter{nodes: nodes}
}

func prefixHash(prompt string) uint64 {
    if len(prompt) > 256 {
        prompt = prompt[:256]
    }
    sum := sha256.Sum256([]byte(prompt))
    return binary.LittleEndian.Uint64(sum[:8])
}

// Route selects the optimal node balancing KV-cache locality and active load.
func (r *PromptLocalityRouter) Route(prompt string) RouteDecision {
    r.mu.Lock()
    defer r.mu.Unlock()
    if len(r.nodes) == 0 {
        return RouteDecision{}
    }

    h := prefixHash(prompt)
    bestIdx := 0
    bestScore := -1e9
    bestHit := false

    for i, n := range r.nodes {
        hit := n.PrefixCache[h]
        localityBonus := 0.0
        if hit {
            localityBonus = 10.0
        }
        score := localityBonus - float64(n.ActiveLoad)*0.5
        if score > bestScore {
            bestScore = score
            bestIdx = i
            bestHit = hit
        }
    }

    chosen := r.nodes[bestIdx]
    chosen.PrefixCache[h] = true
    chosen.ActiveLoad++
    return RouteDecision{
        NodeID:   chosen.ID,
        CacheHit: bestHit,
        Score:    bestScore,
    }
}