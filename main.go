package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Node represents a node in an expression tree.
// It can be a leaf (a number) or an internal node (an operation).
type Node struct {
	op    string // +, -, *, /
	value float64
	left  *Node
	right *Node
}

// Expression represents a single valid solution found.
type Expression struct {
	formula string
	value   float64
}

var operations = []string{"+", "-", "*", "/"}

func calculate(a, b float64, op string) (float64, bool) {
	switch op {
	case "+":
		return a + b, true
	case "-":
		return a - b, true
	case "*":
		return a * b, true
	case "/":
		if math.Abs(b) < 1e-9 {
			return 0, false // Avoid division by zero.
		}
		return a / b, true
	}
	return 0, false
}

func isApproximately24(value float64) bool {
	return math.Abs(value-24.0) < 1e-9
}

func generatePermutations(nums []float64) [][]float64 {
	if len(nums) <= 1 {
		return [][]float64{nums}
	}
	var result [][]float64
	for i, num := range nums {
		remaining := make([]float64, 0, len(nums)-1)
		remaining = append(remaining, nums[:i]...)
		remaining = append(remaining, nums[i+1:]...)
		for _, perm := range generatePermutations(remaining) {
			newPerm := append([]float64{num}, perm...)
			result = append(result, newPerm)
		}
	}
	return result
}

func generateOperations() [][]string {
	var result [][]string
	for _, op1 := range operations {
		for _, op2 := range operations {
			for _, op3 := range operations {
				result = append(result, []string{op1, op2, op3})
			}
		}
	}
	return result
}

// numStr is a helper to convert a float to a string for keys.
func numStr(n float64) string {
	return strconv.FormatFloat(n, 'g', -1, 64)
}

// collectOperands traverses chains of the same associative operator (like a + b + c)
// to flatten the structure for normalization.
func collectOperands(node *Node, op string, operands *[]string) {
	// If the child node is part of the same associative chain, recurse.
	if node.op == op {
		if node.left != nil {
			collectOperands(node.left, op, operands)
		}
		if node.right != nil {
			collectOperands(node.right, op, operands)
		}
	} else {
		// Otherwise, it's a new sub-expression, get its key.
		*operands = append(*operands, getCanonicalKey(node))
	}
}

// getCanonicalKey generates a unique, normalized string representation from an expression tree.
// This key ignores differences in operator order (commutativity) and grouping (associativity).
func getCanonicalKey(node *Node) string {
	// Base case: leaf node (a number)
	if node.left == nil && node.right == nil {
		return numStr(node.value)
	}

	// Recursive step: get keys for children
	keyL := getCanonicalKey(node.left)
	keyR := getCanonicalKey(node.right)

	// --- Normalization Rules ---

	// 1. Identity operations: simplify expressions with *1 or /1.
	if node.op == "*" {
		if keyL == "1" { return keyR }
		if keyR == "1" { return keyL }
	}
	if node.op == "/" && keyR == "1" {
		return keyL
	}

	// 2. Associativity & Commutativity: for + and *, flatten the expression,
	// sort the operands, and join them. This treats (a+b)+c and c+(a+b) as identical.
	if node.op == "+" || node.op == "*" {
		operands := []string{}
		collectOperands(node, node.op, &operands)
		sort.Strings(operands) // Sort for commutativity.
		return "(" + strings.Join(operands, node.op) + ")"
	}

	// 3. For non-commutative/associative operations (-, /), the order matters.
	return "(" + keyL + node.op + keyR + ")"
}

