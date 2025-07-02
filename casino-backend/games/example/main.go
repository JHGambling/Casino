package main

import "jhgambling/protocol"

type ExampleProvider struct{}

func (p *ExampleProvider) GetID() string {
	return "example"
}
func (p *ExampleProvider) GetName() string {
	return "Example Provider"
}
func (p *ExampleProvider) GetInstances() []protocol.GameInstance {
	return []protocol.GameInstance{}
}

var Provider protocol.GameProvider = &ExampleProvider{}
