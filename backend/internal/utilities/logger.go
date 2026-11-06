package utilities

import "github.com/rs/zerolog"

// LogInfo menambahkan Service, Method, dan ID opsional secara konsisten
func LogInfo(
	logger *zerolog.Logger,
	service,
	method string,
	ids map[string]string,
	msg string,
) {
	evt := logger.Info().
		Str("Service", service).
		Str("Method", method)

	for k, v := range ids {
		if v != "" {
			evt = evt.Str(k, v)
		}
	}

	evt.Msg(msg)
}

// LogWarn menambahkan Service, Method, dan ID opsional secara konsisten
func LogWarn(
	logger *zerolog.Logger,
	service,
	method string,
	ids map[string]string, msg string,
) {
	evt := logger.Warn().
		Str("Service", service).
		Str("Method", method)

	for k, v := range ids {
		if v != "" {
			evt = evt.Str(k, v)
		}
	}

	evt.Msg(msg)
}

// LogError menambahkan Service, Method, ID opsional, dan error secara konsisten
func LogError(
	logger *zerolog.Logger,
	service,
	method string,
	ids map[string]string,
	err error,
	msg string,
) {
	evt := logger.Error().
		Str("Service", service).
		Str("Method", method)

	for k, v := range ids {
		if v != "" {
			evt = evt.Str(k, v)
		}
	}

	evt.Err(err).Msg(msg)
}
