package model

type AuthLevel uint8

const (
	READ AuthLevel = iota
	ADD
	UPDATE
	DEL
)
