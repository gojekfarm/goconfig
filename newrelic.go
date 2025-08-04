package goconfig

import (
	"github.com/newrelic/go-agent"
)

func getNewRelicConfigOrPanic(accessor *YamlConfigAccessor) newrelic.Config {
	config := newrelic.NewConfig(getStringOrPanic(accessor, "NEW_RELIC_APP_NAME"), getStringOrPanic(accessor, "NEW_RELIC_LICENCE_KEY"))
	config.Enabled = getFeature(accessor, "NEW_RELIC_ENABLED")
	labels, err := parseNewRelicLabels(getStringOrPanic(accessor, "NEW_RELIC_LABELS"))
	if err == nil {
		config.Labels = labels
	}
	return config
}