// findSolutions builds expression trees for all 5 parenthesis patterns,
// then generates a canonical key to find truly unique solutions.
func findSolutions(perm []float64, ops []string, seenKeys map[string]bool) []Expression {
	var results []Expression
	n := make([]*Node, 4)
	for i := 0; i < 4; i++ {
		n[i] = &Node{value: perm[i]}
	}
	op := ops

	formulas := map[int]string{
		1: "((%.0f %s %.0f) %s %.0f) %s %.0f",
		2: "(%.0f %s (%.0f %s %.0f)) %s %.0f",
		3: "%.0f %s ((%.0f %s %.0f) %s %.0f)",
		4: "%.0f %s (%.0f %s (%.0f %s %.0f))",
		5: "(%.0f %s %.0f) %s (%.0f %s %.0f)",
	}
	var trees []*Node

	// Pattern 1: ((a op b) op c) op d
	if v1, ok := calculate(n[0].value, n[1].value, op[0]); ok {
		if v2, ok := calculate(v1, n[2].value, op[1]); ok {
			if v3, ok := calculate(v2, n[3].value, op[2]); ok && isApproximately24(v3) {
				node1 := &Node{op: op[0], value: v1, left: n[0], right: n[1]}
				node2 := &Node{op: op[1], value: v2, left: node1, right: n[2]}
				trees = append(trees, &Node{op: op[2], value: v3, left: node2, right: n[3]})
			}
		}
	}
	// Pattern 2: (a op (b op c)) op d
	if v1, ok := calculate(n[1].value, n[2].value, op[1]); ok {
		if v2, ok := calculate(n[0].value, v1, op[0]); ok {
			if v3, ok := calculate(v2, n[3].value, op[2]); ok && isApproximately24(v3) {
				node1 := &Node{op: op[1], value: v1, left: n[1], right: n[2]}
				node2 := &Node{op: op[0], value: v2, left: n[0], right: node1}
				trees = append(trees, &Node{op: op[2], value: v3, left: node2, right: n[3]})
			}
		}
	}
	// Pattern 3: a op ((b op c) op d)
	if v1, ok := calculate(n[1].value, n[2].value, op[1]); ok {
		if v2, ok := calculate(v1, n[3].value, op[2]); ok {
			if v3, ok := calculate(n[0].value, v2, op[0]); ok && isApproximately24(v3) {
				node1 := &Node{op: op[1], value: v1, left: n[1], right: n[2]}
				node2 := &Node{op: op[2], value: v2, left: node1, right: n[3]}
				trees = append(trees, &Node{op: op[0], value: v3, left: n[0], right: node2})
			}
		}
	}
	// Pattern 4: a op (b op (c op d))
	if v1, ok := calculate(n[2].value, n[3].value, op[2]); ok {
		if v2, ok := calculate(n[1].value, v1, op[1]); ok {
			if v3, ok := calculate(n[0].value, v2, op[0]); ok && isApproximately24(v3) {
				node1 := &Node{op: op[2], value: v1, left: n[2], right: n[3]}
				node2 := &Node{op: op[1], value: v2, left: n[1], right: node1}
				trees = append(trees, &Node{op: op[0], value: v3, left: n[0], right: node2})
			}
		}
	}
	// Pattern 5: (a op b) op (c op d)
	if v1, ok1 := calculate(n[0].value, n[1].value, op[0]); ok1 {
		if v2, ok2 := calculate(n[2].value, n[3].value, op[2]); ok2 {
			if v3, ok3 := calculate(v1, v2, op[1]); ok3 && isApproximately24(v3) {
				node1 := &Node{op: op[0], value: v1, left: n[0], right: n[1]}
				node2 := &Node{op: op[2], value: v2, left: n[2], right: n[3]}
				trees = append(trees, &Node{op: op[1], value: v3, left: node1, right: node2})
			}
		}
	}

	for i, tree := range trees {
		key := getCanonicalKey(tree)
		if !seenKeys[key] {
			seenKeys[key] = true
			formula := fmt.Sprintf(formulas[i+1], perm[0], ops[0], perm[1], ops[1], perm[2], ops[2], perm[3])
			results = append(results, Expression{formula: formula, value: tree.value})
		}
	}
	return results
}

func parseInput(input string) ([]float64, error) {
	input = strings.TrimSpace(input)
	var parts []string
	if strings.Contains(input, ",") {
		parts = strings.Split(input, ",")
	} else if strings.Contains(input, " ") {
		parts = strings.Fields(input)
	} else if len(input) == 4 {
		parts = make([]string, 4)
		for i, char := range input {
			if char < '0' || char > '9' {
				return nil, fmt.Errorf("input must be numeric if no spaces/commas are used")
			}
			parts[i] = string(char)
		}
	} else {
		parts = strings.Fields(input)
	}
	if len(parts) != 4 {
		return nil, fmt.Errorf("you must enter exactly 4 numbers")
	}
	var nums []float64
	for _, part := range parts {
		part = strings.TrimSpace(part)
		num, err := strconv.ParseFloat(part, 64)
		if err != nil {
			return nil, fmt.Errorf("'%s' is not a valid number", part)
		}
		if num < 1 || num > 9 || num != math.Floor(num) {
			return nil, fmt.Errorf("numbers must be digits 1-9, found: %g", num)
		}
		nums = append(nums, num)
	}
	return nums, nil
}

func main() {
	fmt.Println("WELCOME TO THE 24 GAME SOLVER")
	fmt.Println("===============================")
	fmt.Println("Rules:")
	fmt.Println("- Enter 4 numbers (digits 1-9)")
	fmt.Println("- Format: 1 2 3 4 or 1,2,3,4 or 1234")
	fmt.Println("- The program will find all unique ways to make 24.")
	fmt.Println("- Supports: +, -, *, /")
	fmt.Println("===============================")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\nEnter 4 numbers (or 'quit' to exit): ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		if strings.ToLower(strings.TrimSpace(input)) == "quit" {
			fmt.Println("Thank you for playing!")
			break
		}
		nums, err := parseInput(input)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}
		fmt.Printf("\nSearching for solutions with: %.0f, %.0f, %.0f, %.0f\n", nums[0], nums[1], nums[2], nums[3])
		fmt.Println("===============================")

		var uniqueSolutions []Expression
		seenKeys := make(map[string]bool)
		permutations := generatePermutations(nums)
		operationCombos := generateOperations()

		for _, perm := range permutations {
			for _, ops := range operationCombos {
				solutions := findSolutions(perm, ops, seenKeys)
				uniqueSolutions = append(uniqueSolutions, solutions...)
			}
		}
		if len(uniqueSolutions) == 0 {
			fmt.Println("No solutions found for these numbers.")
		} else {
			fmt.Printf("Found %d unique solution(s):\n\n", len(uniqueSolutions))
			for i, solution := range uniqueSolutions {
				fmt.Printf("%d. %s = %.0f\n", i+1, solution.formula, solution.value)
			}
		}
		fmt.Println("\n===============================")
	}
}


