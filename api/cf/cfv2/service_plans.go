package cfv2

import (
	"fmt"

	"github.com/IBM-Bluemix/bluemix-go/bmxerror"
	"github.com/IBM-Bluemix/bluemix-go/client"
	"github.com/IBM-Bluemix/bluemix-go/rest"
)

//ErrCodeServicePlanDoesNotExist ...
const ErrCodeServicePlanDoesNotExist = "ServicePlanDoesNotExist"

//ServicePlan ...
type ServicePlan struct {
	GUID                string
	Name                string `json:"name"`
	Description         string `json:"description"`
	IsFree              bool   `json:"free"`
	IsPublic            bool   `json:"public"`
	IsActive            bool   `json:"active"`
	ServiceGUID         string `json:"service_guid"`
	UniqueID            string `json:"unique_id"`
	ServiceInstancesURL string `json:"service_instances_url"`
}

//ServicePlanResource ...
type ServicePlanResource struct {
	Resource
	Entity ServicePlanEntity
}

//ServicePlanEntity ...
type ServicePlanEntity struct {
	Name                string `json:"name"`
	Description         string `json:"description"`
	IsFree              bool   `json:"free"`
	IsPublic            bool   `json:"public"`
	IsActive            bool   `json:"active"`
	ServiceGUID         string `json:"service_guid"`
	UniqueID            string `json:"unique_id"`
	ServiceInstancesURL string `json:"service_instances_url"`
}

//ToFields ...
func (resource ServicePlanResource) ToFields() ServicePlan {
	entity := resource.Entity

	return ServicePlan{
		GUID:                resource.Metadata.GUID,
		Name:                entity.Name,
		Description:         entity.Description,
		IsFree:              entity.IsFree,
		IsPublic:            entity.IsPublic,
		IsActive:            entity.IsActive,
		ServiceGUID:         entity.ServiceGUID,
		UniqueID:            entity.UniqueID,
		ServiceInstancesURL: entity.ServiceInstancesURL,
	}
}

//ServicePlans ...
type ServicePlans interface {
	GetServicePlan(serviceOfferingGUID string, planType string) (*ServicePlan, error)
}

type servicePlan struct {
	client *client.Client
}

func newServicePlanAPI(c *client.Client) ServicePlans {
	return &servicePlan{
		client: c,
	}
}

func (r *servicePlan) GetServicePlan(serviceOfferingGUID string, planType string) (*ServicePlan, error) {
	req := rest.GetRequest("/v2/service_plans")
	if serviceOfferingGUID != "" {
		req.Query("q", "service_guid:"+serviceOfferingGUID)
	}
	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()
	plans, err := r.listServicesPlanWithPath(path)
	if err != nil {
		return nil, err
	}
	if len(plans) == 0 {
		return nil, bmxerror.New(ErrCodeServicePlanDoesNotExist,
			fmt.Sprintf("Given plan %q doesn't  exist for the service %q", planType, serviceOfferingGUID))
	}
	for _, p := range plans {
		if p.Name == planType {
			return &p, nil
		}

	}
	return nil, bmxerror.New(ErrCodeServicePlanDoesNotExist,
		fmt.Sprintf("Given plan %q doesn't  exist for the service %q", planType, serviceOfferingGUID))

}

func (r *servicePlan) listServicesPlanWithPath(path string) ([]ServicePlan, error) {
	var servicePlans []ServicePlan
	_, err := r.client.GetPaginated(path, ServicePlanResource{}, func(resource interface{}) bool {
		if servicePlanResource, ok := resource.(ServicePlanResource); ok {
			servicePlans = append(servicePlans, servicePlanResource.ToFields())
			return true
		}
		return false
	})
	return servicePlans, err
}
