package gql

import (
	authentication "checkpoint/modules/authentication"
	imeiconfiguration "checkpoint/modules/imei-configuration"
	imsiconfiguration "checkpoint/modules/imsi-configuration"
	project "checkpoint/modules/project"
	projectRole "checkpoint/modules/project-role"
)

type Resolver struct {
	project.ProjectResolver
	authentication.AuthenticationResolver
	projectRole.ProjectRoleResolver
	imsiconfiguration.ImsiConfigurationResolver
	imeiconfiguration.ImeiConfigurationResolver
}
