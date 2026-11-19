package ensorouting

import "testing"

func TestPromptLocalityRouter_AffinityHit(t *testing.T) {
    router := NewPromptLocalityRouter([]string{"node-a", "node-b", "node-c"})
    prompt := "system: you are a high-throughput inference router for HIP-1345"

    first := router.Route(prompt)
    if first.NodeID == "" {
        t.Fatalf("expected valid node ID, got empty")
    }
    if first.CacheHit {
        t.Fatalf("expected cold miss on first request")
    }

    second := router.Route(prompt)
    if second.NodeID != first.NodeID {
        t.Fatalf("expected locality affinity to %s, got %s", first.NodeID, second.NodeID)
    }
    if !second.CacheHit {
        t.Fatalf("expected KV-cache hit on repeated prefix")
    }
}