package trainer

import (
	"fmt"
	"sort"
	"strings"
)

type (
	Cluster struct {
		Activities ActivityList
		Coords     Point
	}
	ClusterList []*Cluster
)

func findClusters(activities ActivityList, distance float64) (clusters ClusterList) {
	clusters = ClusterList{}
	for _, activity := range activities {
		if len(activity.DataPoints()) == 0 {
			continue
		}
		coords := activity.DataPoints()[0].Coords
		matchFound := false
		for _, cluster := range clusters {
			distance := coords.DistanceTo(cluster.Coords)
			if distance < distance {
				matchFound = true
				cluster.Activities = append(cluster.Activities, activity)
				break
			}
		}
		if !matchFound {
			newCluster := &Cluster{
				Activities: ActivityList{activity},
				Coords:     coords,
			}
			clusters = append(clusters, newCluster)
		}
	}
	sort.Slice(clusters, func(i, j int) bool {
		return len(clusters[i].Activities) > len(clusters[j].Activities)
	})
	return clusters
}

func (cluster *Cluster) String() string {
	return fmt.Sprintf("%s, n=%d", cluster.Coords, len(cluster.Activities))
}

func (cl ClusterList) String() string {
	lines := make([]string, 0, len(cl))
	for _, cluster := range cl {
		lines = append(lines, cluster.String())
	}
	return strings.Join(lines, "\n")
}
