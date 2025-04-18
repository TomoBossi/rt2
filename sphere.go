package main

import "math"

type sphere struct {
	center vec3
	radius float64
	mat    material
}

func (s sphere) hit(r ray, tInterval interval, hr *hitRecord) bool {
	oc := s.center.subtract(r.ori)
	a := r.dir.l2Squared()
	h := r.dir.dot(oc)
	c := oc.l2Squared() - s.radius*s.radius
	discriminant := h*h - a*c

	if discriminant < 0 {
		return false
	}

	discriminantSqrt := math.Sqrt(discriminant)
	root := (h - discriminantSqrt) / a
	if !tInterval.surrounds(root) {
		root = (h + discriminantSqrt) / a
		if !tInterval.surrounds(root) {
			return false
		}
	}

	hr.t = root
	hr.point = r.at(root)
	outwardNormal := hr.point.subtract(s.center).divide(s.radius)
	hr.setFaceNormal(r, outwardNormal)
	hr.mat = s.mat
	return true
}
