/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

package model

import (
	"sort"
	"testing"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/threagile/threagile/pkg/input"
	"github.com/threagile/threagile/pkg/risks/builtin"
	"github.com/threagile/threagile/pkg/types"
)

func TestDefaultInputNotFail(t *testing.T) {
	parsedModel, err := ParseModel(&mockConfig{}, createInputModel(make(map[string]input.TechnicalAsset), make(map[string]input.DataAsset)), make(types.RiskRules), make(types.RiskRules))

	assert.NoError(t, err)
	assert.NotNil(t, parsedModel)
}

func TestInferConfidentiality_NotSet_NoOthers_ExpectTODO(t *testing.T) {
	ta := make(map[string]input.TechnicalAsset)
	da := make(map[string]input.DataAsset)

	_, err := ParseModel(&mockConfig{}, createInputModel(ta, da), make(types.RiskRules), make(types.RiskRules))
	// TODO: rename test and check if everyone agree that by default it should be public if there are no other assets

	assert.NoError(t, err)
}

func TestInferConfidentiality_ExpectHighestConfidentiality(t *testing.T) {
	ta := make(map[string]input.TechnicalAsset)
	da := make(map[string]input.DataAsset)

	daConfidentialConfidentiality := createDataAsset(types.Confidential, types.Critical, types.Critical)
	da[daConfidentialConfidentiality.ID] = daConfidentialConfidentiality

	daRestrictedConfidentiality := createDataAsset(types.Restricted, types.Important, types.Important)
	da[daRestrictedConfidentiality.ID] = daRestrictedConfidentiality

	daPublicConfidentiality := createDataAsset(types.Public, types.Archive, types.Archive)
	da[daPublicConfidentiality.ID] = daPublicConfidentiality

	taWithConfidentialConfidentialityDataAsset := createTechnicalAsset(types.Internal, types.Operational, types.Operational)
	taWithConfidentialConfidentialityDataAsset.DataAssetsProcessed = append(taWithConfidentialConfidentialityDataAsset.DataAssetsProcessed, daConfidentialConfidentiality.ID)
	ta[taWithConfidentialConfidentialityDataAsset.ID] = taWithConfidentialConfidentialityDataAsset

	taWithRestrictedConfidentialityDataAsset := createTechnicalAsset(types.Internal, types.Operational, types.Operational)
	taWithRestrictedConfidentialityDataAsset.DataAssetsProcessed = append(taWithRestrictedConfidentialityDataAsset.DataAssetsProcessed, daRestrictedConfidentiality.ID)
	ta[taWithRestrictedConfidentialityDataAsset.ID] = taWithRestrictedConfidentialityDataAsset

	taWithPublicConfidentialityDataAsset := createTechnicalAsset(types.Internal, types.Operational, types.Operational)
	taWithPublicConfidentialityDataAsset.DataAssetsProcessed = append(taWithPublicConfidentialityDataAsset.DataAssetsProcessed, daPublicConfidentiality.ID)
	ta[taWithPublicConfidentialityDataAsset.ID] = taWithPublicConfidentialityDataAsset

	parsedModel, err := ParseModel(&mockConfig{}, createInputModel(ta, da), make(types.RiskRules), make(types.RiskRules))

	assert.NoError(t, err)
	assert.Equal(t, types.Confidential, parsedModel.TechnicalAssets[taWithConfidentialConfidentialityDataAsset.ID].Confidentiality)
	assert.Equal(t, types.Restricted, parsedModel.TechnicalAssets[taWithRestrictedConfidentialityDataAsset.ID].Confidentiality)
	assert.Equal(t, types.Internal, parsedModel.TechnicalAssets[taWithPublicConfidentialityDataAsset.ID].Confidentiality)
}

func TestInferIntegrity_NotSet_NoOthers_ExpectTODO(t *testing.T) {
	ta := make(map[string]input.TechnicalAsset)
	da := make(map[string]input.DataAsset)

	_, err := ParseModel(&mockConfig{}, createInputModel(ta, da), make(types.RiskRules), make(types.RiskRules))
	// TODO: rename test and check if everyone agree that by default it should be public if there are no other assets

	assert.NoError(t, err)
}

