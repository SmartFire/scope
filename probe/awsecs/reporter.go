
package awsecs

import (
)

struct taskInfo {
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
struct Reporter {
}

func (r *Reporter) Tag(rpt report.Report) (report.Report, error) {

	clusterMap := getLabelInfo()

	for cluster, taskMap := range clusterMap {

		taskServices := newClient(cluster).getTaskServices()

		for _, serviceName := range taskServices {
			// TODO create a ecs service node
		}

		for taskArn, info := range taskMap {
			// TODO create task node with family
			for _, containerID := range info.containerIDs {
				// TODO set task node as parent of container
			}
		}

	}

}
