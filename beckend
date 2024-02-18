package main

import (
	"Calculator/internal/storage/sqlite" // вот прямо сама разобралась как в базу записывать :)
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/scanner"
)

// ***********************
// вычислялочка
// ***********************

// Вот так она паникует
func Must(result string, err error) string {
	if err != nil {
		panic(err)
	}

	return result
}

// А это евал. просто евал...
func Eval(expr string) (string, error) {
	return EvalVars(expr, make(map[string]interface{}))
}

// пробежимся по типам
func EvalVars(expr string, vars map[string]interface{}) (string, error) {
	var s scanner.Scanner
	s.Init(strings.NewReader(expr))

	var r []string
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		token := s.TokenText()
		if v, ok := vars[token]; ok {
			switch v.(type) {
			case int:
				token = fmt.Sprintf("%d", v)
			case int64:
				token = fmt.Sprintf("%d", v)
			case float64:
				token = fmt.Sprintf("%f", v)
			default:
				return "", fmt.Errorf("unsupported var type: %T", v)
			}
		}

		r = append(r, token)
	}

	return evalTokens(r...)
}

// а вот что умеет вычислитель
func evalTokens(s ...string) (string, error) {
	var bodmas = []op{
		evalBrackets,
		binaryOp("^", func(a, b float64) float64 { return math.Pow(a, b) }),
		binaryOp("/", func(a, b float64) float64 { return a / b }),
		binaryOp("*", func(a, b float64) float64 { return a * b }),
		binaryOp("+", func(a, b float64) float64 { return a + b }),
		binaryOp("-", func(a, b float64) float64 { return a - b }),
		binaryOp("%", func(a, b float64) float64 { return float64(int(a) % int(b)) }),
	}

	var err error
	for _, op := range bodmas {
		s, err = op(s)
		if err != nil {
			return "", err
		}
	}

	return strings.Join(s, " "), nil
}

// снимаю шляпу перед теми, кто разобрался со скобками! А что, эта шляпа еще и скобки умеет?
func evalBrackets(s []string) ([]string, error) {
	for i := 0; i < len(s); i++ {
		if s[i] == "(" {
			bracketDepth := 0
			for j := i; j < len(s); j++ {
				switch s[j] {
				case "(":
					bracketDepth++
				case ")":
					bracketDepth--
				}

				if s[j] == ")" && bracketDepth == 0 {
					bracketResult, err := evalTokens(s[i+1 : j]...)
					if err != nil {
						return nil, err
					}

					s[i] = bracketResult
					s = append(s[0:i+1], s[j+1:]...)
					break
				}

				if j == len(s)-1 {
					return nil, errors.New("mismatched brackets, expected to find ')' but reached end of tokens")
				}
			}
		}
	}

	return s, nil
}

type op func([]string) ([]string, error)

func binaryOp(symbol string, fn func(float64, float64) float64) op {
	return func(s []string) ([]string, error) {
		for i := 0; i < len(s); i++ {
			if s[i] == symbol {
				lhs, err := strconv.ParseFloat(s[i-1], 64)
				if err != nil {
					return nil, fmt.Errorf("expected number got '%s': %s", s[i-1], err)
				}

				rhs, err := strconv.ParseFloat(s[i+1], 64)
				if err != nil {
					return nil, fmt.Errorf("expected number got '%s': %s", s[i+1], err)
				}

				s[i-1] = strconv.FormatFloat(fn(lhs, rhs), 'f', -1, 64)

				s = append(s[0:i], s[i+2:]...)
				i = i - 2
			}
		}

		return s, nil
	}
}

// ****************************************
// конец вычислялочки *********************
// ****************************************

func main() {
	// базу создаю, базеночку...
	storage, err := sqlite.New("sqliteDB")
	if err != nil {
		os.Exit(1)
	}

	// первая часть Марлезонского балета POST-метод
	// сохраняем выражение в базу и возвращаем ИД выражения
	http.HandleFunc("/calc/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")

		// вытягиваю параметры
		par := r.URL.Path[len("/calc/"):]

		// отправляю на расчет в вычислялочку
		myResult := Must(Eval(par))

		// сохраняю в базу результат
		retVal, err := storage.SaveURL(par, myResult)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Fprintf(w, "%d", retVal)
	})

	// вторая часть Марлезонского балета GET-метод для получения одного результата
	// возвращаем результат по ИД
	http.HandleFunc("/result/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")

		name := r.URL.Path[len("/result/"):]
		calc, err := storage.GetURL(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		strVar, _ := json.Marshal(calc)
		fmt.Fprintf(w, string(strVar))
	})

	// третья часть Марлезонского балета GET-метод для получения всех результатов
	// возвращаем все результаты
	http.HandleFunc("/results/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")

		calc, err := storage.GetAll()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		strVar, _ := json.Marshal(calc)
		fmt.Fprintf(w, string(strVar))
	})

	// стартуем сервис на порту
	http.ListenAndServe(":9990", nil)
}
