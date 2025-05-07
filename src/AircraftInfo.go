package main

type AircraftInfo struct {
	Hexcode      string  `json:"HexCode"`
	Lat          float32 `json:"Lat"`
	Long         float32 `json:"Lon"`
	Callsign     string  `json:"Callsign"`
	AcType       string  `json:"Type"`
	AcReg        string  `json:"Reg"`
	Squawk       string  `json:"Squawk"`
	DbFlag       uint8   `json:"Mil"`
	AcCat        uint8   `json:"Category"`
	AcSrc        uint8   `json:"Source"`
	HbarFt       uint16  `json:"Alt.baro,ft"`
	HgeoFt       uint16  `json:"Alt.geo,ft"`
	SpdGroundKts float32 `json:"Gr.Speed,kts"`
	Timestamp    uint64  `json:"Timestamp"`
}
