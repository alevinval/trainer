package trainer

import (
	"fmt"
	"sort"
	"strings"
)

type (
	// Cluster contains a list of activities that happened close to a center of
	// coordinates.
	Cluster struct {
		Activities ActivityList
		Coords     Point
		TagCloud   *TagCloud
	}

	// ClusterList is a list of clusters.
	ClusterList []*Cluster
)

func findClusters(activities ActivityList, distance float64) (clusters ClusterList) {
	clusters = ClusterList{}
	for _, activity := range activities {
		if len(activity.DataPoints()) == 0 {
			continue
		}
		coords := activity.DataPoints()[0].Coords
		var matchFound bool
		for _, cluster := range clusters {
			if coords.DistanceTo(cluster.Coords) < distance {
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

	for _, cluster := range clusters {
		cluster.TagCloud = tagCloudFromActivities(cluster.Activities)
	}
	return clusters
}

func (cluster *Cluster) String() string {
	tags := strings.Join(cluster.TagCloud.tags, " ")
	return fmt.Sprintf("%s, n=%d, tags: %s", cluster.Coords, len(cluster.Activities), tags)
}

func (cl ClusterList) String() string {
	lines := make([]string, 0, len(cl))
	for _, cluster := range cl {
		lines = append(lines, cluster.String())
	}
	return strings.Join(lines, "\n")
}
