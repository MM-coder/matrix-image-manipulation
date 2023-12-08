package main

import (
	"fmt"
	"matrix-image-manipulation/manipulations"
	"matrix-image-manipulation/utils"
	"path/filepath"
	"strings"
)

func main() {
	var path, choice string

	// Request the file path from the user
	fmt.Print("Qual o path do ficheiro: ")
	_, err := fmt.Scanln(&path)
	if err != nil {
		fmt.Println("Error requesting input from user:", err)
		return
	}

	// Read the image into a matrix
	matrix, err := utils.ReadImageToMatrix(path)
	if err != nil {
		fmt.Println("Error reading image:", err)
		return
	}

	// Ask the user for the operation to perform
	fmt.Println("Escolha uma operação:")
	fmt.Println("1: Filtro Gaussiano")
	fmt.Println("2: Converter para Grayscale")
	fmt.Print("Escolha (1 ou 2): ")
	_, err = fmt.Scanln(&choice)
	if err != nil {
		fmt.Println("Erro a ler a escolha:", err)
		return
	}

	switch choice {
	case "1":
		matrix, err = manipulations.GaussianFilter(matrix, 7, 10.5)
		if err != nil {
			fmt.Println("Erro a aplicar filtro Gaussiano:", err)
			return
		}
	case "2":
		matrix, err = manipulations.ConvertToGreyScale(matrix)
		if err != nil {
			fmt.Println("Erro a converter para grayscale:", err)
			return
		}
	default:
		fmt.Println("Escolha inválida.")
		return
	}

	// Write the modified image back to a file
	outputPath := strings.TrimSuffix(path, filepath.Ext(path)) + "_new.png"
	err = utils.WriteImageFromMatrix(matrix, outputPath)
	if err != nil {
		fmt.Println("Error writing image:", err)
		return
	}

	fmt.Println("Operação completada. Output guardado em:", outputPath)
}
