package gql

import (
	"checkpoint/gql/enum"
	"checkpoint/modules/authentication"
	imeiconfiguration "checkpoint/modules/imei-configuration"
	imsiconfiguration "checkpoint/modules/imsi-configuration"
	mobiledeviceconfiguration "checkpoint/modules/mobile-device-configuration"
	"checkpoint/modules/project"
	projectRole "checkpoint/modules/project-role"
	stationdashboardactivity "checkpoint/modules/station-dashboard-activity"
	note "checkpoint/modules/station-dashboard-report"
	stationDevice "checkpoint/modules/station-device"
	stationDeviceHealthCheck "checkpoint/modules/station-device-health-check"
	stationHealthCheck "checkpoint/modules/station-health-check"
	stationimeiimsiactivity "checkpoint/modules/station-imei-imsi-activity"
	stationlocation "checkpoint/modules/station-location"
	stationOfficer "checkpoint/modules/station-officer"
	vehicleActivity "checkpoint/modules/station-vehicle-activity"
	"checkpoint/modules/tag"
	"checkpoint/modules/upload"
	vehicleLicensePlate "checkpoint/modules/vehicle-license-plate"
	vehicleproperty "checkpoint/modules/vehicle-property"
	vehicleTarget "checkpoint/modules/vehicle-target-configuration"
	vehicletargetconfigurationimage "checkpoint/modules/vehicle-target-configuration-image"
)

type Resolver struct {
	BlacklistPriority    enum.BlacklistPriority
	DevicePermittedLabel enum.DevicePermittedLabel
	ImageType            enum.ImageType
	DeviceStatus         enum.DeviceStatus
	StationStatus        enum.StationStatus
	PermissionAction     enum.PermissionAction
	project.ProjectResolver
	authentication.AuthenticationResolver
	projectRole.ProjectRoleResolver
	imsiconfiguration.ImsiConfigurationResolver
	imeiconfiguration.ImeiConfigurationResolver
	mobiledeviceconfiguration.MobileDeviceConfigurationResolver
	tag.TagResolver
	note.NoteResolver
	stationlocation.StationLocationResolver
	stationOfficer.StationOfficerResolver
	stationDevice.StationDeviceResolver
	stationDeviceHealthCheck.StationDeviceHealthCheckActivityResolver
	stationHealthCheck.StationLocationHealthCheckActivityResolver
	vehicleTarget.VehicleTargetConfigurationResolver
	vehicleActivity.StationVehicleActivityResolver
	vehicleproperty.VehiclePropertyResolver
	vehicleLicensePlate.VehiclelicensePlateResolver
	upload.UploadResolver
	stationimeiimsiactivity.StationImeiImsiActivityResolver
	stationdashboardactivity.StationDashboardActivityResolver
	vehicletargetconfigurationimage.VehicleTargetConfigurationImageResolver
}

func GraphqlResolver() *Resolver {
	r := &Resolver{}

	r.StationVehicleActivityResolver.SetupSubscription()

	go r.StationVehicleActivityResolver.BroadcastStationVehicleActivity()

	return r
}
