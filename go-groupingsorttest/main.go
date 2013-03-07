package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type entry struct {
	//	these categories don't really need to be *ordered* ascending or descending per-se:
	//	just clustered/grouped together hierarchically: all identical catA entries should
	//	be together, within that cluster all identical catB entries, within that all catC.
	catA, catB, catC int

	//	only within properly grouped clusters then, sort by ordering "dist" ascending:
	dist float64
}

type entries struct {
	e   []entry
	sts stats
}

func (e *entries) calcStats() {
	e.sts = make(stats)

	for _, v := range e.e {
		e.sts.addDatum(v.dist, [3]int{v.catA, -1, -1})
		e.sts.addDatum(v.dist, [3]int{v.catA, v.catB, -1})
		e.sts.addDatum(v.dist, [3]int{v.catA, v.catB, v.catC})
	}
}

func (me entries) Len() int {
	return len(me.e)
}

func compare(a, b float64) int {
	switch {
	case a == b:
		return 0
	case a < b:
		return -1
	}
	return 1
}

func (me entries) Less(i, j int) bool {
	iv := me.e[i]
	jv := me.e[j]

	if iv.catA != jv.catA {
		switch compare(me.sts[[3]int{iv.catA, -1, -1}].mean, me.sts[[3]int{jv.catA, -1, -1}].mean) {
		case 0:
			return iv.catA < jv.catA
		case -1:
			return true
		case 1:
			return false
		}
	}

	if iv.catB != jv.catB {
		switch compare(me.sts[[3]int{iv.catA, iv.catB, -1}].mean, me.sts[[3]int{jv.catA, jv.catB, -1}].mean) {
		case 0:
			return iv.catB < jv.catB
		case -1:
			return true
		case 1:
			return false
		}
	}

	if iv.catC != jv.catC {
		switch compare(me.sts[[3]int{iv.catA, iv.catB, iv.catC}].mean, me.sts[[3]int{jv.catA, jv.catB, jv.catC}].mean) {
		case 0:
			return iv.catC < jv.catC
		case -1:
			return true
		case 1:
			return false
		}
	}

	return iv.dist < jv.dist
}

func (me entries) Swap(i, j int) {
	me.e[i], me.e[j] = me.e[j], me.e[i]
}

type stat struct {
	count int
	sum   float64
	mean  float64
}

type stats map[[3]int]stat

func (s stats) addDatum(dist float64, key [3]int) {
	v := s[key]
	v.count++
	v.sum += dist
	v.mean = v.sum / float64(v.count)
	s[key] = v
}

func main() {
	rand.Seed(time.Now().UnixNano())
	catCs := []int{0, 1}
	catAs := []int{2, 4, 6}
	catBs := []int{3, 5, 7}
	dists := []float64{10, 20, 40, 80}
	all := make([]entry, 16)
	for i := 0; i < len(all); i++ {
		all[i].catC = catCs[rand.Intn(len(catCs))]
		all[i].catA = catAs[rand.Intn(len(catAs))]
		all[i].catB = catBs[rand.Intn(len(catBs))]
		all[i].dist = dists[rand.Intn(len(dists))]
	}
	fmt.Println("\nUNSORTED:\n===>")
	printAll(all, false)

	e := entries{
		e: all,
	}
	e.calcStats()
	sort.Sort(e)
	fmt.Println("\nSORTED:\n===>")
	printAll(all, true)
}

func printAll(all []entry, sorted bool) {
	lastA := -1
	for i := 0; i < len(all); i++ {
		if sorted && all[i].catA != lastA {
			lastA = all[i].catA
			fmt.Println("")
		}
		fmt.Printf("catA=%v\tcatB=%v\tcatC=%v\tdist=%v\n", all[i].catA, all[i].catB, all[i].catC, all[i].dist)
	}
}