func TestInferIntegrity_ExpectHighestIntegrity(t *testing.T) {
	ta := make(map[string]input.TechnicalAsset)
	da := make(map[string]input.DataAsset)

	daCriticalIntegrity := createDataAsset(types.Confidential, types.Critical, types.Critical)
	da[daCriticalIntegrity.ID] = daCriticalIntegrity

	daImportantIntegrity := createDataAsset(types.Restricted, types.Important, types.Important)
	da[daImportantIntegrity.ID] = daImportantIntegrity

	daArchiveIntegrity := createDataAsset(types.Public, types.Archive, types.Archive)
	da[daArchiveIntegrity.ID] = daArchiveIntegrity

	taWithCriticalIntegrityDataAsset := createTechnicalAsset(types.Internal, types.Operational, types.Operational)
	taWithCriticalIntegrityDataAsset.DataAssetsProcessed = append(taWithCriticalIntegrityDataAsset.DataAssetsProcessed, daCriticalIntegrity.ID)
	ta[taWithCriticalIntegrityDataAsset.ID] = taWithCriticalIntegrityDataAsset

	taWithImportantIntegrityDataAsset := createTechnicalAsset(types.Internal, types.Operational, types.Operational)
	taWithImportantIntegrityDataAsset.DataAssetsProcessed = append(taWithImportantIntegrityDataAsset.DataAssetsProcessed, daImportantIntegrity.ID)
	ta[taWithImportantIntegrityDataAsset.ID] = taWithImportantIntegrityDataAsset

	taWithArchiveIntegrityDataAsset := createTechnicalAsset(types.Internal, types.Operational, types.Operational)
	taWithArchiveIntegrityDataAsset.DataAssetsProcessed = append(taWithArchiveIntegrityDataAsset.DataAssetsProcessed, daArchiveIntegrity.ID)
	ta[taWithArchiveIntegrityDataAsset.ID] = taWithArchiveIntegrityDataAsset

	parsedModel, err := ParseModel(&mockConfig{}, createInputModel(ta, da), make(types.RiskRules), make(types.RiskRules))

	assert.NoError(t, err)
	assert.Equal(t, types.Critical, parsedModel.TechnicalAssets[taWithCriticalIntegrityDataAsset.ID].Integrity)
	assert.Equal(t, types.Important, parsedModel.TechnicalAssets[taWithImportantIntegrityDataAsset.ID].Integrity)
	assert.Equal(t, types.Operational, parsedModel.TechnicalAssets[taWithArchiveIntegrityDataAsset.ID].Integrity)
}

func TestInferAvailability_NotSet_NoOthers_ExpectTODO(t *testing.T) {
	ta := make(map[string]input.TechnicalAsset)
	da := make(map[string]input.DataAsset)

	_, err := ParseModel(&mockConfig{}, createInputModel(ta, da), make(types.RiskRules), make(types.RiskRules))

	assert.NoError(t, err)
}

func TestInferAvailability_ExpectHighestAvailability(t *testing.T) {
	ta := make(map[string]input.TechnicalAsset)
	da := make(map[string]input.DataAsset)

	daCriticalAvailability := createDataAsset(types.Confidential, types.Critical, types.Critical)
	da[daCriticalAvailability.ID] = daCriticalAvailability

	daImportantAvailability := createDataAsset(types.Restricted, types.Important, types.Important)
	da[daImportantAvailability.ID] = daImportantAvailability

	daArchiveAvailability := createDataAsset(types.Public, types.Archive, types.Archive)
	da[daArchiveAvailability.ID] = daArchiveAvailability

	taWithCriticalAvailabilityDataAsset := createTechnicalAsset(types.Internal, types.Operational, types.Operational)
	taWithCriticalAvailabilityDataAsset.DataAssetsProcessed = append(taWithCriticalAvailabilityDataAsset.DataAssetsProcessed, daCriticalAvailability.ID)
	ta[taWithCriticalAvailabilityDataAsset.ID] = taWithCriticalAvailabilityDataAsset

	taWithImportantAvailabilityDataAsset := createTechnicalAsset(types.Internal, types.Operational, types.Operational)
	taWithImportantAvailabilityDataAsset.DataAssetsProcessed = append(taWithImportantAvailabilityDataAsset.DataAssetsProcessed, daImportantAvailability.ID)
	ta[taWithImportantAvailabilityDataAsset.ID] = taWithImportantAvailabilityDataAsset

	taWithArchiveAvailabilityDataAsset := createTechnicalAsset(types.Internal, types.Operational, types.Operational)
	taWithArchiveAvailabilityDataAsset.DataAssetsProcessed = append(taWithArchiveAvailabilityDataAsset.DataAssetsProcessed, daArchiveAvailability.ID)
	ta[taWithArchiveAvailabilityDataAsset.ID] = taWithArchiveAvailabilityDataAsset

	parsedModel, err := ParseModel(&mockConfig{}, createInputModel(ta, da), make(types.RiskRules), make(types.RiskRules))

	assert.NoError(t, err)
	assert.Equal(t, types.Critical, parsedModel.TechnicalAssets[taWithCriticalAvailabilityDataAsset.ID].Availability)
	assert.Equal(t, types.Important, parsedModel.TechnicalAssets[taWithImportantAvailabilityDataAsset.ID].Availability)
	assert.Equal(t, types.Operational, parsedModel.TechnicalAssets[taWithArchiveAvailabilityDataAsset.ID].Availability)
}

