package goconfig

import (
	"github.com/newrelic/go-agent"
)

func getNewRelicConfigOrPanic(accessor ConfigAccessor) newrelic.Config {
	config := newrelic.NewConfig(getString(accessor, "NEW_RELIC_APP_NAME"), getString(accessor, "NEW_RELIC_LICENCE_KEY"))
	config.Enabled = getFeature(accessor, "NEW_RELIC_ENABLED")
	labels, err := parseNewRelicLabels(getString(accessor, "NEW_RELIC_LABELS"))
	if err == nil {
		config.Labels = labels
	}
	return config
}
