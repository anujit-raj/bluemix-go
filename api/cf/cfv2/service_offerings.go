package cfv2

import (
	"fmt"

	"github.com/IBM-Bluemix/bluemix-go/client"

	"github.com/IBM-Bluemix/bluemix-go/bmxerror"
	"github.com/IBM-Bluemix/bluemix-go/rest"
)

//ErrCodeServiceDoesnotExist ...
const ErrCodeServiceDoesnotExist = "ServiceDoesnotExist"

//ServiceOffering model
type ServiceOffering struct {
	GUID              string
	Label             string   `json:"label"`
	Provider          string   `json:"provider"`
	Description       string   `json:"description"`
	LongDescription   string   `json:"long_description"`
	Version           string   `json:"version"`
	URL               string   `json:"url"`
	InfoURL           string   `json:"info_url"`
	DocumentURL       string   `json:"documentation_url"`
	Timeout           string   `json:"timeout"`
	UniqueID          string   `json:"unique_id"`
	ServiceBrokerGUID string   `json:"service_broker_guid"`
	ServicePlansURL   string   `json:"service_plans_url"`
	Tags              []string `json:"tags"`
	Requires          []string `json:"requires"`
	IsActive          bool     `json:"active"`
	IsBindable        bool     `json:"bindable"`
	IsPlanUpdateable  bool     `json:"plan_updateable"`
}

//ServiceOfferingResource ...
type ServiceOfferingResource struct {
	Resource
	Entity ServiceOfferingEntity
}

//ServiceOfferingEntity ...
type ServiceOfferingEntity struct {
	Label             string   `json:"label"`
	Provider          string   `json:"provider"`
	Description       string   `json:"description"`
	LongDescription   string   `json:"long_description"`
	Version           string   `json:"version"`
	URL               string   `json:"url"`
	InfoURL           string   `json:"info_url"`
	DocumentURL       string   `json:"documentation_url"`
	Timeout           string   `json:"timeout"`
	UniqueID          string   `json:"unique_id"`
	ServiceBrokerGUID string   `json:"service_broker_guid"`
	ServicePlansURL   string   `json:"service_plans_url"`
	Tags              []string `json:"tags"`
	Requires          []string `json:"requires"`
	IsActive          bool     `json:"active"`
	IsBindable        bool     `json:"bindable"`
	IsPlanUpdateable  bool     `json:"plan_updateable"`
}

//ToFields ...
func (resource ServiceOfferingResource) ToFields() ServiceOffering {
	entity := resource.Entity

	return ServiceOffering{
		GUID:              resource.Metadata.GUID,
		Label:             entity.Label,
		Provider:          entity.Provider,
		Description:       entity.Description,
		LongDescription:   entity.LongDescription,
		Version:           entity.Version,
		URL:               entity.URL,
		InfoURL:           entity.InfoURL,
		DocumentURL:       entity.DocumentURL,
		Timeout:           entity.Timeout,
		UniqueID:          entity.UniqueID,
		ServiceBrokerGUID: entity.ServiceBrokerGUID,
		ServicePlansURL:   entity.ServicePlansURL,
		Tags:              entity.Tags,
		Requires:          entity.Requires,
		IsActive:          entity.IsActive,
		IsBindable:        entity.IsBindable,
		IsPlanUpdateable:  entity.IsPlanUpdateable,
	}
}

//ServiceOfferings ...
type ServiceOfferings interface {
	FindByLabel(serviceName string) (*ServiceOffering, error)
}

type serviceOfferrings struct {
	client *client.Client
}

func newServiceOfferingAPI(c *client.Client) ServiceOfferings {
	return &serviceOfferrings{
		client: c,
	}
}

func (r *serviceOfferrings) FindByLabel(serviceName string) (*ServiceOffering, error) {
	req := rest.GetRequest("v2/services")
	if serviceName != "" {
		req.Query("q", "label:"+serviceName)
	}
	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()
	var services ServiceOffering
	var found bool
	err = r.listServicesOfferingWithPath(path, func(serviceOfferingResource ServiceOfferingResource) bool {
		services = serviceOfferingResource.ToFields()
		found = true
		return false
	})

	if err != nil {
		return nil, err
	}

	if found {
		return &services, err
	}
	//May not be found and no error

	return nil, bmxerror.New(ErrCodeServiceDoesnotExist,
		fmt.Sprintf("Given service %q doesn't exist", serviceName))

}

func (r *serviceOfferrings) listServicesOfferingWithPath(path string, cb func(ServiceOfferingResource) bool) error {
	_, err := r.client.GetPaginated(path, ServiceOfferingResource{}, func(resource interface{}) bool {
		if serviceOfferingResource, ok := resource.(ServiceOfferingResource); ok {
			return cb(serviceOfferingResource)
		}
		return false
	})
	return err
}
