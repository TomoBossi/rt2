package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
)

type camera struct {
	aspectRatio               float64 // Ratio of image width over height
	imgWidth                  int     // Rendered image width in pixel count
	imgHeight                 int     // Rendered image height
	center                    vec3    // Camera center
	lookAt                    vec3    // Point in space where the camera is looking
	upDir                     vec3    // "Up" direction
	u, v, w                   vec3    // Camera frame of reference versors
	verticalFOV               float64 // Vertical view angle
	defocusAngle              float64 // Variation angle of rays through each pixel
	focalDistance             float64 // Distance from camera lookfrom point to plane of perfect focus
	defocusDiskU              vec3    // Defocus disk horizontal radius
	defocusDiskV              vec3    // Defocus disk vertical radius
	viewportUpperLeft         vec3    // Location of top-left corner of pixel 0, 0
	interPixelDeltaHorizontal vec3    // Offset to pixel to the right
	interPixelDeltaVertical   vec3    // Offset to pixel below
	antiAliasing              int     // Level of antialiasing
	maxDepth                  int     // Maximum number of ray bounces into scene
	pixels                    []byte  // Flattened image last rendered by the camera
}

type cameraParams struct {
	aspectRatio   float64 // Ratio of image width over height
	imgWidth      int     // Rendered image width in pixel count
	lookFrom      vec3    // Point in space where the camera eye is located
	lookAt        vec3    // Point in space where the camera is looking
	verticalFOV   float64 // Vertical view angle
	defocusAngle  float64 // Variation angle of rays through each pixel
	focalDistance float64 // Distance from camera lookfrom point to plane of perfect focus
	antiAliasing  int     // Level of antialiasing
	maxDepth      int     // Maximum number of ray bounces into scene
}

func cameraInit(params cameraParams) camera {
	center := params.lookFrom
	upDir := vec3{0, 1, 0}

	imgHeight := int(float64(params.imgWidth) / params.aspectRatio)
	viewportHeight := 2 * math.Tan(deg2rad(params.verticalFOV)/2) * params.focalDistance
	viewportWidth := viewportHeight * (float64(params.imgWidth) / float64(imgHeight))

	w := center.subtract(params.lookAt).normalize()
	u := upDir.cross(w).normalize()
	v := w.cross(u)

	viewportEdgeHorizontal := u.scale(viewportWidth)
	viewportEdgeVertical := v.scale(-viewportHeight)
	interPixelDeltaHorizontal := viewportEdgeHorizontal.divide(float64(params.imgWidth))
	interPixelDeltaVertical := viewportEdgeVertical.divide(float64(imgHeight))

	viewportUpperLeft := center.
		subtract(w.scale(params.focalDistance)).
		subtract((viewportEdgeHorizontal.add(viewportEdgeVertical)).scale(0.5))

	defocusRadius := params.focalDistance * math.Tan(deg2rad(params.defocusAngle/2))
	defocusDiskU := u.scale(defocusRadius)
	defocusDiskV := v.scale(defocusRadius)

	pixels := make([]byte, 4*params.imgWidth*imgHeight)
	for i := range pixels {
		pixels[i] = 255
	}

	return camera{
		aspectRatio:               params.aspectRatio,
		imgWidth:                  params.imgWidth,
		imgHeight:                 imgHeight,
		center:                    center,
		lookAt:                    params.lookAt,
		upDir:                     upDir,
		u:                         u,
		v:                         v,
		w:                         w,
		verticalFOV:               params.verticalFOV,
		defocusAngle:              params.defocusAngle,
		focalDistance:             params.focalDistance,
		defocusDiskU:              defocusDiskU,
		defocusDiskV:              defocusDiskV,
		viewportUpperLeft:         viewportUpperLeft,
		interPixelDeltaHorizontal: interPixelDeltaHorizontal,
		interPixelDeltaVertical:   interPixelDeltaVertical,
		antiAliasing:              params.antiAliasing,
		maxDepth:                  params.maxDepth,
		pixels:                    pixels,
	}
}

