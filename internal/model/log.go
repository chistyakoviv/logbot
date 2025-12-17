package model

import (
	"time"
)

type Log struct {
	Id            int
	Token         string
	Data          string
	Service       string
	ContainerName string
	ContainerId   string
	Node          string
	NodeId        string
	Hash          string
	CreatedAt     time.Time
}

type LogInfo struct {
	Token         string
	Data          string
	Service       string
	ContainerName string
	ContainerId   string
	Node          string
	NodeId        string
}
