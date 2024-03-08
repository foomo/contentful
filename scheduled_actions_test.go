package contentful

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ExampleScheduledActionsService_Get(t *testing.T) {

	fmt.Println("HII")

	cma := NewCMA("CFPAT-fNuUQMfJ4q5Lo2hX0_hI1EoJp1-dMZN4K5rdtV3bWn8")

	assert.NotNil(t, cma)
	assert.NotNil(t, cma.ScheduledActions)

	scheduledActions, err := cma.ScheduledActions.List("qfsyzz7ytbcy", "2zgxTOq8CGHMGuusEpnJDq", "master").Next()

	if err != nil {
		log.Fatal(err)
	}
	
	spaces := scheduledActions.ToScheduledAction()
	for _, space := range spaces {
		fmt.Println(space.Sys.ID, space.Fields)
	}
}
