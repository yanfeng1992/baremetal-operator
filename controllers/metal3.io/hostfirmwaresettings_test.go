package controllers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	metal3v1alpha1 "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	"github.com/metal3-io/baremetal-operator/pkg/provisioner"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	hostName      string = "myHostName"
	hostNamespace string = "myHostNamespace"
	schemaName    string = "schema-4bcc035f" // Hash generated from schema, change this if the schema is changed
)

var (
	iTrue      bool = true
	iFalse     bool = false
	minLength  int  = 0
	maxLength  int  = 20
	lowerBound int  = 0
	upperBound int  = 20
)

// Test support for HostFirmwareSettings in the HostFirmwareSettingsReconciler
func getTestHFSReconciler(host *metal3v1alpha1.HostFirmwareSettings) *HostFirmwareSettingsReconciler {

	c := fakeclient.NewFakeClient(host)
	reconciler := &HostFirmwareSettingsReconciler{
		Client: c,
		Log:    ctrl.Log.WithName("test_reconciler").WithName("HostFirmwareSettings"),
	}

	return reconciler
}

func getMockProvisioner(settings metal3v1alpha1.SettingsMap, schema map[string]metal3v1alpha1.SettingSchema) *hsfMockProvisioner {
	return &hsfMockProvisioner{
		Settings: settings,
		Schema:   schema,
		Error:    nil,
	}
}

type hsfMockProvisioner struct {
	Settings metal3v1alpha1.SettingsMap
	Schema   map[string]metal3v1alpha1.SettingSchema
	Error    error
}

func (m *hsfMockProvisioner) HasCapacity() (result bool, err error) {
	return
}

func (m *hsfMockProvisioner) ValidateManagementAccess(data provisioner.ManagementAccessData, credentialsChanged, force bool) (result provisioner.Result, provID string, err error) {
	return
}

func (m *hsfMockProvisioner) InspectHardware(data provisioner.InspectData, force, refresh bool) (result provisioner.Result, started bool, details *metal3v1alpha1.HardwareDetails, err error) {
	return
}

func (m *hsfMockProvisioner) UpdateHardwareState() (hwState provisioner.HardwareState, err error) {
	return
}

func (m *hsfMockProvisioner) Prepare(data provisioner.PrepareData, unprepared bool) (result provisioner.Result, started bool, err error) {
	return
}

func (m *hsfMockProvisioner) Adopt(data provisioner.AdoptData, force bool) (result provisioner.Result, err error) {
	return
}

func (m *hsfMockProvisioner) Provision(data provisioner.ProvisionData) (result provisioner.Result, err error) {
	return
}

func (m *hsfMockProvisioner) Deprovision(force bool) (result provisioner.Result, err error) {
	return
}

func (m *hsfMockProvisioner) Delete() (result provisioner.Result, err error) {
	return
}

func (m *hsfMockProvisioner) Detach() (result provisioner.Result, err error) {
	return
}

func (m *hsfMockProvisioner) PowerOn(force bool) (result provisioner.Result, err error) {
	return
}

func (m *hsfMockProvisioner) PowerOff(rebootMode metal3v1alpha1.RebootMode, force bool) (result provisioner.Result, err error) {
	return
}

func (m *hsfMockProvisioner) IsReady() (result bool, err error) {
	return
}

func (m *hsfMockProvisioner) GetFirmwareSettings(includeSchema bool) (settings metal3v1alpha1.SettingsMap, schema map[string]metal3v1alpha1.SettingSchema, err error) {

	return m.Settings, m.Schema, m.Error
}

func getSchema() *metal3v1alpha1.FirmwareSchema {

	schema := &metal3v1alpha1.FirmwareSchema{
		TypeMeta: metav1.TypeMeta{
			Kind:       "FirmwareSchema",
			APIVersion: "metal3.io/v1alpha1"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      schemaName,
			Namespace: hostNamespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "metal3.io/v1alpha1",
					Kind:       "HostFirmwareSettings",
					Name:       "dummyhfs",
				},
			},
		},
	}

	return schema
}

