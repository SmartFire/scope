
package awsecs

import (
	"github.com/weaveworks/scope/report"
)

type taskInfo struct {
	containerIDs []string
	family string
}

// return map from cluster to map of task arns to task infos
func getLabelInfo(rpt report.Report) map[string][string]taskInfo {
	results := make(map[string][string]taskInfo)
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

			task, ok := results[taskArn]
			if !ok {
				task = taskInfo{containerIDs: make([]string), family: family}
				results[taskArn] = task
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

	now = time.Now()

	clusterMap := getLabelInfo()

	for cluster, taskMap := range clusterMap {

		taskServices := newClient(cluster).getTaskServices()

		// Create all the services first
		for _, serviceName := range taskServices {
			rpt.ECSServices.AddNode(report.MakeNode(serviceNodeID(serviceName)))
		}

		for taskArn, info := range taskMap {

			// new task node
			node := report.MakeNode(taskNodeID(taskArn))
			node.Latest.Set("family", now, info.family)

			for _, containerID := range info.containerIDs {
				// TODO set task node as parent of container
			}

			if serviceName, ok := taskServices[taskArn]; ok {
				// TODO set service node as parent of task node
			}
		}

	}

}

func serviceNodeID(id string) string {
	return fmt.Sprintf("%s;ECSService", id)
}

func taskNodeID(id string) string {
	return fmt.Sprintf("%s;ECSTask", id)
}
