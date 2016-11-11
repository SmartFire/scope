
package awsecs

import (
)

struct labelInfo {
	containerID string
	cluster string
	taskArn string
	family string
}

func getLabelInfos() []labelInfo {
	// TODO get labels for all containers and look for com.amazonaws.ecs.{cluster,task-arn,task-definition-family}
}
