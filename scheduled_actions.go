package contentful

import (
	"fmt"
	"net/url"
)

// SpacesService model
type ScheduledActionsService service

// ScheduledActions model
type ScheduledActions struct {
	Sys         *Sys   `json:"sys,omitempty"`
	Action      string `json:"action,omitempty"`
	ScheduleFor struct {
		Datetime string `json:"datetime,omitempty"`
		TimeZone string `json:"timezone,omitempty"`
	} `json:"scheduledFor,omitempty"`
}

// Get returns a single scheduledActions entity
func (service *ScheduledActionsService) Get(spaceID string, entryID string, environmentID string) (*ScheduledActions, error) {
	path := fmt.Sprintf("/spaces/%s/scheduled_actions", spaceID)

	fmt.Println(path)
	fmt.Println(service.c.BaseURL)

	query := url.Values{}

	query.Add("entity.sys.id", entryID)
	query.Add("environment.sys.id", environmentID)
	query.Add("sys.status[in]", "scheduled")

	fmt.Println(query)

	method := "GET"

	req, err := service.c.newRequest(method, path, query, nil)
	if err != nil {
		return &ScheduledActions{}, err
	}

	col := NewCollection(&CollectionOptions{})
	col.c = service.c
	col.req = req

	if ok := service.c.do(req, &col); ok != nil {
		return &ScheduledActions{}, ok
	}

	for _, ct := range col.ToScheduledAction() {
		fmt.Println(ct)
	}
	return &ScheduledActions{}, nil
}
