package todoreminder

import (
	core_feature "github.com/moond0wner/morning-digest/internal/core/feature"
)

type Feature struct {
	core_feature.BaseFeature
	handler *handler
}
