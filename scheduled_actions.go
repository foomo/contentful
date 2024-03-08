package contentful

import (
	"fmt"
	"net/url"
)

// SpacesService model
type ScheduledActionsService service

// ScheduledActions model
type ScheduledActions struct {
	Sys           *Sys   `json:"sys,omitempty"`
	Fields map[string]interface{} `json:"fields,omitempty"`
}

// // MarshalJSON for custom json marshaling
// func (scheduledActions *ScheduledActions) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(&struct {
// 		Name          string `json:"name,omitempty"`
// 		DefaultLocale string `json:"defaultLocale,omitempty"`
// 	}{
// 		Name:          scheduledActions.Name,
// 		DefaultLocale: scheduledActions.DefaultLocale,
// 	})
// }

// // GetVersion returns entity version
// func (scheduledActions *ScheduledActions) GetVersion() int {
// 	version := 1
// 	if scheduledActions.Sys != nil {
// 		version = scheduledActions.Sys.Version
// 	}

// 	return version
// }

// // List creates a spaces collection
// func (service *ScheduledActionsService) List() *Collection {
// 	req, _ := service.c.newRequest("GET", "/spaces", nil, nil)

// 	col := NewCollection(&CollectionOptions{})
// 	col.c = service.c
// 	col.req = req

// 	return col
// }

// Get returns a single scheduledActions entity
func (service *ScheduledActionsService) Get(spaceID string, entryID string, environmentID string) (*ScheduledActions, error) {
	// path := "/spaces/qfsyzz7ytbcy/scheduled_actions"
	// path := fmt.Sprintf("/spaces/%s%s/scheduled_actions/%s", spaceID, getEnvPath(service.c), entryID)
	path := fmt.Sprintf("/spaces/%s/scheduled_actions", spaceID)
	
	fmt.Println(path)
	fmt.Println(service.c.BaseURL)

	query := url.Values{}
	
	query.Add("entity.sys.id", "2zgxTOq8CGHMGuusEpnJDq")
	query.Add("environment.sys.id", "master")
	query.Add("status[in]", "scheduled")

	fmt.Println(query)

	
	method := "GET"
	

	req, err := service.c.newRequest(method, path, query, nil)

	if err != nil {
		return &ScheduledActions{}, err
	}

	var scheduled_actions ScheduledActions
	if ok := service.c.do(req, &scheduled_actions); ok != nil {
		return &ScheduledActions{}, ok
	}

	return &scheduled_actions, nil
}

// // Upsert updates or creates a new scheduledActions
// func (service *ScheduledActionsService) Upsert(scheduledActions *ScheduledActions) error {
// 	bytesArray, err := json.Marshal(scheduledActions)
// 	if err != nil {
// 		return err
// 	}

// 	var path string
// 	var method string

// 	if scheduledActions.Sys != nil && scheduledActions.Sys.CreatedAt != "" {
// 		path = fmt.Sprintf("/spaces/%s%s", scheduledActions.Sys.ID, getEnvPath(service.c))
// 		method = "PUT"
// 	} else {
// 		path = "/spaces"
// 		method = "POST"
// 	}

// 	req, err := service.c.newRequest(method, path, nil, bytes.NewReader(bytesArray))
// 	if err != nil {
// 		return err
// 	}

// 	req.Header.Set("X-Contentful-Version", strconv.Itoa(scheduledActions.GetVersion()))

// 	return service.c.do(req, scheduledActions)
// }

// // Delete the given scheduledActions
// func (service *ScheduledActionsService) Delete(scheduledActions *ScheduledActions) error {
// 	path := fmt.Sprintf("/spaces/%s", scheduledActions.Sys.ID)
// 	method := "DELETE"

// 	req, err := service.c.newRequest(method, path, nil, nil)
// 	if err != nil {
// 		return err
// 	}

// 	version := strconv.Itoa(scheduledActions.Sys.Version)
// 	req.Header.Set("X-Contentful-Version", version)

// 	return service.c.do(req, nil)
// }