func createInputModel(technicalAssets map[string]input.TechnicalAsset, dataAssets map[string]input.DataAsset) *input.Model {
	return &input.Model{
		TechnicalAssets: technicalAssets,
		DataAssets:      dataAssets,

		// set some dummy values to bypass validation
		BusinessCriticality: "archive",
	}
}

func createTechnicalAsset(confidentiality types.Confidentiality, integrity types.Criticality, availability types.Criticality) input.TechnicalAsset {
	return input.TechnicalAsset{
		ID: uuid.New().String(),
		// those values are required to bypass validation
		Usage:           "business",
		Type:            "process",
		Size:            "system",
		Technology:      "unknown-technology",
		Encryption:      "none",
		Machine:         "virtual",
		Confidentiality: confidentiality.String(),
		Integrity:       integrity.String(),
		Availability:    availability.String(),
	}
}

func createDataAsset(confidentiality types.Confidentiality, integrity types.Criticality, availability types.Criticality) input.DataAsset {
	return input.DataAsset{
		ID:              uuid.New().String(),
		Usage:           "business",
		Quantity:        "few",
		Confidentiality: confidentiality.String(),
		Integrity:       integrity.String(),
		Availability:    availability.String(),
	}
}

type mockConfig struct {
}

func (m *mockConfig) GetAppFolder() string {
	return ""
}

func (m *mockConfig) GetTechnologyFilename() string {
	return ""
}

// =========================================================================
// trust-boundary topology
// =========================================================================

func topologyTechnicalAsset(id, assetType string) input.TechnicalAsset {
	return input.TechnicalAsset{
		ID:              id,
		Usage:           types.Business.String(),
		Type:            assetType,
		Size:            "system",
		Technology:      "unknown-technology",
		Encryption:      types.NoneEncryption.String(),
		Machine:         types.Virtual.String(),
		Confidentiality: types.Public.String(),
		Integrity:       types.Archive.String(),
		Availability:    types.Archive.String(),
	}
}

func topologyTrustBoundary(id string, tbType types.TrustBoundaryType, assetsInside []string, nested []string) input.TrustBoundary {
	return input.TrustBoundary{
		ID:                    id,
		Type:                  tbType.String(),
		TechnicalAssetsInside: assetsInside,
		TrustBoundariesNested: nested,
	}
}

func createNestedTrustBoundaryInputModel(trustBoundaries map[string]input.TrustBoundary) *input.Model {
	web := topologyTechnicalAsset("web", types.Process.String())
	web.CommunicationLinks = map[string]input.CommunicationLink{
		"Call DB": {
			Target:         "db",
			Protocol:       types.HTTP.String(),
			Authentication: types.NoneAuthentication.String(),
			Authorization:  types.NoneAuthorization.String(),
			Usage:          types.Business.String(),
			DataAssetsSent: []string{"secret"},
		},
	}
	db := topologyTechnicalAsset("db", types.Datastore.String())
	db.DataAssetsStored = []string{"secret"} // raises its availability rating to mission-critical
	lonely := topologyTechnicalAsset("lonely", types.Process.String())
	return &input.Model{
		BusinessCriticality: types.Archive.String(),
		DataAssets: map[string]input.DataAsset{
			"secret": {
				ID:              "secret",
				Usage:           types.Business.String(),
				Quantity:        types.Few.String(),
				Confidentiality: types.StrictlyConfidential.String(),
				Integrity:       types.MissionCritical.String(),
				Availability:    types.MissionCritical.String(),
			},
		},
		TechnicalAssets: map[string]input.TechnicalAsset{
			"web":    web,
			"db":     db,
			"lonely": lonely,
		},
		TrustBoundaries: trustBoundaries,
	}
}

