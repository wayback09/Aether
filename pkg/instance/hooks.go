package instance

// ModLoaderHook is set by the extensions package after init to avoid import cycles.
// It receives a launch context and returns a modified context.
var ModLoaderHook func(loaderID string, ctx map[string]interface{}) (map[string]interface{}, error)
