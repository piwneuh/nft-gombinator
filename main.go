package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
)

func main() {
	dir := "./images"

	imagePaths := make(map[string][]string)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			layer := filepath.Base(filepath.Dir(path))
			imagePaths[layer] = append(imagePaths[layer], path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	combinations := generateCombinations(imagePaths)
	for i, combination := range combinations {
		result := mergeImages(combination)
		saveImage(result, "output", fmt.Sprintf("output_%d.png", i+1))
	}
}

func generateCombinations(images map[string][]string) [][]string {
	// magic slice convert
	var layerImages [][]string
	for _, imgs := range images {
		layerImages = append(layerImages, imgs)
	}

	var result [][]string
	generate(layerImages, 0, []string{}, &result)
	return result
}

func generate(layerImages [][]string, layerIndex int, currentCombination []string, result *[][]string) {
	if layerIndex == len(layerImages) {
		if len(currentCombination) > 0 { // Ensure non-empty combinations are added.
			*result = append(*result, append([]string{}, currentCombination...))
		}
		return
	}

	// Option to skip the current layer to allow combinations of varying lengths.
	generate(layerImages, layerIndex+1, currentCombination, result)

	for _, img := range layerImages[layerIndex] {
		newCombination := append([]string{}, currentCombination...) // Clone currentCombination to avoid mutation.
		newCombination = append(newCombination, img)                // Add the new image.
		generate(layerImages, layerIndex+1, newCombination, result)
	}
}

func mergeImages(combination []string) image.Image {
	if len(combination) == 0 {
		fmt.Println("No images to merge.")
		return nil
	}

	var images []image.Image
	for _, path := range combination {
		file, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		img, err := png.Decode(file)
		if err != nil {
			panic(err)
		}
		images = append(images, img)
	}

	// Create a blank image to draw the layers on
	result := image.NewRGBA(image.Rect(0, 0, images[0].Bounds().Dx(), images[0].Bounds().Dy()))
	for _, img := range images {
		draw.Draw(result, result.Bounds(), img, image.Point{0, 0}, draw.Over)
	}

	return result
}

func saveImage(img image.Image, directory, filename string) {
	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		panic(err)
	}

	outFile, err := os.Create(filepath.Join(directory, filename))
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	err = png.Encode(outFile, img)
	if err != nil {
		panic(err)
	}
	fmt.Println("Image saved:", filepath.Join(directory, filename))
}
