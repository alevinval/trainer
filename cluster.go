package trainer

import (
	"fmt"
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
