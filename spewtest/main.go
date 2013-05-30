package main

import (
	"github.com/goforks/go-spew/spew"

	u3d "github.com/go3d/go-3dutil"
)

func main() {
	var fc u3d.FrustumCoords
	fc.X, fc.Y = 10, 20
	fc.BR.X, fc.BR.Y, fc.BR.Z = 5, 7, 9
	spew.Dump(fc)
}
