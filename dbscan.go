package patentcluster

type DBSCAN struct {
    set_of_points map[string](*Patent)
    epsilon float64
    min_cluster_points int
}

/**
    Changes cluster_id of Patent
*/
func (db *DBSCAN) ChangeClusterID(point *Patent, cluster_id string) {
    point.cluster_id = cluster_id
}

/**
    Given a list of points, iterates through them and changes the cluster_id
    of each. This is called in order to identify a set of patents as 
    being in a cluster
*/
func (db *DBSCAN) ChangeClusterIDs(points [](*Patent), cluster_id string) {
    for _, val := range points {
        db.ChangeClusterID(val, cluster_id)
    }
}

/**
    For the given patent, loops through all other patents in the DBSCAN
    intance's set_of_points and returns a slice of all Patents within
    distance `epsilon` of that patent.

    At some point, this should be optimized to use a modified R* tree.
    Because we're using Jaccard distance to compare patents, which isn't
    an absolute distance, a normal R* tree might not work.
*/
func (db *DBSCAN) RegionQuery(point *Patent, epsilon float64) [](*Patent) {
    returned_points := [](*Patent){}
    for _, patent := range db.set_of_points {
        if point.JaccardDistance(patent) <= epsilon {
            returned_points = append(returned_points, patent)
        }
    }
    return returned_points
}
