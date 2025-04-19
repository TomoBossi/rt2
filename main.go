package main

func main() {
	world := &world{
		objects: []hittable{
			sphere{
				center: vec3{0, 0, -1.2},
				radius: 0.5,
				mat:    lambertian{albedo: vec3{0.1, 0.2, 0.5}},
			},
			sphere{
				center: vec3{-1, 0, -1},
				radius: 0.5,
				mat:    dielectric{refractionIndex: 1.5},
			},
			sphere{
				center: vec3{-1, 0, -1},
				radius: 0.4,
				mat:    dielectric{refractionIndex: 1.0 / 1.5},
			},
			sphere{
				center: vec3{1, 0, -1},
				radius: 0.5,
				mat:    metal{albedo: vec3{0.8, 0.6, 0.2}, fuzz: 0.2},
			},
			sphere{
				center: vec3{0, -100.5, -1},
				radius: 100,
				mat:    metal{albedo: vec3{0.8, 0.8, 0.0}, fuzz: 0.0},
			},
		},
	}

	camera := cameraInit(cameraParams{
		imgWidth:      200,
		aspectRatio:   16.0 / 9.0,
		verticalFov:   60.00,
		lookFrom:      vec3{-0.183, -0.168, -0.463},
		lookAt:        vec3{0.572, -0.365, -1.088},
		defocusAngle:  0,
		focalDistance: 1,
		antiAliasing:  1,
		maxDepth:      10,
	})

	defer close(camera.renderJobQueue)

	if game := true; game {
		gameInit(gameParams{camera: camera, world: world, fpsCap: 30, fullscreen: true})
	} else {
		camera.render(world)
		err := camera.screenshot("./out/", "image.png")
		if err != nil {
			panic(err)
		}
	}
}
