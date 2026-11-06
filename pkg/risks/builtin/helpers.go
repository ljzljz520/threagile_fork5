package builtin

import (
	"strings"

	"github.com/threagile/threagile/pkg/types"
)

// isAcrossTrustBoundaryNetworkOnly reports whether the communication link leaves the network
// domain of its source asset and enters a different one. The nearest network trust boundary is
// resolved from the validated trust-boundary topology, transparently walking through nested
// execution environments and other logical (non-network) boundaries on either side.
func isAcrossTrustBoundaryNetworkOnly(parsedModel *types.Model, communicationLink *types.CommunicationLink) bool {
	trustBoundaryOfSourceAsset := parsedModel.NearestNetworkTrustBoundary(communicationLink.SourceId)
	trustBoundaryOfTargetAsset := parsedModel.NearestNetworkTrustBoundary(communicationLink.TargetId)
	if trustBoundaryOfSourceAsset == nil || trustBoundaryOfTargetAsset == nil {
		return false
	}
	return trustBoundaryOfSourceAsset.Id != trustBoundaryOfTargetAsset.Id
}

func contains(as []string, b string) bool {
	for _, a := range as {
		if b == a {
			return true
		}
	}
	return false
}

func containsCaseInsensitiveAny(as []string, bs ...string) bool {
	for _, a := range as {
		for _, b := range bs {
			if strings.TrimSpace(strings.ToLower(b)) == strings.TrimSpace(strings.ToLower(a)) {
				return true
			}
		}
	}
	return false
}

func isSameExecutionEnvironment(parsedModel *types.Model, ta *types.TechnicalAsset, otherAssetId string) bool {
	trustBoundaryOfMyAsset, trustBoundaryOfMyAssetOk := parsedModel.DirectContainingTrustBoundaryMappedByTechnicalAssetId[ta.Id]
	trustBoundaryOfOtherAsset, trustBoundaryOfOtherAssetOk := parsedModel.DirectContainingTrustBoundaryMappedByTechnicalAssetId[otherAssetId]
	if trustBoundaryOfMyAssetOk != trustBoundaryOfOtherAssetOk {
		return false
	}
	if !trustBoundaryOfMyAssetOk {
		return true
	}
	if trustBoundaryOfMyAsset.Type == types.ExecutionEnvironment && trustBoundaryOfOtherAsset.Type == types.ExecutionEnvironment {
		return trustBoundaryOfMyAsset.Id == trustBoundaryOfOtherAsset.Id
	}
	return false
}

// isSameTrustBoundaryNetworkOnly reports whether both assets reside within the same network
// domain, using the nearest network trust boundary from the validated trust-boundary topology
// (nested execution environments and other logical boundaries are walked through). When neither
// asset resolves to any network boundary they are considered within the same domain, while a
// single-sided resolution counts as separated. This network-domain strategy intentionally
// differs from isSameExecutionEnvironment's exact execution-environment match.
func isSameTrustBoundaryNetworkOnly(parsedModel *types.Model, ta *types.TechnicalAsset, otherAssetId string) bool {
	trustBoundaryOfMyAsset := parsedModel.NearestNetworkTrustBoundary(ta.Id)
	trustBoundaryOfOtherAsset := parsedModel.NearestNetworkTrustBoundary(otherAssetId)
	if (trustBoundaryOfMyAsset == nil) != (trustBoundaryOfOtherAsset == nil) {
		return false
	}
	if trustBoundaryOfMyAsset == nil {
		return true
	}
	return trustBoundaryOfMyAsset.Id == trustBoundaryOfOtherAsset.Id
}
