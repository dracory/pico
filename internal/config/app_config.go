package config

func appConfig(env *envValidator) appSettings {
	name := env.GetStringOrDefault(KEY_APP_NAME, "Pico")
	url := env.GetStringOrDefault(KEY_APP_URL, "http://localhost:8080")
	host := env.GetStringOrError(KEY_APP_HOST, "set the application host address")
	port := env.GetStringOrError(KEY_APP_PORT, "set the application port")
	appEnv := env.GetStringOrError(KEY_APP_ENVIRONMENT, "set the application environment")
	debug := env.GetBool(KEY_APP_DEBUG)

	return appSettings{
		name:  name,
		url:   url,
		host:  host,
		port:  port,
		env:   appEnv,
		debug: debug,
	}
}

type appSettings struct {
	name  string
	url   string
	host  string
	port  string
	env   string
	debug bool
}
