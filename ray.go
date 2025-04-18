package main

type ray struct {
	ori, dir vec3
}

func (r ray) at(t float64) vec3 {
	return r.ori.add(r.dir.scale(t))
}