// Mock settings to return from provisioner
func getCurrentSettings() metal3v1alpha1.SettingsMap {

	return metal3v1alpha1.SettingsMap{
		"L2Cache":               "10x512 KB",
		"NetworkBootRetryCount": "20",
		"ProcVirtualization":    "Disabled",
		"SecureBoot":            "Enabled",
		"AssetTag":              "X45672917",
	}
}

// Mock schema to return from provisioner
func getCurrentSchemaSettings() map[string]metal3v1alpha1.SettingSchema {

	return map[string]metal3v1alpha1.SettingSchema{
		"AssetTag": {
			AttributeType: "String",
			MinLength:     &minLength,
			MaxLength:     &maxLength,
			Unique:        &iTrue,
		},
		"CustomPostMessage": {
			AttributeType: "String",
			MinLength:     &minLength,
			MaxLength:     &maxLength,
			Unique:        &iFalse,
			ReadOnly:      &iFalse,
		},
		"L2Cache": {
			AttributeType: "String",
			MinLength:     &minLength,
			MaxLength:     &maxLength,
			ReadOnly:      &iTrue,
		},
		"NetworkBootRetryCount": {
			AttributeType: "Integer",
			LowerBound:    &lowerBound,
			UpperBound:    &upperBound,
			ReadOnly:      &iFalse,
		},
		"ProcVirtualization": {
			AttributeType:   "Enumeration",
			AllowableValues: []string{"Enabled", "Disabled"},
			ReadOnly:        &iFalse,
		},
		"SecureBoot": {
			AttributeType:   "Enumeration",
			AllowableValues: []string{"Enabled", "Disabled"},
			ReadOnly:        &iTrue,
		},
	}
}

// Create the baremetalhost reconciler and use that to create bmh in same namespace
func createBaremetalHost() *metal3v1alpha1.BareMetalHost {

	bmh := &metal3v1alpha1.BareMetalHost{}
	bmh.ObjectMeta = metav1.ObjectMeta{Name: hostName, Namespace: hostNamespace}
	c := fakeclient.NewFakeClient(bmh)

	reconciler := &BareMetalHostReconciler{
		Client:             c,
		ProvisionerFactory: nil,
		Log:                ctrl.Log.WithName("bmh_reconciler").WithName("BareMetalHost"),
	}

	reconciler.Create(context.TODO(), bmh)

	return bmh
}

func getExpectedSchema() *metal3v1alpha1.FirmwareSchema {
	firmwareSchema := getSchema()
	firmwareSchema.ObjectMeta.ResourceVersion = "1"
	firmwareSchema.ObjectMeta.OwnerReferences = []metav1.OwnerReference{
		{
			APIVersion: "metal3.io/v1alpha1",
			Kind:       "HostFirmwareSettings",
			Name:       hostName,
		},
	}
	firmwareSchema.Spec.Schema = getCurrentSchemaSettings()

	return firmwareSchema
}

func getExpectedSchemaTwoOwners() *metal3v1alpha1.FirmwareSchema {
	firmwareSchema := getSchema()
	firmwareSchema.ObjectMeta.ResourceVersion = "2"
	firmwareSchema.ObjectMeta.OwnerReferences = []metav1.OwnerReference{
		{
			APIVersion: "metal3.io/v1alpha1",
			Kind:       "HostFirmwareSettings",
			Name:       "dummyhfs",
		},
		{
			APIVersion: "metal3.io/v1alpha1",
			Kind:       "HostFirmwareSettings",
			Name:       hostName,
		},
	}
	firmwareSchema.Spec.Schema = getCurrentSchemaSettings()

	return firmwareSchema
}

