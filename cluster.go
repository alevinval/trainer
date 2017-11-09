package trainer

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// DistanceCriteria clusters activities that fall within a distance
// threshold between them.
var DistanceCriteria = func(distance float64) clusteringCriteria {
	return func(c *Cluster, a *Activity) bool {
		startingPoint := a.DataPoints()[0].Coords
		return startingPoint.DistanceTo(c.Coords) < distance
	}
}

type (
	// Cluster contains a list of activities that happened close to a center of
	// coordinates.
	Cluster struct {
		Activities ActivityList
		Coords     Point
	}

	// ClusterList is a list of clusters.
	ClusterList []*Cluster

	clusteringCriteria func(*Cluster, *Activity) bool
)

// GetClusters clusters a list of activities that match a certain criteria.
func GetClusters(activities ActivityList, criteria clusteringCriteria) ClusterList {
	clusters := ClusterList{}
	for _, activity := range activities {
		if len(activity.DataPoints()) == 0 {
			continue
		}
		cluster, err := getMatchingCluster(clusters, activity, criteria)
		if err == nil {
			cluster.Activities = append(cluster.Activities, activity)
		} else {
			newCluster := &Cluster{
				Activities: ActivityList{activity},
				Coords:     activity.DataPoints()[0].Coords,
			}
			clusters = append(clusters, newCluster)
		}
	}
	sort.Slice(clusters, func(i, j int) bool {
		return len(clusters[i].Activities) > len(clusters[j].Activities)
	})
	return clusters
}

func getMatchingCluster(clusters ClusterList, candidate *Activity, criteria clusteringCriteria) (*Cluster, error) {
	for _, cluster := range clusters {
		if criteria(cluster, candidate) {
			return cluster, nil
		}
	}
	return nil, errors.New("no cluster matches criteria")
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
