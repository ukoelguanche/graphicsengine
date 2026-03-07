package core

type Point struct {
	X, Y float64
}

func (p1 Point) Equals(p2 Point) bool { return p1.X == p2.X && p1.Y == p2.Y }
func (p Point) SetX(x float64) Point  { return Point{X: x, Y: p.Y} }
func (p Point) SetY(y float64) Point  { return Point{X: p.X, Y: y} }

func (p Point) IncX(x float64) Point { return Point{X: p.X + x, Y: p.Y} }
func (p Point) IncY(y float64) Point { return Point{X: p.X, Y: p.Y + y} }
