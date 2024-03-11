package contentful

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ExampleScheduledActionsService_Get(t *testing.T) {

	cma := NewCMA("cma-token")
	assert.NotNil(t, cma)
	assert.NotNil(t, cma.ScheduledActions)

	scheduledActions, err := cma.ScheduledActions.Get("space-id", "entry-id", "env")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v", scheduledActions)
}