func nestedCrossNetworkTrustBoundaries() map[string]input.TrustBoundary {
	return map[string]input.TrustBoundary{
		"dmz": topologyTrustBoundary("dmz", types.NetworkCloudProvider, nil, []string{"web-exec"}),
		"web-exec": topologyTrustBoundary("web-exec", types.ExecutionEnvironment,
			[]string{"web"}, nil),
		"internal": topologyTrustBoundary("internal", types.NetworkOnPrem, nil, []string{"db-namespace"}),
		"db-namespace": topologyTrustBoundary("db-namespace", types.ExecutionEnvironment,
			nil, []string{"db-pod"}),
		"db-pod": topologyTrustBoundary("db-pod", types.ExecutionEnvironment,
			[]string{"db"}, nil),
	}
}

func TestParseModel_TrustBoundaryTopology_NestedExecutionEnvironmentsResolveNetworkDomain(t *testing.T) {
	parsedModel, err := ParseModel(&mockConfig{}, createNestedTrustBoundaryInputModel(nestedCrossNetworkTrustBoundaries()),
		make(types.RiskRules), make(types.RiskRules))

	assert.NoError(t, err)

	// validated parent mapping
	assert.Equal(t, "dmz", parsedModel.FindParentTrustBoundary(parsedModel.TrustBoundaries["web-exec"]).Id)
	assert.Equal(t, "db-namespace", parsedModel.FindParentTrustBoundary(parsedModel.TrustBoundaries["db-pod"]).Id)
	assert.Nil(t, parsedModel.FindParentTrustBoundary(parsedModel.TrustBoundaries["dmz"]))

	// complete ancestor chains (direct parent first, outermost root last)
	assert.Equal(t, []string{"db-namespace", "internal"},
		trustBoundaryIds(parsedModel.TrustBoundaryAncestors(parsedModel.TrustBoundaries["db-pod"])))
	assert.Equal(t, []string{"web-exec", "dmz"},
		parsedModel.AllParentTrustBoundaryIDs(parsedModel.TrustBoundaries["web-exec"]))

	// the containment facts used by the diagram must match the validated topology
	assert.Contains(t, parsedModel.RecursivelyAllTechnicalAssetIDsInside(parsedModel.TrustBoundaries["internal"]), "db")
	assert.Contains(t, parsedModel.RecursivelyAllTechnicalAssetIDsInside(parsedModel.TrustBoundaries["dmz"]), "web")
}

func TestParseModel_TrustBoundaryTopology_LineageAndNearestNetwork(t *testing.T) {
	parsedModel, err := ParseModel(&mockConfig{}, createNestedTrustBoundaryInputModel(nestedCrossNetworkTrustBoundaries()),
		make(types.RiskRules), make(types.RiskRules))

	assert.NoError(t, err)

	// complete containment chain of an asset (direct boundary first, outermost root last)
	assert.Equal(t, []string{"db-pod", "db-namespace", "internal"},
		technicalAssetTrustBoundaryIds(parsedModel, "db"))
	assert.Equal(t, []string{"web-exec", "dmz"},
		technicalAssetTrustBoundaryIds(parsedModel, "web"))

	// nearest network domain is resolved through arbitrary nesting depth
	assert.Equal(t, "internal", parsedModel.NearestNetworkTrustBoundary("db").Id)
	assert.Equal(t, "dmz", parsedModel.NearestNetworkTrustBoundary("web").Id)

	// legitimate boundary-less asset: no lineage, no network domain, but no error
	assert.Empty(t, parsedModel.TrustBoundaryLineage("lonely"))
	assert.Nil(t, parsedModel.NearestNetworkTrustBoundary("lonely"))
}

