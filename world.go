package main

type world struct {
	objects []hittable
}

func (w world) hit(r ray, tInterval interval, hr *hitRecord) bool {
	var tempHr hitRecord
	hitAnything := false
	closest := tInterval.max
	for _, object := range w.objects {
		if object.hit(r, interval{tInterval.min, closest}, &tempHr) {
			hitAnything = true
			closest = tempHr.t
			*hr = tempHr
		}
	}
	return hitAnything
}
