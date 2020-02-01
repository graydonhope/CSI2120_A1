package main

import (
	"fmt"
	"math"
	"os"
	"bufio"
	"strings"
	"strconv"
)

func sigmoid(inputs []float64, coefficients []float64) float64 {
	bias := coefficients[0]
	sum := bias
	for i := 0; i < len(inputs); i++ {
		sum += inputs[i] * coefficients[i + 1]
	}

	result := 1 / (1 + (math.Exp(-1 * sum)))
	return result
}

func main() {

	// Define Beta values inside main in case we want to accept values in from the console in the future.
	Betas := []float64{0.5, 0.3, 0.7, 0.1}
	var coefficients = [][]float64{ {0.1, 0.3, 0.4}, {0.5, 0.8, 0.3}, {0.7, 0.6, 0.6} }

	// Get N input from user
	reader := bufio.NewReader(os.Stdin)
	var input_N int

	for {
		fmt.Println("Please enter the desired value for N:")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		input = strings.TrimRight(input, "\r\n")
	
		N, err := strconv.Atoi(input)
	
		if err == nil {
			input_N = N
			break
		} else {
			// Invalid input
			fmt.Println("Invalid input for N. Try again")
		}
	}

	// Calculate neural network outputs from 0 to N-1
	for i:= 0; i < input_N; i++ {
		var X[]float64
		var calculated_Zs[]float64

		x1 := math.Sin( (2 * math.Pi * float64(i - 1)) / float64(input_N) )
		x2 := math.Cos( (2 * math.Pi * float64(i - 1)) / float64(input_N) )
		X = append(X, x1)
		X = append(X, x2)

		channel := make(chan float64, 3)
		write := make(chan bool)

		for j := 0; j < 3; j++ {
			if j > 0 {
				// block and wait for channel
				<- write
			} 
			go func(input int) {
				z := sigmoid(X, coefficients[input])
				channel <- float64(z)
				write <- true
			}(j)
		}
		
		for d := 0; d < 3; d++ {
			calculated_Zs = append(calculated_Zs, <- channel)
		}

		output := sigmoid(calculated_Zs, Betas)
		fmt.Println("Neural Network Output", (i + 1), ":", output)
	}
}