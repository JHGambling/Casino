package plugins

import (
	"errors"
	"io/ioutil"
	"jhgambling/backend/core/utils"
	"path"
	"plugin"
	"strings"

	"jhgambling/protocol"
)

type PluginManager struct {
}

func NewPluginManager() *PluginManager {
	return &PluginManager{}
}

func (pm *PluginManager) LoadPlugins() {
	files := pm.ListAvailablePlugins()

	for _, f := range files {
		provider, err := pm.LoadGamePlugin(f)
		if err != nil {
			utils.Log("error", "casino::plugins", "failed to load plugin: ", err)
			continue
		}

		utils.Log("ok", "casino::plugins", "loaded plugin '", provider.GetID(), "' with name '", provider.GetName(), "'")
	}
}

func (pm *PluginManager) ListAvailablePlugins() []string {
	pluginPath := "../games/"

	result := []string{}

	files, err := ioutil.ReadDir(pluginPath)
	if err != nil {
		utils.Log("error", "casino::plugins", "failed to list plugins directory: ", err)
		return result
	}

	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".so") {
			continue
		}

		result = append(result, path.Join(pluginPath, f.Name()))
	}

	utils.Log("info", "casino::plugins", "discovered ", len(result), " plugin(s)")

	return result
}

func (pm *PluginManager) LoadGamePlugin(path string) (protocol.GameProvider, error) {
	p, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}
	sym, err := p.Lookup("Provider")
	if err != nil {
		return nil, err
	}
	providerPtr, ok := sym.(*protocol.GameProvider)
	if !ok {
		return nil, errors.New("plugin does not implement provider interface")
	}
	return *providerPtr, nil
}
