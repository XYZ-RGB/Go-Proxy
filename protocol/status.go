package protocol

import (
	"bytes"
	"encoding/json"
	"io"
)

type ClientStatusRequest struct {
}

func (c ClientStatusRequest) Write(buffer *bytes.Buffer) error {
	return nil
}

func (c *ClientStatusRequest) Read(session io.Reader) error {
	return nil
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

func (s ServerStatusResponse) Write(buffer *bytes.Buffer) error {
	statusData, err := json.Marshal(s.Status)
	if err == nil {
		String(statusData).Write(buffer)
	} else {
		return err
	}
	return nil
}

func (s *ServerStatusResponse) Read(session io.Reader) error {
	var statusData String;
	statusData.Read(session);
	return json.Unmarshal([]byte(statusData), &s.Status)
}

func (s ServerStatusResponse) Id() VarInt {
	return 0x00
}
