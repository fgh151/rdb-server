package plugins

type PluginInterface interface {
	Run(params PluginParams) PluginResult
}

type PluginParams struct {
	Data interface{}
}

type PluginResult struct {
	Payload interface{}
	Err     error
}
