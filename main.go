package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/rossmacarthur/cases"
	"os"
	"slices"
	"strconv"
	"strings"
)

type FormatTargets []int

func (f *FormatTargets) String() string {
	return fmt.Sprintf("%v", *f)
}

func (f *FormatTargets) Set(value string) error {
	num, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	*f = append(*f, num)
	return nil
}

type Case int

const (
	Unknown Case = iota
	CamelCase
	PascalCase
	SnakeCase
	ScreamingSnakeCase
	KebabCase
	ScreamingKebabCase
	TrainCase
	TitleCase
	LowerCase
	UpperCase
)

func availableCases() []Case {
	return []Case{
		CamelCase, PascalCase, SnakeCase, ScreamingSnakeCase,
		KebabCase, ScreamingKebabCase, TrainCase,
		TitleCase, LowerCase, UpperCase,
	}
}

func (c Case) String() string {
	switch c {
	case Unknown:
		return "unknown"
	case CamelCase:
		return "CamelCase"
	case PascalCase:
		return "PascalCase"
	case SnakeCase:
		return "SnakeCase"
	case ScreamingSnakeCase:
		return "ScreamingSnakeCase"
	case KebabCase:
		return "KebabCase"
	case ScreamingKebabCase:
		return "ScreamingKebabCase"
	case TrainCase:
		return "TrainCase"
	case TitleCase:
		return "TitleCase"
	case LowerCase:
		return "LowerCase"
	case UpperCase:
		return "UpperCase"
	}
	return fmt.Sprintf("Case(%d)", int(c))
}

func (c Case) Apply(text string) (string, error) {
	switch c {
	case CamelCase:
		return cases.ToCamel(text), nil
	case PascalCase:
		return cases.ToPascal(text), nil
	case SnakeCase:
		return cases.ToSnake(text), nil
	case ScreamingSnakeCase:
		return cases.ToScreamingSnake(text), nil
	case KebabCase:
		return cases.ToKebab(text), nil
	case ScreamingKebabCase:
		return cases.ToScreamingKebab(text), nil
	case TrainCase:
		return cases.ToTrain(text), nil
	case TitleCase:
		return cases.ToTitle(text), nil
	case LowerCase:
		return cases.ToLower(text), nil
	case UpperCase:
		return cases.ToUpper(text), nil
	default:
		return "", fmt.Errorf("unknown case: %s", c)
	}
}

type CaseFlag struct {
	Value Case
}

func (cf *CaseFlag) String() string {
	return cf.Value.String()
}

func (cf *CaseFlag) Set(s string) error {
	s = strings.ToLower(s)
	for i := Case(0); i <= UpperCase; i++ {
		name := strings.ToLower(i.String())
		if name == s || strings.HasPrefix(name, s) {
			cf.Value = i
			return nil
		}
	}
	return fmt.Errorf("invalid case: %s", s)
}

func main() {
	d := flag.String("delim", ":", "Delimiter to use for delimiter")
	var fmtTargets FormatTargets
	flag.Var(&fmtTargets, "targets", "Comma-separated list of targets to format")
	var caseFlag CaseFlag
	caseFlag.Value = Unknown
	flag.Var(&caseFlag, "case", "case style")
	flag.Parse()

	if *d == "" || len(fmtTargets) == 0 || caseFlag.Value == Unknown {
		flag.Usage()
		if caseFlag.Value == Unknown {
			var sb strings.Builder
			sb.WriteString("Available values for 'case' is ...\n")
			for _, c := range availableCases() {
				sb.WriteString(fmt.Sprintf("\t%s\n", c.String()))
			}
			_, _ = fmt.Fprintln(os.Stderr, sb.String())
		}
		os.Exit(1)
	}

	slices.Sort(fmtTargets)
	convert := func(n int, line string) (string, error) {
		split := strings.Split(line, *d)
		length := len(split)
		r := make([]string, length)
		for i, w := range split {
			for _, target := range fmtTargets {
				if 1 <= target && target <= length && i == target-1 {
					c, err := caseFlag.Value.Apply(w)
					if err != nil {
						return "", fmt.Errorf("error at: line=%d col=%d word=%s case=%s, %w", n, target, w, caseFlag.Value, err)
					}
					r[i] = c
				} else {
					r[i] = w
				}
			}
		}
		return strings.Join(r, *d), nil
	}

	scanner := bufio.NewScanner(os.Stdin)
	current := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNum := current + 1
		result, err := convert(lineNum, line)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(result)
		current = lineNum
	}
	if err := scanner.Err(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(1)
	}
}
