package patentcluster

import (
    "os"
    "fmt"
    "bufio"
    "sort"
    "sync"
)

var NOISE = "NOISE" // cluster_id string for noisy patents
var UNCLASSIFIED = "UNCLASSIFIED" // cluster_id string for unclassified patents
var wg sync.WaitGroup

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
    distance `epsilon` of that patent. This includes the given patent.

    At some point, this should be optimized to use a modified R* tree.
    Because we're using Jaccard distance to compare patents, which isn't
    an absolute distance, a normal R* tree might not work.
*/
func (db *DBSCAN) RegionQuery(point *Patent) [](*Patent) {
    returned_points := [](*Patent){}
    for _, patent := range db.set_of_points {
        if point.JaccardDistance(patent) <= db.epsilon {
            returned_points = append(returned_points, patent)
        }
    }
    return returned_points
}

func (db *DBSCAN) PRegionQuery(point *Patent) [](*Patent) {
    returned_points := [](*Patent){}
    results := make(chan *Patent)
    done := make(chan bool)
    go func(results chan *Patent, done chan bool) {
        for {
            select {
            case val := <- results:
                returned_points = append(returned_points, val)
            case <- done:
                break
            }
        }
    }(results, done)

    for _, patent := range db.set_of_points {
        wg.Add(1)
        go func(point, patent *Patent, results chan *Patent) {
            if point.JaccardDistance(patent) <= db.epsilon {
                results <- patent
            }
            wg.Done()
        }(point, patent, results)
    }
    wg.Wait()
    done <- true
    return returned_points
}

/**
    Given a point and a list of points, returns a copy of the list
    with every instance of `point` removed
*/
func remove_point_from_seeds(point *Patent, seeds [](*Patent)) [](*Patent) {
    newseeds := [](*Patent){}
    for _, patent := range seeds {
        if point.number == patent.number {
            continue
        }
        newseeds = append(newseeds, patent)
    }
    return newseeds
}

/**
   Attempts to classify the set of points within the region surrounding the
   argument `point`. If the point's region does not contain the requisite
   number of points, it is classified as NOISE. Otherwise, we classify all
   region points as belonging to the same cluster. We then loop through all the
   region points and attempt to associate them with the same cluster.

   Returns TRUE if `point` belongs to a cluster, and FALSE otherwise
*/
func (db *DBSCAN) ExpandCluster(point *Patent, cluster_id string) bool {
    seeds := db.RegionQuery(point)
    if len(seeds) < db.min_cluster_points {
        db.ChangeClusterID(point, NOISE);
        return false
    }
    db.ChangeClusterIDs(seeds, cluster_id)
    seeds = remove_point_from_seeds(point, seeds)
    for {
        if len(seeds) == 0 {
            break
        }
        current_point := seeds[0]
        result := db.RegionQuery(current_point)
        if len(result) > db.min_cluster_points {
            for _, result_point := range result {
                if result_point.cluster_id == NOISE ||
                   result_point.cluster_id == UNCLASSIFIED {
                       if result_point.cluster_id == UNCLASSIFIED {
                           seeds = append(seeds, result_point)
                       }
                       db.ChangeClusterID(result_point, cluster_id)
               }
           }
       }
       seeds = seeds[1:]
   }
   return true;
}

/**
    returns the number of the first Patent in the DBSCAN.set_of_points
    that is classified as the cluster_id `clid`
*/
func (db *DBSCAN) nextClusterID(clid string) (number string) {
    for _, pat := range db.set_of_points {
        if pat.cluster_id == clid {
            return pat.number
        }
    }
    return ""
}

/**
    Runs the DBSCAN algorithm, classifying all points in
    DBSCAN.set_of_points as belonging to a cluster or as NOISE.
*/
func (db *DBSCAN) Run() {
    /* find first patent classified as NOISE */
    cluster_id := db.nextClusterID(UNCLASSIFIED)
    for _, point := range db.set_of_points {
        if point.cluster_id == UNCLASSIFIED {
            if db.ExpandCluster(point, cluster_id) {
                cluster_id = db.nextClusterID(cluster_id)
            }
        }
    }
}

/**
    Dumps db.set_of_points to a CSV file:
    patentnumber, cluster_id

    Places all points that are part of a cluster
    at the beginning of the file. All NOISE points
    are listed at the end.
*/
func (db *DBSCAN) To_file(filename string) {
    outfile, err := os.Create(filename)
    if err != nil {
        fmt.Println("Could not output to file", filename, ":", err)
        return
    }
    defer outfile.Close()
    writer := bufio.NewWriter(outfile)

    for _, patent := range db.set_of_points {
        if patent.cluster_id != UNCLASSIFIED && patent.cluster_id != NOISE {
            line := patent.number + ", " + patent.cluster_id + "\n"
            writer.WriteString(line)
        }
    }
    for _, patent := range db.set_of_points {
        if patent.cluster_id == UNCLASSIFIED || patent.cluster_id == NOISE {
            line := patent.number + ", " + patent.cluster_id + "\n"
            writer.WriteString(line)
        }
    }
    writer.Flush()
}

/**
    For an instance of DBSCAN (after Run() has been called), returns
    * number of clusters
    * mean cluster size
    * median cluster size
    * size of largest Cluster
    * list of patents in largest Cluster

*/
func (db *DBSCAN) Compute_Stats() (int, float64, int, int, [](*Patent)) {
    largest_cluster := [](*Patent){}
    largest_cluster_key := ""
    cluster_counts := make(map[string]int)
    for _, v := range db.set_of_points {
        if v.cluster_id != NOISE {
            cluster_counts[v.cluster_id] += 1
        }
    }
    list_of_counts := []int{}
    max := 0
    sum := 0.0
    for k, v:= range cluster_counts {
        if v > max {
            max = v
            largest_cluster_key = k
        }
        sum += float64(v)
        list_of_counts = append(list_of_counts, v)
    }
    for _, v := range db.set_of_points {
        if v.cluster_id == largest_cluster_key {
            largest_cluster = append(largest_cluster, v)
        }
    }

    mean_cluster_size := sum / float64(len(cluster_counts))
    sort.Ints(list_of_counts)
    median_key := len(list_of_counts) / 2
    median_cluster_size := 0
    if len(list_of_counts) > 0 {
        median_cluster_size = list_of_counts[median_key]
    }

    return len(cluster_counts), mean_cluster_size, median_cluster_size, len(largest_cluster), largest_cluster
}

/**
    Takes a slice of patent pointers and initializes an instance of the
    DBSCAN algorithm. Does not run the algorithm
*/
func Init_DBSCAN(points [](*Patent), epsilon float64, min_cluster_points int) (*DBSCAN) {
    db := new(DBSCAN)
    db.epsilon = epsilon
    db.min_cluster_points = min_cluster_points
    db.set_of_points = make(map[string](*Patent))
    for _, patent := range points {
        patent.cluster_id = UNCLASSIFIED
        db.set_of_points[patent.number] = patent
    }
    return db
}
