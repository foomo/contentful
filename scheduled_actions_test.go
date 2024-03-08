package contentful

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ExampleScheduledActionsService_Get(t *testing.T) {

	fmt.Println("HII")

	cma := NewCMA("TOKEN HERE")

	assert.NotNil(t, cma)
	assert.NotNil(t, cma.ScheduledActions)

	scheduledActions, err := cma.ScheduledActions.Get("qfsyzz7ytbcy", "2zgxTOq8CGHMGuusEpnJDq", "master")

	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("%v",scheduledActions.Fields)
}