func rayColor(r ray, depth int, w world) vec3 {
	if depth <= 0 {
		return vec3{0, 0, 0}
	}

	var hr hitRecord
	if w.hit(r, interval{0.0001, math.Inf(1)}, &hr) {
		var rOut ray
		var colorAttenuation vec3
		if hr.mat.scatter(r, &hr, &colorAttenuation, &rOut) {
			return rayColor(rOut, depth-1, w).multiply(colorAttenuation)
		}
		return vec3{0, 0, 0}
	}

	unitDir := r.dir.normalize()
	a := 0.5 * (unitDir.y + 1.0)
	return vec3{1.0, 1.0, 1.0}.scale(1.0 - a).add(vec3{0.5, 0.7, 1.0}.scale(a))
}

func (c camera) render(w world) {
	intensity := interval{0, 0.999}
	antiAliasingDeltaHorizontal := c.interPixelDeltaHorizontal.divide(float64(c.antiAliasing + 1))
	antiAliasingDeltaVertical := c.interPixelDeltaVertical.divide(float64(c.antiAliasing + 1))
	for y := range c.imgHeight {
		for x := range c.imgWidth {
			pixelCorner := c.viewportUpperLeft.
				add(c.interPixelDeltaHorizontal.scale(float64(x))).
				add(c.interPixelDeltaVertical.scale(float64(y)))

			rayCol := vec3{0, 0, 0}
			for j := 1; j < c.antiAliasing+1; j++ {
				for i := 1; i < c.antiAliasing+1; i++ {
					viewportPoint := pixelCorner.
						add(antiAliasingDeltaHorizontal.scale(float64(i))).
						add(antiAliasingDeltaVertical.scale(float64(j)))
					rayOri := c.center
					if c.defocusAngle > 0 {
						rayOri = c.randomPointOnDefocusDisk()
					}
					rayDir := viewportPoint.subtract(rayOri)
					rayCol = rayCol.add(rayColor(ray{ori: rayOri, dir: rayDir}, c.maxDepth, w))
				}
			}

			col := rayCol.divide(float64(c.antiAliasing) * float64(c.antiAliasing))
			r := uint8(math.Floor(256 * intensity.clamp(math.Sqrt(col.x))))
			g := uint8(math.Floor(256 * intensity.clamp(math.Sqrt(col.y))))
			b := uint8(math.Floor(256 * intensity.clamp(math.Sqrt(col.z))))

			idx := 4 * (y*c.imgWidth + x)
			c.pixels[idx] = r
			c.pixels[idx+1] = g
			c.pixels[idx+2] = b
		}
	}
}

func (c camera) randomPointOnDefocusDisk() vec3 {
	v := randomVecOnUnitDisk()
	return c.center.add(c.defocusDiskU.scale(v.x)).add(c.defocusDiskV.scale(v.y))
}

func (c camera) screenshot(directory, fileName string) error {
	ext := filepath.Ext(fileName)
	path := filepath.Join(directory, fileName)

	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return err
	}

	switch ext {
	case ".ppm":
		err = savePpm(c.pixels, c.imgWidth, c.imgHeight, path)
	case ".png":
		err = savePng(c.pixels, c.imgWidth, c.imgHeight, path)
	}
	return err
}

func savePpm(pixels []byte, imgWidth, imgHeight int, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(fmt.Sprintf("P3\n%d %d\n255\n", imgWidth, imgHeight)); err != nil {
		return err
	}
	for y := range imgHeight {
		for x := range imgWidth {
			idx := 4 * (y*imgWidth + x)
			if _, err := file.WriteString(fmt.Sprintf("%d %d %d\n", pixels[idx], pixels[idx+1], pixels[idx+2])); err != nil {
				return err
			}
		}
	}
	return nil
}

func savePng(pixels []byte, imgWidth, imgHeight int, path string) error {
	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	for y := range imgHeight {
		for x := range imgWidth {
			idx := 4 * (y*imgWidth + x)
			img.Set(x, y, color.RGBA{R: pixels[idx], G: pixels[idx+1], B: pixels[idx+2], A: 255})
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}
