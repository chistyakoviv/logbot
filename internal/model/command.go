package model

import "time"

const NoStage int = -1

type Command struct {
	Name      string
	UserId    int64
	ChatId    int64
	Stage     int
	Data      map[string]interface{}
	UpdatedAt time.Time
}

func (c *Command) IsInProgress() bool {
	return c.Stage != NoStage
}

type CommandInfo struct {
	Name   string
	ChatId int64
	UserId int64
	Stage  int
	Data   map[string]interface{}
}

type CommandKey struct {
	ChatId int64
	UserId int64
}

type CommandData struct {
	Name  string
	Stage int
	Data  map[string]interface{}
}
