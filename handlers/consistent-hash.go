package main

import (
    "hash/crc32"
    "sort"
    "strconv"
)

type Node struct {
    ID   string
    Addr string
}

type HashRing struct {
    Nodes []Node
    Ring  map[uint32]Node
}

func NewHashRing(nodes []Node) *HashRing {
    ring := &HashRing{
        Nodes: nodes,
        Ring:  make(map[uint32]Node),
    }

    for _, node := range nodes {
        hash := ring.hashKey(node.ID)
        ring.Ring[hash] = node
    }

    return ring
}

func (h *HashRing) hashKey(key string) uint32 {
    return crc32.ChecksumIEEE([]byte(key))
}

func (h *HashRing) GetNode(key string) Node {
    hashedKey := h.hashKey(key)
    keys := []uint32{}
    for k := range h.Ring {
        keys = append(keys, k)
    }

    sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

    for _, k := range keys {
        if hashedKey <= k {
            return h.Ring[k]
        }
    }

    return h.Ring[keys[0]]
}
