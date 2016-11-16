
package awsecs

import (
	"fmt"
	"time"

	"github.com/weaveworks/scope/report"
	"github.com/weaveworks/scope/probe/docker"
)

type taskInfo struct {
	containerIDs []string
	family string
}

// return map from cluster to map of task arns to task infos
func getLabelInfo(rpt report.Report) map[string]map[string]taskInfo {
	results := make(map[string]map[string]taskInfo)
	for nodeID, node := range rpt.Container.Nodes {

		taskArn, taskArnOk := node.Latest.Lookup(docker.LabelPrefix + "com.amazonaws.ecs.task-arn")
		cluster, clusterOk := node.Latest.Lookup(docker.LabelPrefix + "com.amazonaws.ecs.cluster")
		family, familyOk := node.Latest.Lookup(docker.LabelPrefix + "com.amazonaws.ecs.task-definition-family")

		if taskArnOk && clusterOk && familyOk {
			taskMap, ok := results[cluster]
			if !ok {
				taskMap = make(map[string]taskInfo)
				results[cluster] = taskMap
			}

			task, ok := taskMap[taskArn]
			if !ok {
				task = taskInfo{containerIDs: make([]string, 0), family: family}
				taskMap[taskArn] = task
			}

			task.containerIDs = append(task.containerIDs, nodeID)
		}
	}
	return results
}

// implements Tagger
type Reporter struct {
}

func (r *Reporter) Tag(rpt report.Report) (report.Report, error) {

	now := time.Now()

	clusterMap := getLabelInfo(rpt)

	for cluster, taskMap := range clusterMap {

		client, err := newClient(cluster)
		if err != nil {
			return rpt, err
		}

		taskArns := make([]string, 0, len(taskMap))
		for taskArn, _ := range taskMap {
			taskArns = append(taskArns, taskArn)
		}

		taskServices, err := client.getTaskServices(taskArns)
		if err != nil {
			return rpt, err
		}

		// Create all the services first
		for _, serviceName := range taskServices {
			rpt.ECSServices.AddNode(report.MakeNode(serviceNodeID(serviceName)))
		}

		for taskArn, info := range taskMap {

			// new task node
			node := report.MakeNode(taskNodeID(taskArn))
			node.Latest.Set("family", now, info.family)

			rpt.ECSTask.AddNode(node)

			for _, containerID := range info.containerIDs {
				// TODO set task node as parent of container
			}

			if serviceName, ok := taskServices[taskArn]; ok {
				// TODO set service node as parent of task node
			}
		}

	}

	return rpt, nil

}

func serviceNodeID(id string) string {
	return fmt.Sprintf("%s;ECSService", id)
}

func taskNodeID(id string) string {
	return fmt.Sprintf("%s;ECSTask", id)
}
