package main

type hitRecord struct {
	point     vec3
	normal    vec3
	t         float64
	frontFace bool
	mat       material
}

type hittable interface {
	hit(r ray, tInterval interval, record *hitRecord) bool
}

func (hr *hitRecord) setFaceNormal(r ray, outwardUnitNormal vec3) {
	hr.frontFace = r.dir.dot(outwardUnitNormal) < 0
	if hr.frontFace {
		hr.normal = outwardUnitNormal
	} else {
		hr.normal = outwardUnitNormal.scale(-1)
	}
}
