package main

import uhash "github.com/metaleap/go-util/hash"

func main() {
	hash := uhash.Fnv1a
	println("3s:")
	println(hash([]int{1, 2, 3}))
	println(hash([]int{1, 3, 2}))
	println(hash([]int{2, 3, 1}))
	println(hash([]int{2, 1, 3}))
	println(hash([]int{3, 1, 2}))
	println(hash([]int{3, 2, 1}))
	println("2s:")
	println(hash([]int{3, 2}))
	println(hash([]int{2, 3}))
	println(hash([]int{1, 2}))
	println(hash([]int{2, 1}))
	println(hash([]int{1, 3}))
	println(hash([]int{3, 1}))
	println("0 / -1:")
	hash = uhash.Fnv1a
	println(hash([]int{-1}))
	println(hash([]int{0, -1}))
	println(hash([]int{-1, -1}))
	println(hash([]int{0, -1, -1}))
}
