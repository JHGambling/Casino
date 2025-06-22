package main

import "jhgambling/protocol"

type ExampleProvider struct{}

func (p *ExampleProvider) GetID() string {
	return "example"
}
func (p *ExampleProvider) GetName() string {
	return "Example Provider"
}

var Provider protocol.GameProvider = &ExampleProvider{}
