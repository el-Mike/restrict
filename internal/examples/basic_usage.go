package examples

import (
	"fmt"
	"log"

	"github.com/el-mike/restrict"
	"github.com/el-mike/restrict/adapters"
)

var policy = &restrict.PolicyDefinition{
	Roles: restrict.Roles{
		"User": {
			Grants: restrict.GrantsMap{
				"Conversation": {
					&restrict.Permission{Action: "read"},
					&restrict.Permission{Action: "create"},
				},
			},
		},
	},
}

func main() {
	// Create an instance of PolicyManager, which will be responsible for handling given PolicyDefinition.
	// You can use one of the built-in persistence adapters (in-memory or json/yaml file adapters), or provide your own.
	policyManager, err := restrict.NewPolicyManager(adapters.NewInMemoryAdapter(policy), true)
	if err != nil {
		log.Fatal(err)
	}

	manager := restrict.NewAccessManager(policyManager)

	if err = manager.Authorize(&restrict.AccessRequest{
		Subject:        &User{},
		Resource:       &Conversation{},
		Actions:        []string{"read", "delete"},
		SkipConditions: false,
	}); err != nil {
		fmt.Print(err) // Access denied for actions: "delete". Reason: Permission for action: "delete" is not granted for Resource: "Conversation"
	}

	if accessError, ok := err.(*restrict.AccessDeniedError); ok {
		// Error() implementation. Returns a message in a form:
		// Access denied for action: "...". Reason: Permission for action: "..." is not granted for Resource: "..."
		fmt.Println(accessError)
		// Returns an AccessRequest that failed.
		fmt.Println(accessError.Request)
		// We can use Errors.First() to get the first encountered PermissionError.
		fmt.Println(accessError.Errors.First())

		// We can use Errors property to loop over all PermissionErrors.
		for _, permissionErr := range accessError.Errors {
			fmt.Println(permissionErr)
			fmt.Println(permissionErr.Action)
			fmt.Println(permissionErr.RoleName)
			fmt.Println(permissionErr.ResourceName)

			// If the reason of an Permission was failed Condition,
			// this helper method returns it directly. Otherwise, nil will be returned.
			failedCondition := permissionErr.FailedCondition()

			// It can be later cast to the type you want.
			if emptyCondition, ok := failedCondition.(*restrict.EmptyCondition); failedCondition != nil && ok {
				fmt.Print(emptyCondition.ID)
			}
		}
	}
}
