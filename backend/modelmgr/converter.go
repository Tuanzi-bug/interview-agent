package modelmgr

import "ai-eino-interview-agent/modelmgr/manager"

// convertToModel 将内部模型转换为公共模型。
func convertToModel(im *manager.InternalModel) *Model {
	description := ""
	if im.Description != nil {
		if im.Description.EN != "" {
			description = im.Description.EN
		} else if im.Description.ZH != "" {
			description = im.Description.ZH
		}
	}

	displayName := im.DisplayName
	if displayName == "" {
		displayName = im.Name
	}

	model := &Model{
		ID:          im.ID,
		Name:        im.Name,
		DisplayName: displayName,
		IconURL:     im.IconURL,
		Protocol:    Protocol(im.Meta.Protocol),
		Status:      ModelStatus(im.Meta.Status),
		Description: description,
	}

	if im.Meta.Capability != nil {
		model.Capability = &ModelCapability{
			SupportFunctionCall: im.Meta.Capability.FunctionCall,
			SupportStreaming:    im.Meta.Capability.Streaming,
			SupportVision:       im.Meta.Capability.Vision,
			MaxInputTokens:      im.Meta.Capability.InputTokens,
			MaxOutputTokens:     im.Meta.Capability.OutputTokens,
			SupportedModalities: im.Meta.Capability.InputModal,
		}
	}

	return model
}
