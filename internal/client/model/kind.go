package model

type Kind string

const (
	KindConfiguration Kind = "Configuration"
	KindSource        Kind = "Source"
	KindProcessor     Kind = "Processor"
	KindDestination   Kind = "Destination"
)
