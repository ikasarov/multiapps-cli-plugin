package mtaclient_v2

import (
	models "github.com/cloudfoundry-incubator/multiapps-cli-plugin/clients/models"
	strfmt "github.com/go-openapi/strfmt"
)

// MtaClientOperations drun drun drun
type MtaV2ClientOperations interface {
	GetMtas(name, namespace, spaceGuid string) ([]*models.Mta, error)
}

// ResponseHeader response header
type ResponseHeader struct {
	Location strfmt.URI
}
