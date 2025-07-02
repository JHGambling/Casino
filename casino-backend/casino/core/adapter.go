package core

import "jhgambling/protocol"

type CasinoPluginAdapter struct {
	core *CasinoCore
}

func NewCasinoPluginAdapter(core *CasinoCore) *CasinoPluginAdapter {
	return &CasinoPluginAdapter{
		core: core,
	}
}

func (a *CasinoPluginAdapter) Table(id string) (protocol.Table, error) {
	return a.core.Database.GetTable(id)
}