// Create an HFS with input spec settings
func getHFS(spec metal3v1alpha1.HostFirmwareSettingsSpec) *metal3v1alpha1.HostFirmwareSettings {

	hfs := &metal3v1alpha1.HostFirmwareSettings{}

	hfs.Status = metal3v1alpha1.HostFirmwareSettingsStatus{
		Settings: metal3v1alpha1.SettingsMap{
			"CustomPostMessage":     "All tests passed",
			"L2Cache":               "10x256 KB",
			"NetworkBootRetryCount": "10",
			"ProcVirtualization":    "Enabled",
			"SecureBoot":            "Enabled",
			"AssetTag":              "X45672917",
		},
	}
	hfs.TypeMeta = metav1.TypeMeta{
		Kind:       "HostFirmwareSettings",
		APIVersion: "metal3.io/v1alpha1"}
	hfs.ObjectMeta = metav1.ObjectMeta{
		Name:      hostName,
		Namespace: hostNamespace}

	hfs.Spec = spec

	return hfs
}

// Test the hostfirmwaresettings reconciler functions
func TestStoreHostFirmwareSettings(t *testing.T) {

	testCases := []struct {
		Scenario string
		// the resource that the reconciler is managing
		CurrentHFSResource *metal3v1alpha1.HostFirmwareSettings
		// whether to create a schema resource before calling reconciler
		CreateSchemaResource bool
		// the expected created or updated resource
		ExpectedSettings *metal3v1alpha1.HostFirmwareSettings
		// whether the spec values pass the validity test
		SpecIsValid bool
	}{
		{
			Scenario: "initial hfs resource with no schema",
			CurrentHFSResource: &metal3v1alpha1.HostFirmwareSettings{
				TypeMeta: metav1.TypeMeta{
					Kind:       "HostFirmwareSettings",
					APIVersion: "metal3.io/v1alpha1"},
				ObjectMeta: metav1.ObjectMeta{
					Name:            hostName,
					Namespace:       hostNamespace,
					ResourceVersion: "1"},
				Spec: metal3v1alpha1.HostFirmwareSettingsSpec{
					Settings: metal3v1alpha1.DesiredSettingsMap{},
				},
				Status: metal3v1alpha1.HostFirmwareSettingsStatus{},
			},
			CreateSchemaResource: false,
			ExpectedSettings: &metal3v1alpha1.HostFirmwareSettings{
				Spec: metal3v1alpha1.HostFirmwareSettingsSpec{
					Settings: metal3v1alpha1.DesiredSettingsMap{},
				},
				Status: metal3v1alpha1.HostFirmwareSettingsStatus{
					FirmwareSchema: &metal3v1alpha1.SchemaReference{
						Name:      schemaName,
						Namespace: hostNamespace,
					},
					Settings: metal3v1alpha1.SettingsMap{
						"AssetTag":              "X45672917",
						"L2Cache":               "10x512 KB",
						"NetworkBootRetryCount": "20",
						"ProcVirtualization":    "Disabled",
						"SecureBoot":            "Enabled",
					},
					Conditions: []metav1.Condition{
						{Type: "Valid", Status: "True", Reason: "Success"},
					},
				},
			},
			SpecIsValid: true,
		},
		{
			Scenario: "initial hfs resource with existing schema",
			CurrentHFSResource: &metal3v1alpha1.HostFirmwareSettings{
				TypeMeta: metav1.TypeMeta{
					Kind:       "HostFirmwareSettings",
					APIVersion: "metal3.io/v1alpha1"},
				ObjectMeta: metav1.ObjectMeta{
					Name:            hostName,
					Namespace:       hostNamespace,
					ResourceVersion: "1"},
				Spec: metal3v1alpha1.HostFirmwareSettingsSpec{
					Settings: metal3v1alpha1.DesiredSettingsMap{},
				},
				Status: metal3v1alpha1.HostFirmwareSettingsStatus{},
			},
			CreateSchemaResource: true,
			ExpectedSettings: &metal3v1alpha1.HostFirmwareSettings{
				Spec: metal3v1alpha1.HostFirmwareSettingsSpec{
					Settings: metal3v1alpha1.DesiredSettingsMap{},
				},
				Status: metal3v1alpha1.HostFirmwareSettingsStatus{
					FirmwareSchema: &metal3v1alpha1.SchemaReference{
						Name:      schemaName,
						Namespace: hostNamespace,
					},
					Settings: metal3v1alpha1.SettingsMap{
						"AssetTag":              "X45672917",
						"L2Cache":               "10x512 KB",
						"NetworkBootRetryCount": "20",
						"ProcVirtualization":    "Disabled",
						"SecureBoot":            "Enabled",
					},
					Conditions: []metav1.Condition{
						{Type: "Valid", Status: "True", Reason: "Success"},
					},
				},
			},
			SpecIsValid: true,
		},
		{
			Scenario: "updated settings",
			CurrentHFSResource: &metal3v1alpha1.HostFirmwareSettings{
				TypeMeta: metav1.TypeMeta{
					Kind:       "HostFirmwareSettings",
					APIVersion: "metal3.io/v1alpha1"},
				ObjectMeta: metav1.ObjectMeta{
					Name:            hostName,
					Namespace:       hostNamespace,
					ResourceVersion: "1"},
				Spec: metal3v1alpha1.HostFirmwareSettingsSpec{
					Settings: metal3v1alpha1.DesiredSettingsMap{
						"NetworkBootRetryCount": intstr.FromString("10"),
						"ProcVirtualization":    intstr.FromString("Enabled"),
						"AssetTag":              intstr.FromString("Z98765432"),
					},
				},
				Status: metal3v1alpha1.HostFirmwareSettingsStatus{
					FirmwareSchema: &metal3v1alpha1.SchemaReference{
						Name:      schemaName,
						Namespace: hostNamespace,
					},
					Settings: metal3v1alpha1.SettingsMap{
						"AssetTag":              "Z98765432",
						"L2Cache":               "10x256 KB",
						"NetworkBootRetryCount": "10",
						"ProcVirtualization":    "Enabled",
					},
				},
			},
			CreateSchemaResource: true,
			ExpectedSettings: &metal3v1alpha1.HostFirmwareSettings{
				Spec: metal3v1alpha1.HostFirmwareSettingsSpec{
					Settings: metal3v1alpha1.DesiredSettingsMap{
						"NetworkBootRetryCount": intstr.FromString("10"),
						"ProcVirtualization":    intstr.FromString("Enabled"),
						"AssetTag":              intstr.FromString("Z98765432"),
					},
				},
				Status: metal3v1alpha1.HostFirmwareSettingsStatus{
					FirmwareSchema: &metal3v1alpha1.SchemaReference{
						Name:      schemaName,
						Namespace: hostNamespace,
					},
					Settings: metal3v1alpha1.SettingsMap{
						"AssetTag":              "X45672917",
						"L2Cache":               "10x512 KB",
						"NetworkBootRetryCount": "20",
						"ProcVirtualization":    "Disabled",
						"SecureBoot":            "Enabled",
					},
					Conditions: []metav1.Condition{
						{Type: "ChangeDetected", Status: "True", Reason: "Success"},
						{Type: "Valid", Status: "True", Reason: "Success"},
					},
				},
			},
			SpecIsValid: true,
		},
		{
			Scenario: "spec updated with invalid setting",
			CurrentHFSResource: &metal3v1alpha1.HostFirmwareSettings{
				TypeMeta: metav1.TypeMeta{
					Kind:       "HostFirmwareSettings",
					APIVersion: "metal3.io/v1alpha1"},
				ObjectMeta: metav1.ObjectMeta{
					Name:            hostName,
					Namespace:       hostNamespace,
					ResourceVersion: "1"},
				Spec: metal3v1alpha1.HostFirmwareSettingsSpec{
					Settings: metal3v1alpha1.DesiredSettingsMap{
						"NetworkBootRetryCount": intstr.FromString("1000"),
						"ProcVirtualization":    intstr.FromString("Enabled"),
					},
				},
				Status: metal3v1alpha1.HostFirmwareSettingsStatus{
					FirmwareSchema: &metal3v1alpha1.SchemaReference{
						Name:      schemaName,
						Namespace: hostNamespace,
					},
					Settings: metal3v1alpha1.SettingsMap{
						"L2Cache":               "10x256 KB",
						"NetworkBootRetryCount": "10",
						"ProcVirtualization":    "Enabled",
					},
				},
			},
			CreateSchemaResource: true,
			ExpectedSettings: &metal3v1alpha1.HostFirmwareSettings{
				Spec: metal3v1alpha1.HostFirmwareSettingsSpec{
					Settings: metal3v1alpha1.DesiredSettingsMap{
						"NetworkBootRetryCount": intstr.FromString("1000"),
						"ProcVirtualization":    intstr.FromString("Enabled"),
					},
				},
				Status: metal3v1alpha1.HostFirmwareSettingsStatus{
					FirmwareSchema: &metal3v1alpha1.SchemaReference{
						Name:      schemaName,
						Namespace: hostNamespace,
					},
					Settings: metal3v1alpha1.SettingsMap{
						"AssetTag":              "X45672917",
						"L2Cache":               "10x512 KB",
						"NetworkBootRetryCount": "20",
						"ProcVirtualization":    "Disabled",
						"SecureBoot":            "Enabled",
					},
					Conditions: []metav1.Condition{
						{Type: "ChangeDetected", Status: "True", Reason: "Success"},
						{Type: "Valid", Status: "False", Reason: "ConfigurationError", Message: "Invalid BIOS setting"},
					},
				},
			},
			SpecIsValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {

			ctx := context.TODO()
			prov := getMockProvisioner(getCurrentSettings(), getCurrentSchemaSettings())

			tc.ExpectedSettings.TypeMeta = metav1.TypeMeta{
				Kind:       "HostFirmwareSettings",
				APIVersion: "metal3.io/v1alpha1"}
			tc.ExpectedSettings.ObjectMeta = metav1.ObjectMeta{
				Name:            hostName,
				Namespace:       hostNamespace,
				ResourceVersion: "2"}

			hfs := tc.CurrentHFSResource
			r := getTestHFSReconciler(hfs)
			// Create bmh resource needed by hfs reconciler
			bmh := createBaremetalHost()

			info := &rInfo{
				log: logf.Log.WithName("controllers").WithName("HostFirmwareSettings"),
				hfs: tc.CurrentHFSResource,
				bmh: bmh,
			}

			if tc.CreateSchemaResource {
				// Create an existing schema with different hfs owner
				firmwareSchema := getSchema()
				firmwareSchema.Spec.Schema = getCurrentSchemaSettings()

				r.Client.Create(ctx, firmwareSchema)
			}

			currentSettings, schema, err := prov.GetFirmwareSettings(true)
			assert.Equal(t, nil, err)

			err = r.updateHostFirmwareSettings(currentSettings, schema, info)
			assert.Equal(t, nil, err)

			// Check that resources get created or updated
			key := client.ObjectKey{
				Namespace: hfs.ObjectMeta.Namespace, Name: hfs.ObjectMeta.Name}
			actualSettings := &metal3v1alpha1.HostFirmwareSettings{}
			err = r.Client.Get(ctx, key, actualSettings)
			assert.Equal(t, nil, err)

			// Use the same time for expected and actual
			currentTime := metav1.Now()
			tc.ExpectedSettings.Status.LastUpdated = &currentTime
			actualSettings.Status.LastUpdated = &currentTime
			for i := range tc.ExpectedSettings.Status.Conditions {
				tc.ExpectedSettings.Status.Conditions[i].LastTransitionTime = currentTime
				actualSettings.Status.Conditions[i].LastTransitionTime = currentTime
			}
			assert.Equal(t, tc.ExpectedSettings, actualSettings)

			key = client.ObjectKey{
				Namespace: hfs.ObjectMeta.Namespace, Name: schemaName}
			actualSchema := &metal3v1alpha1.FirmwareSchema{}
			err = r.Client.Get(ctx, key, actualSchema)
			assert.Equal(t, nil, err)
			var expectedSchema *metal3v1alpha1.FirmwareSchema
			if tc.CreateSchemaResource {
				expectedSchema = getExpectedSchemaTwoOwners()
			} else {
				expectedSchema = getExpectedSchema()
			}
			assert.Equal(t, expectedSchema, actualSchema)
		})
	}
}