func TestParseModel_TrustBoundaryTopology_InvalidNestingRejected(t *testing.T) {
	tests := []struct {
		name            string
		trustBoundaries map[string]input.TrustBoundary
		errorContains   string
	}{
		{
			name: "cyclic nesting",
			trustBoundaries: map[string]input.TrustBoundary{
				"dmz":      topologyTrustBoundary("dmz", types.NetworkCloudProvider, nil, []string{"web-exec"}),
				"web-exec": topologyTrustBoundary("web-exec", types.ExecutionEnvironment, []string{"web"}, []string{"dmz"}),
			},
			errorContains: "cyclic",
		},
		{
			name: "self nesting",
			trustBoundaries: map[string]input.TrustBoundary{
				"dmz": topologyTrustBoundary("dmz", types.NetworkCloudProvider, []string{"web"}, []string{"dmz"}),
			},
			errorContains: "cyclic",
		},
		{
			name: "multiple parent boundaries",
			trustBoundaries: map[string]input.TrustBoundary{
				"dmz":         topologyTrustBoundary("dmz", types.NetworkCloudProvider, nil, []string{"shared-exec"}),
				"other":       topologyTrustBoundary("other", types.NetworkOnPrem, nil, []string{"shared-exec"}),
				"shared-exec": topologyTrustBoundary("shared-exec", types.ExecutionEnvironment, []string{"web"}, nil),
			},
			errorContains: "multiple parent trust boundaries",
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := ParseModel(&mockConfig{}, createNestedTrustBoundaryInputModel(testCase.trustBoundaries),
				make(types.RiskRules), make(types.RiskRules))
			assert.Error(t, err)
			assert.Contains(t, err.Error(), testCase.errorContains)
		})
	}
}

func TestParseModel_TrustBoundaryTopology_AssetsWithoutBoundaryAllowed(t *testing.T) {
	parsedModel, err := ParseModel(&mockConfig{},
		createNestedTrustBoundaryInputModel(make(map[string]input.TrustBoundary)),
		make(types.RiskRules), make(types.RiskRules))

	assert.NoError(t, err)
	assert.NotNil(t, parsedModel.ParentTrustBoundaryMappedByTrustBoundaryId)
	assert.NotNil(t, parsedModel.TrustBoundaryAncestorsMappedByTrustBoundaryId)
	assert.NotNil(t, parsedModel.NearestNetworkTrustBoundaryMappedByTechnicalAssetId)
	assert.Empty(t, parsedModel.NearestNetworkTrustBoundaryMappedByTechnicalAssetId)
	assert.Nil(t, parsedModel.NearestNetworkTrustBoundary("web"))
	assert.Nil(t, parsedModel.NearestNetworkTrustBoundary("db"))
}

// =========================================================================
// network-risk invariance when adding pure logical nesting
// =========================================================================

func flatCrossNetworkTrustBoundaries() map[string]input.TrustBoundary {
	return map[string]input.TrustBoundary{
		"dmz":      topologyTrustBoundary("dmz", types.NetworkCloudProvider, []string{"web"}, nil),
		"internal": topologyTrustBoundary("internal", types.NetworkOnPrem, []string{"db"}, nil),
	}
}

func flatSameNetworkTrustBoundaries() map[string]input.TrustBoundary {
	return map[string]input.TrustBoundary{
		"dmz": topologyTrustBoundary("dmz", types.NetworkCloudProvider, []string{"web", "db"}, nil),
	}
}

func nestedSameNetworkTrustBoundaries() map[string]input.TrustBoundary {
	return map[string]input.TrustBoundary{
		"dmz":      topologyTrustBoundary("dmz", types.NetworkCloudProvider, nil, []string{"web-exec", "db-exec"}),
		"web-exec": topologyTrustBoundary("web-exec", types.ExecutionEnvironment, []string{"web"}, nil),
		"db-exec":  topologyTrustBoundary("db-exec", types.ExecutionEnvironment, []string{"db"}, nil),
	}
}

type networkRiskFingerprint struct {
	category   string
	synthetic  string
	severity   types.RiskSeverity
	likelihood types.RiskExploitationLikelihood
}

