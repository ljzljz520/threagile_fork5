/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

package types

import (
	"fmt"
	"sort"
	"strings"
)

type TrustBoundary struct {
	Id                    string            `json:"id,omitempty" yaml:"id,omitempty"`
	Title                 string            `json:"title,omitempty" yaml:"title,omitempty"`
	Description           string            `json:"description,omitempty" yaml:"description,omitempty"`
	Type                  TrustBoundaryType `json:"type,omitempty" yaml:"type,omitempty"`
	Tags                  []string          `json:"tags,omitempty" yaml:"tags,omitempty"`
	TechnicalAssetsInside []string          `json:"technical_assets_inside,omitempty" yaml:"technical_assets_inside,omitempty"`
	TrustBoundariesNested []string          `json:"trust_boundaries_nested,omitempty" yaml:"trust_boundaries_nested,omitempty"`
}

func (what TrustBoundary) IsTaggedWithAny(tags ...string) bool {
	return containsCaseInsensitiveAny(what.Tags, tags...)
}

// BuildTrustBoundaryTopology derives and validates the trust-boundary containment topology of
// the model from the declared TrustBoundariesNested relationships and the directly containing
// trust boundary of each technical asset.
//
// It rejects ambiguous or cyclic nesting (a trust boundary nested within multiple parent trust
// boundaries, self-nesting, and longer parent cycles), while assets without any surrounding
// trust boundary are perfectly valid and simply stay unlisted. On success the model holds the
// validated parent mapping, complete ancestor chains, and the nearest network domain per asset,
// which network-related risk rules and the diagram rendering share as common facts.
func (model *Model) BuildTrustBoundaryTopology() error {
	parentById := make(map[string]*TrustBoundary)
	// iterate in sorted order so that any validation error is reproducible
	for _, parentId := range model.sortedTrustBoundaryIds() {
		parent := model.TrustBoundaries[parentId]
		for _, nestedBoundaryId := range parent.TrustBoundariesNested {
			if nestedBoundaryId == parentId {
				return fmt.Errorf("trust boundary %q nests within itself (cyclic trust boundary nesting)", nestedBoundaryId)
			}
			if existingParent, alreadyNested := parentById[nestedBoundaryId]; alreadyNested && existingParent.Id != parentId {
				return fmt.Errorf("trust boundary %q is nested within multiple parent trust boundaries (%q and %q): "+
					"a trust boundary can only be nested within a single parent trust boundary",
					nestedBoundaryId, existingParent.Id, parentId)
			}
			parentById[nestedBoundaryId] = parent
		}
	}

	ancestorsById := make(map[string][]*TrustBoundary)
	for _, id := range model.sortedTrustBoundaryIds() {
		chain := make([]*TrustBoundary, 0)
		visited := map[string]bool{id: true}
		path := []string{id}
		current := parentById[id]
		for current != nil {
			if visited[current.Id] {
				return fmt.Errorf("cyclic trust boundary nesting detected: %s", strings.Join(append(path, current.Id), " -> "))
			}
			visited[current.Id] = true
			path = append(path, current.Id)
			chain = append(chain, current)
			current = parentById[current.Id]
		}
		ancestorsById[id] = chain
	}

	nearestNetworkByAssetId := make(map[string]*TrustBoundary)
	for _, technicalAssetId := range model.SortedTechnicalAssetIDs() {
		direct := model.DirectContainingTrustBoundaryMappedByTechnicalAssetId[technicalAssetId]
		if direct == nil {
			continue // legitimate: asset is modeled without any surrounding trust boundary
		}
		if direct.Type.IsNetworkBoundary() {
			nearestNetworkByAssetId[technicalAssetId] = direct
			continue
		}
		// execution environments and other logical boundaries are transparent for network
		// domain purposes: walk the complete ancestor chain until a network boundary is found
		for _, ancestor := range ancestorsById[direct.Id] {
			if ancestor.Type.IsNetworkBoundary() {
				nearestNetworkByAssetId[technicalAssetId] = ancestor
				break
			}
		}
	}

	model.ParentTrustBoundaryMappedByTrustBoundaryId = parentById
	model.TrustBoundaryAncestorsMappedByTrustBoundaryId = ancestorsById
	model.NearestNetworkTrustBoundaryMappedByTechnicalAssetId = nearestNetworkByAssetId
	return nil
}

// TrustBoundaryAncestors returns all trust boundaries containing the given one, ordered from
// the direct parent outwards to the outermost root. Execution environments and network
// boundaries are both included; use NearestNetworkTrustBoundary to resolve the network domain.
func (model *Model) TrustBoundaryAncestors(what *TrustBoundary) []*TrustBoundary {
	result := make([]*TrustBoundary, 0)
	if what == nil {
		return result
	}
	if model.TrustBoundaryAncestorsMappedByTrustBoundaryId != nil {
		return append(result, model.TrustBoundaryAncestorsMappedByTrustBoundaryId[what.Id]...)
	}
	// fallback for models constructed directly in code (e.g. tests) without parsed topology
	visited := map[string]bool{what.Id: true}
	current := model.FindParentTrustBoundary(what)
	for current != nil {
		if visited[current.Id] {
			break
		}
		visited[current.Id] = true
		result = append(result, current)
		current = model.FindParentTrustBoundary(current)
	}
	return result
}

// TrustBoundaryLineage returns the complete containment chain of a technical asset, starting
// with its directly containing trust boundary and continuing outwards to the outermost root.
func (model *Model) TrustBoundaryLineage(technicalAssetId string) []*TrustBoundary {
	lineage := make([]*TrustBoundary, 0)
	direct := model.DirectContainingTrustBoundaryMappedByTechnicalAssetId[technicalAssetId]
	if direct == nil {
		return lineage
	}
	lineage = append(lineage, direct)
	lineage = append(lineage, model.TrustBoundaryAncestors(direct)...)
	return lineage
}

// NearestNetworkTrustBoundary returns the nearest (deepest) network trust boundary on the
// containment chain of the technical asset, transparently walking through nested execution
// environments and other logical boundaries. It returns nil when neither the asset nor any of
// its ancestor boundaries is a network boundary (including for assets without any boundary).
func (model *Model) NearestNetworkTrustBoundary(technicalAssetId string) *TrustBoundary {
	if model.NearestNetworkTrustBoundaryMappedByTechnicalAssetId != nil {
		return model.NearestNetworkTrustBoundaryMappedByTechnicalAssetId[technicalAssetId]
	}
	// fallback for models constructed directly in code (e.g. tests) without parsed topology
	for _, candidate := range model.TrustBoundaryLineage(technicalAssetId) {
		if candidate.Type.IsNetworkBoundary() {
			return candidate
		}
	}
	return nil
}

func (model *Model) sortedTrustBoundaryIds() []string {
	ids := make([]string, 0, len(model.TrustBoundaries))
	for id := range model.TrustBoundaries {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}