// Test the function to validate hostFirmwareSettings
func TestValidateHostFirmwareSettings(t *testing.T) {

	testCases := []struct {
		Scenario      string
		SpecSettings  metal3v1alpha1.HostFirmwareSettingsSpec
		ExpectedError string
	}{
		{
			Scenario: "valid spec changes with schema",
			SpecSettings: metal3v1alpha1.HostFirmwareSettingsSpec{
				Settings: metal3v1alpha1.DesiredSettingsMap{
					"CustomPostMessage":     intstr.FromString("All tests passed"),
					"ProcVirtualization":    intstr.FromString("Disabled"),
					"NetworkBootRetryCount": intstr.FromString("20"),
				},
			},
			ExpectedError: "",
		},
		{
			Scenario: "invalid string",
			SpecSettings: metal3v1alpha1.HostFirmwareSettingsSpec{
				Settings: metal3v1alpha1.DesiredSettingsMap{
					"CustomPostMessage":     intstr.FromString("A really long POST message"),
					"ProcVirtualization":    intstr.FromString("Disabled"),
					"NetworkBootRetryCount": intstr.FromString("20"),
				},
			},
			ExpectedError: "Setting CustomPostMessage is invalid, string A really long POST message length is above maximum length 20",
		},
		{
			Scenario: "invalid int",
			SpecSettings: metal3v1alpha1.HostFirmwareSettingsSpec{
				Settings: metal3v1alpha1.DesiredSettingsMap{
					"CustomPostMessage":     intstr.FromString("All tests passed"),
					"ProcVirtualization":    intstr.FromString("Disabled"),
					"NetworkBootRetryCount": intstr.FromString("2000"),
				},
			},
			ExpectedError: "Setting NetworkBootRetryCount is invalid, integer 2000 is above maximum value 20",
		},
		{
			Scenario: "invalid enum",
			SpecSettings: metal3v1alpha1.HostFirmwareSettingsSpec{
				Settings: metal3v1alpha1.DesiredSettingsMap{
					"CustomPostMessage":     intstr.FromString("All tests passed"),
					"ProcVirtualization":    intstr.FromString("Not enabled"),
					"NetworkBootRetryCount": intstr.FromString("20"),
				},
			},
			ExpectedError: "Setting ProcVirtualization is invalid, unknown enumeration value - Not enabled",
		},
		{
			Scenario: "invalid name",
			SpecSettings: metal3v1alpha1.HostFirmwareSettingsSpec{
				Settings: metal3v1alpha1.DesiredSettingsMap{
					"SomeNewSetting": intstr.FromString("foo"),
				},
			},
			ExpectedError: "Setting SomeNewSetting is not in the Status field",
		},
		{
			Scenario: "invalid password in spec",
			SpecSettings: metal3v1alpha1.HostFirmwareSettingsSpec{
				Settings: metal3v1alpha1.DesiredSettingsMap{
					"CustomPostMessage":     intstr.FromString("All tests passed"),
					"ProcVirtualization":    intstr.FromString("Disabled"),
					"NetworkBootRetryCount": intstr.FromString("20"),
					"SysPassword":           intstr.FromString("Pa%$word"),
				},
			},
			ExpectedError: "Cannot set Password field",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {

			hfs := getHFS(tc.SpecSettings)
			r := getTestHFSReconciler(hfs)
			info := &rInfo{
				log: logf.Log.WithName("controllers").WithName("HostFirmwareSettings"),
				hfs: hfs,
			}

			errors := r.validateHostFirmwareSettings(info, getExpectedSchema())
			if len(errors) == 0 {
				assert.Equal(t, tc.ExpectedError, "")
			} else {
				for _, error := range errors {
					assert.Equal(t, tc.ExpectedError, error.Error())
				}
			}
		})
	}
}
