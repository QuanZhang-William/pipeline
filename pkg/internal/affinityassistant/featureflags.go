package affinityassistant

import (
	"context"

	"github.com/tektoncd/pipeline/pkg/apis/config"
)

type AffinityAssitantBehavior string

const (
	AffinityAssistantDisabled                    = AffinityAssitantBehavior("AffinityAssistantDisabled")
	AffinityAssistantPerWorkspace                = AffinityAssitantBehavior("AffinityAssistantPerWorkspace")
	AffinityAssistantPerPipelineRun              = AffinityAssitantBehavior("AffinityAssistantPerPipelineRun")
	AffinityAssistantPerPipelineRunWithIsolation = AffinityAssitantBehavior("AffinityAssistantPerPipelineRunWithIsolation")
)

// TODO: add string doc
func GetAffinityAssistantBehavior(ctx context.Context) AffinityAssitantBehavior {
	cfg := config.FromContextOrDefaults(ctx)

	if !cfg.FeatureFlags.DisableAffinityAssistant {
		return AffinityAssistantPerWorkspace
	}

	switch cfg.FeatureFlags.Coscheduling {
	case config.CoscheduleWorkspaces:
		return AffinityAssistantPerWorkspace
	case config.CoschedulePipelineRuns:
		return AffinityAssistantPerPipelineRun
	case config.CoscheduleIsolatePipelineRuns:
		return AffinityAssistantPerPipelineRunWithIsolation
	case config.CoscheduleDisabled:
		return AffinityAssistantDisabled
	}

	return ""
}