func networkRiskFingerprints(parsedModel *types.Model) []networkRiskFingerprint {
	rules := []types.RiskRule{
		builtin.NewUnencryptedCommunicationRule(),
		builtin.NewDosRiskyAccessAcrossTrustBoundaryRule(),
	}
	fingerprints := make([]networkRiskFingerprint, 0)
	for _, rule := range rules {
		risks, err := rule.GenerateRisks(parsedModel)
		if err != nil {
			panic(err)
		}
		for _, risk := range risks {
			fingerprints = append(fingerprints, networkRiskFingerprint{
				category:   risk.CategoryId,
				synthetic:  risk.SyntheticId,
				severity:   risk.Severity,
				likelihood: risk.ExploitationLikelihood,
			})
		}
	}
	sort.Slice(fingerprints, func(i, j int) bool {
		return fingerprints[i].synthetic < fingerprints[j].synthetic
	})
	return fingerprints
}

func parseTopologyScenario(t *testing.T, trustBoundaries map[string]input.TrustBoundary) *types.Model {
	parsedModel, err := ParseModel(&mockConfig{}, createNestedTrustBoundaryInputModel(trustBoundaries),
		make(types.RiskRules), make(types.RiskRules))
	assert.NoError(t, err)
	return parsedModel
}

func TestNetworkRisks_PureLogicalNestingDoesNotChangeNetworkRisks(t *testing.T) {
	flatCross := parseTopologyScenario(t, flatCrossNetworkTrustBoundaries())
	nestedCross := parseTopologyScenario(t, nestedCrossNetworkTrustBoundaries())
	flatSame := parseTopologyScenario(t, flatSameNetworkTrustBoundaries())
	nestedSame := parseTopologyScenario(t, nestedSameNetworkTrustBoundaries())

	// adding execution-environment nesting must not change the generated network risks
	assert.Equal(t, networkRiskFingerprints(flatCross), networkRiskFingerprints(nestedCross))
	assert.Equal(t, networkRiskFingerprints(flatSame), networkRiskFingerprints(nestedSame))

	// and the two network situations must still be distinguishable: a real cross-network link
	// stays detected (elevated likelihood for the encryption rule and the DoS rule fires), ...
	crossFingerprints := networkRiskFingerprints(nestedCross)
	assert.Contains(t, crossFingerprints, networkRiskFingerprint{
		category:   "unencrypted-communication",
		synthetic:  "unencrypted-communication@web>call-db@web@db",
		severity:   types.CalculateSeverity(types.Likely, types.HighImpact),
		likelihood: types.Likely,
	})
	assertHasRiskCategory(t, crossFingerprints, "dos-risky-access-across-trust-boundary")

	// ... while same-network links keep their lower, non-crossing risk profile
	sameFingerprints := networkRiskFingerprints(nestedSame)
	assert.Contains(t, sameFingerprints, networkRiskFingerprint{
		category:   "unencrypted-communication",
		synthetic:  "unencrypted-communication@web>call-db@web@db",
		severity:   types.CalculateSeverity(types.Unlikely, types.HighImpact),
		likelihood: types.Unlikely,
	})
	assertNotHasRiskCategory(t, sameFingerprints, "dos-risky-access-across-trust-boundary")
}

func assertHasRiskCategory(t *testing.T, fingerprints []networkRiskFingerprint, category string) {
	t.Helper()
	for _, fingerprint := range fingerprints {
		if fingerprint.category == category {
			return
		}
	}
	assert.Fail(t, "expected risk category "+category)
}

func assertNotHasRiskCategory(t *testing.T, fingerprints []networkRiskFingerprint, category string) {
	t.Helper()
	for _, fingerprint := range fingerprints {
		if fingerprint.category == category {
			assert.Fail(t, "did not expect risk category "+category)
			return
		}
	}
}

func trustBoundaryIds(boundaries []*types.TrustBoundary) []string {
	result := make([]string, len(boundaries))
	for i, boundary := range boundaries {
		result[i] = boundary.Id
	}
	return result
}

func technicalAssetTrustBoundaryIds(parsedModel *types.Model, technicalAssetId string) []string {
	return trustBoundaryIds(parsedModel.TrustBoundaryLineage(technicalAssetId))
}
