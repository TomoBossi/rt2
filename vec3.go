package main

import (
	"fmt"
	"math"
)

type vec3 struct {
	x, y, z float64
}

// func randomVec() vec3 {
// 	return vec3{random(), random(), random()}
// }

func randomVecIn(min, max float64) vec3 {
	return vec3{randomIn(min, max), randomIn(min, max), randomIn(min, max)}
}

func randomUnitVec() vec3 {
	for {
		v := randomVecIn(-1, 1)
		l2Squared := v.l2Squared()
		if 1e-160 < l2Squared && l2Squared <= 1 {
			return v.divide(math.Sqrt(l2Squared))
		}
	}
}

func randomVecOnHemisphere(normal vec3) vec3 {
	unitVec := randomUnitVec()
	if unitVec.dot(normal) > 0.0 {
		return unitVec
	}
	return unitVec.scale(-1)
}

func randomVecOnUnitDisk() vec3 {
	for {
		v := vec3{randomIn(-1, 1), randomIn(-1, 1), 0}
		if v.l2Squared() < 1 {
			return v
		}
	}
}

func (v vec3) add(u vec3) vec3 {
	return vec3{v.x + u.x, v.y + u.y, v.z + u.z}
}

func (v vec3) subtract(u vec3) vec3 {
	return vec3{v.x - u.x, v.y - u.y, v.z - u.z}
}

func (v vec3) multiply(u vec3) vec3 {
	return vec3{v.x * u.x, v.y * u.y, v.z * u.z}
}

func (v vec3) scale(t float64) vec3 {
	return vec3{v.x * t, v.y * t, v.z * t}
}

func (v vec3) divide(t float64) vec3 {
	return v.scale(1 / t)
}

func (v vec3) dot(u vec3) float64 {
	return v.x*u.x + v.y*u.y + v.z*u.z
}

func (v vec3) l2() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
}

func (v vec3) l2Squared() float64 {
	return v.x*v.x + v.y*v.y + v.z*v.z
}

func (v vec3) normalize() vec3 {
	return v.scale(1 / v.l2())
}

func (v vec3) nearZero() bool {
	min := 10e-8
	return abs(v.x) < min && abs(v.y) < min && abs(v.z) < min
}

func (v vec3) reflect(normal vec3) vec3 {
	return v.subtract(normal.scale(2 * v.dot(normal)))
}

func (uv vec3) refract(normal vec3, refractionIndex float64) vec3 {
	cosTheta := interval{math.Inf(-1), 1.0}.clamp(uv.scale(-1).dot(normal))
	rOutPerpendicular := uv.add(normal.scale(cosTheta)).scale(refractionIndex)
	rOutParallel := normal.scale(-math.Sqrt(abs(1.0 - rOutPerpendicular.l2Squared())))
	return rOutPerpendicular.add(rOutParallel)
}

func (v vec3) cross(u vec3) vec3 {
	return vec3{
		v.y*u.z - v.z*u.y,
		v.z*u.x - v.x*u.z,
		v.x*u.y - v.y*u.x,
	}
}

func (v vec3) rotateAroundAxis(axis vec3, angle float64) vec3 {
	axis = axis.normalize()
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	return v.scale(cos).
		add(axis.cross(v).scale(sin)).
		add(axis.scale(axis.dot(v) * (1 - cos)))
}

func (v vec3) String() string {
	return fmt.Sprintf("vec3{%.3f, %.3f, %.3f}", v.x, v.y, v.z)
}
