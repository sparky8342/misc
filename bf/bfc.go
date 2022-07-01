package main

/*
A simple bf compiler - a go program is
generated from the bf, then shell out
and use the go compiler to create a binary

build this, then use: ./bfc [filename].bf
which will generate an executable [filename]
*/

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Token struct {
	cmd    rune
	amount int
}

var MEMSIZE int = 30000

func main() {
	program_file := os.Args[1]

	parts := strings.Split(program_file, ".")
	if len(parts) != 2 || parts[1] != "bf" {
		fmt.Println("usage: bfc [file].bf")
		os.Exit(0)
	}

	data, err := ioutil.ReadFile(program_file)
	if err != nil {
		panic(err)
	}

	program := string(data)
	tokens := get_tokens(program)
	if !validate(tokens) {
		fmt.Println("unbalanced brackets")
		os.Exit(0)
	}
	tokens = optimise(tokens)
	compile_program(tokens, parts[0])
}

func get_tokens(program string) []Token {
	tokens := []Token{}
	for _, b := range program {
		if b == '>' || b == '<' || b == '+' || b == '-' || b == '.' || b == ',' || b == '[' || b == ']' {
			tokens = append(tokens, Token{cmd: b, amount: 1})
		}
	}
	return tokens
}

func validate(tokens []Token) bool {
	brackets := 0
	for _, token := range tokens {
		if token.cmd == '[' {
			brackets++
		} else if token.cmd == ']' {
			brackets--
			if brackets < 0 {
				return false
			}
		}
	}
	return brackets == 0
}

func optimise(tokens []Token) []Token {
	new_tokens := []Token{}

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		if token.cmd == '+' || token.cmd == '-' || token.cmd == '<' || token.cmd == '>' {
			amount := 1
			for j := i + 1; j < len(tokens); j++ {
				if tokens[j].cmd != token.cmd {
					break
				} else {
					amount++
				}
			}
			new_tokens = append(new_tokens, Token{cmd: token.cmd, amount: amount})
			i += amount - 1
		} else {
			new_tokens = append(new_tokens, token)
		}
	}

	return new_tokens
}

func compile_program(tokens []Token, name string) {
	output := "package main\nimport (\n\"os\"\n\"fmt\"\n)\n"
	output += "func getin() int {\nb := make([]byte, 1)\nos.Stdin.Read(b)\nreturn int(b[0])}\n"
	output += "func main() {\nmem := make([]byte, " + strconv.Itoa(MEMSIZE) + ")\npos := 0\n"

	for _, token := range tokens {
		switch token.cmd {
		case '>':
			output += "pos+=" + strconv.Itoa(token.amount)
		case '<':
			output += "pos-=" + strconv.Itoa(token.amount)
		case '+':
			output += "mem[pos]+=" + strconv.Itoa(token.amount)
		case '-':
			output += "mem[pos]-=" + strconv.Itoa(token.amount)
		case '.':
			output += "fmt.Print(string(mem[pos]))"
		case ',':
			output += "mem[pos] = getin()"
		case '[':
			output += "for mem[pos] != 0 {"
		case ']':
			output += "}"
		}
		output += "\n"
	}
	output += "}\n"

	filename := "/tmp/" + name + ".go"

	err := ioutil.WriteFile(filename, []byte(output), 0644)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("bash", "-c", "go build "+filename)
	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	err = os.Remove(filename)
	if err != nil {
		panic(err)
	}
}
