package gql

import (
	"checkpoint/gql/enum"
	"checkpoint/modules/authentication"
	imeiconfiguration "checkpoint/modules/imei-configuration"
	imsiconfiguration "checkpoint/modules/imsi-configuration"
	mobiledeviceconfiguration "checkpoint/modules/mobile-device-configuration"
	"checkpoint/modules/notification"
	"checkpoint/modules/project"
	projectRole "checkpoint/modules/project-role"
	stationDevice "checkpoint/modules/station-device"
	stationlocation "checkpoint/modules/station-location"
	stationOfficer "checkpoint/modules/station-officer"
	"checkpoint/modules/tag"
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
	notification.NotificationResolver
	stationlocation.StationLocationResolver
	stationOfficer.StationOfficerResolver
	stationDevice.StationDeviceResolver
}

func GraphqlResolver() *Resolver {
	r := &Resolver{}
	r.NotificationResolver.WsResolver()

	return r
}
