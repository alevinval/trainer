package trainer

import (
	"fmt"
	"sort"
	"strings"
)

// Clustering distance in meters
const clusteringDistance = 5000

type (
	Cluster struct {
		Activities ActivityList
		Coords     Point
	}
	ClusterList []*Cluster
)

func findClusters(activities ActivityList) (clusters ClusterList) {
	clusters = ClusterList{}
	for _, activity := range activities {
		if len(activity.DataPoints()) == 0 {
			continue
		}
		coords := activity.DataPoints()[0].Coords
		matchFound := false
		for _, cluster := range clusters {
			distance := coords.DistanceTo(cluster.Coords)
			if distance < clusteringDistance {
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
	lines := []string{}
	for _, cluster := range cl {
		lines = append(lines, cluster.String())
	}
	return strings.Join(lines, "\n")
}
