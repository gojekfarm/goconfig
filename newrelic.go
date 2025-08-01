package goconfig

import (
	"github.com/newrelic/go-agent"
)

func getNewRelicConfigOrPanic(loader *ConfigLoader) newrelic.Config {
	config := newrelic.NewConfig(getStringOrPanic(loader, "NEW_RELIC_APP_NAME"), getStringOrPanic(loader, "NEW_RELIC_LICENCE_KEY"))
	config.Enabled = getFeature(loader, "NEW_RELIC_ENABLED")
	labels, err := parseNewRelicLabels(getStringOrPanic(loader, "NEW_RELIC_LABELS"))
	if err == nil {
		config.Labels = labels
	}
	return config
}
