package main

import "math"

type material interface {
	scatter(rIn ray, hr *hitRecord, colorAttenuation *vec3, rOut *ray) bool
}

type lambertian struct {
	albedo vec3
}

func (l lambertian) scatter(rIn ray, hr *hitRecord, colorAttenuation *vec3, rOut *ray) bool {
	scatterDir := hr.normal.add(randomVecOnHemisphere(hr.normal))
	if scatterDir.nearZero() {
		scatterDir = hr.normal
	}
	*rOut = ray{hr.point, scatterDir}
	*colorAttenuation = l.albedo
	return true
}

type metal struct {
	albedo vec3
	fuzz   float64
}

func (m metal) scatter(rIn ray, hr *hitRecord, colorAttenuation *vec3, rOut *ray) bool {
	reflectDir := rIn.dir.reflect(hr.normal).normalize().add(randomUnitVec().scale(m.fuzz))
	*rOut = ray{hr.point, reflectDir}
	*colorAttenuation = m.albedo
	return rOut.dir.dot(hr.normal) > 0
}

type dielectric struct {
	refractionIndex float64
}

func (d dielectric) scatter(rIn ray, hr *hitRecord, colorAttenuation *vec3, rOut *ray) bool {
	*colorAttenuation = vec3{1, 1, 1}
	refractionIndex := d.refractionIndex
	if hr.frontFace {
		refractionIndex = 1.0 / refractionIndex
	}
	unitDir := rIn.dir.normalize()
	cosTheta := interval{math.Inf(-1), 1.0}.clamp(unitDir.scale(-1).dot(hr.normal))
	sinTheta := math.Sqrt(1.0 - cosTheta*cosTheta)
	cannotRefract := refractionIndex*sinTheta > 1.0
	var dir vec3
	if cannotRefract || d.reflectance(cosTheta, refractionIndex) > random() {
		dir = unitDir.reflect(hr.normal)
	} else {
		dir = unitDir.refract(hr.normal, refractionIndex)
	}
	*rOut = ray{hr.point, dir}
	return true
}

func (d dielectric) reflectance(cos, refractionIndex float64) float64 {
	r0 := (1 - refractionIndex) / (1 + refractionIndex)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow((1-cos), 5.0)
}
