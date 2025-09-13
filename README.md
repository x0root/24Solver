# 24 Game Solver  
[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://go.dev/) [![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

A simple solver for the classic **24 Game**.  
It takes 4 digits (1–9) and tries to make the number **24** using the basic math operators: `+ - * /`.

## What is the 24 Game?

The 24 Game is a math puzzle played with 4 numbers.  
Rules:
- Use **all 4 numbers exactly once**.  
- You can use the operations `+ - * /`.  
- Parentheses can be added anywhere to control the order of operations.  
- The goal: make the result equal **24**.  

Example: with numbers `3, 3, 8, 8` → one valid solution is `(8 / (3 - 8/3)) = 24`.

This program automates the search for all possible valid solutions.

## Notice
While the solver removes obvious duplicates, some results may still look similar because of different parenthesis placements or mathematically equivalent expressions.

## Run the Program

1. Make sure [Go](https://go.dev/dl/) is installed.  
2. Clone this repository:
   ```bash
   git clone https://github.com/x0root/24Solver.git
   cd 24Solver
   ```
3. Run the solver:
   ```bash
   go run main.go
   ```
4. Enter 4 digits (example: `1 2 3 4` or `1234`), and the program will search for all valid solutions.

## Contributing

Contributions are welcome. You can help with:  
- Bug fixes  
- Code optimization  
- New features (e.g. more operators, GUI, or testing)  

Fork the repo, create a branch, and open a pull request. Suggestions and discussions are also encouraged.  

## License
This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.
