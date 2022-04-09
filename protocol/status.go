package protocol

import (
	"bytes"
	"encoding/json"
)

type ClientStatusRequest struct {
}

func (c ClientStatusRequest) Write(buffer *bytes.Buffer) {
}

func (c *ClientStatusRequest) Read(session Session) {
}

func (c ClientStatusRequest) Id() VarInt {
	return 0x00
}

type VersionStatusData struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}
type PlayersStatusData struct {
	Max    int                `json:"max"`
	Online int                `json:"online"`
	Sample []PlayerDataStatus `json:"sample"`
}

type PlayerDataStatus struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type DescriptionStatusData struct {
	Text string `json:"text"`
}

type StatusData struct {
	Version     VersionStatusData     `json:"version"`
	Players     PlayersStatusData     `json:"players"`
	Description DescriptionStatusData `json:"description"`
	Favicon     string                `json:"favicon"`
}

type ServerStatusResponse struct {
	Status StatusData
}

func (s ServerStatusResponse) Write(buffer *bytes.Buffer) {
	statusData, _ := json.Marshal(s.Status)
	String(statusData).Write(buffer)

}

func (s *ServerStatusResponse) Read(session Session) {
	//todo
}

func (s ServerStatusResponse) Id() VarInt {
	return 0x00
}
