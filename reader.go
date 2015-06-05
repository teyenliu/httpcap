package main

type RAWData struct {
	Data       []byte
	SrcPort    uint16
	DestPort   uint16
	LocalAddr  string
	RemoteAddr string
}

type InputReader interface {
	Read(data []byte) (int, uint16, uint16, string, string, error)
}