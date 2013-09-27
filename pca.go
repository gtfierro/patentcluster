package patentcluster

import (
    "fmt"
    "math"
    "github.com/skelterjohn/go.matrix"
)

// Tagset is map[string]int of all the tags
/*
Need to do the following
    * construct the sparse data matrix of all the patents, P
    * construct the covariance matrix P'P, C
    * make the columns of C zero-mean
    * compute the eigenvalues and eigenvectors of C
    * sort these, and return the top 3
    * sum up the eigenvalues so we can compute the variance of the eigenvalues
    * provide convenience function to get the 2d/3d coordinates of a given patent

*/

/** 
    takes as an argument the Tagset created by dbscan.go and 
    returns a mapping of tag => vector index
*/
func define_tag_order(tagset map[string]int) map[string]int {
    idx := 0
    tagorder := make(map[string]int)
    for v, _ := range tagset {
        tagorder[v] = idx
        idx += 1
    }
    return tagorder
}

/**
    Given the 
*/
func Create_sparse_data_matrix(tagorder map[string]int, data map[string]([]string)) *matrix.SparseMatrix {
    number_of_tags := len(tagorder)
    number_of_points := len(data)
    matrix := matrix.ZerosSparse(number_of_tags, number_of_points) // rows, columns
    col := 0
    for _, taglist := range data {
        for _, tag := range taglist {
            row := tagorder[tag]
            matrix.Set(row, col, 1)
        }
        col += 1
    }
    return matrix
}

/**
    Given a pointer to a dense matrix, converts
    each of the columns to be zero-mean
*/
func demean(cov *matrix.DenseMatrix, num_rows int) {
    /* compute average column */
    avg_column := make([]float64, num_rows)
    tmp_column := make([]float64, num_rows)
    /* sum up columns */
    for j := 0; j < num_rows; j+=1 {
        cov.BufferCol(j, tmp_column)
        for idx, val := range tmp_column {
            avg_column[idx] += val
        }
    }
    /* divide by number of columns */
    for idx, val := range avg_column {
        avg_column[idx] = val / float64(num_rows)
    }
    /* subtract avg_column from all columns in cov */
    for i := 0; i < num_rows; i+=1 {
        for j := 0; j < num_rows; j+= 1 {
            cov.Set(i, j, cov.Get(i,j) - avg_column[i])
        }
    }
}

/**
    Takes in the sparse data matrix created by Create_sparse_data_matrix, the rank of that
    matrix (which should be just be the number of data points), and `n`, the number of top
    eigenvalues/eigenvectors we want to return
*/
func Compute_n_eigenstuffs(matrix *matrix.SparseMatrix, rank, num_columns, n int) ([]float64, []([]float64)) {
    fmt.Println("Computing covariance data matrix...")
    matrix_t := matrix.Transpose()
    cov, _ := matrix_t.TimesSparse(matrix)
    dense := cov.DenseMatrix()
    demean(dense, rank)
    fmt.Println("Computing eigenstuffs...")
    V, D, _ := dense.Eigen() // should be V, D, _ := .... V will contain the eigenvectors
    eigenvalues := make([]float64, rank)
    eigenvectors := []([]float64){}
    for i := 0; i < n; i+=1 {
        ev := make([]float64, num_columns)
        V.BufferCol(rank-1-i, ev)
        eigenvectors = append(eigenvectors, ev)
    }
    D.BufferDiagonal(eigenvalues)
    max := float64(0)
    for val := range eigenvalues {
        max = math.Max(max, float64(val))
    }
    fmt.Println(max)
    fmt.Println(len(eigenvalues))
    fmt.Println(eigenvalues[len(eigenvalues)-n:])
    fmt.Println(len(eigenvectors))
    return eigenvalues, eigenvectors
}

func Compute_coordinates(data map[string]([]string)) {
    fmt.Println("Making patents...")
    patents := Make_patents(data)
    fmt.Println("Defining tag order...")
    tagorder := define_tag_order(Tagset)
    fmt.Println(len(patents),"patents found")
    fmt.Println("Constructing sparse matrix...")
    matrix := Create_sparse_data_matrix(tagorder, data)
    Compute_n_eigenstuffs(matrix, len(data), len(tagorder), 3)
}
