package log

import (
	"fmt"
	"log"
)

const (
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	reset  = "\033[0m"
)

func Println(in ...any) {
	log.Println(in...)
}

func Printf(format string, v ...any) {
	log.Printf(format, v...)
}

func Panic(in ...any) {
	log.Panic(in...)
}

func Blue(in ...any) {
	log.Println(blue + fmt.Sprint(in...) + reset)
}
func Yellow(in ...any) {
	log.Println(yellow + fmt.Sprint(in...) + reset)
}
func Green(in ...any) {
	log.Println(green + fmt.Sprint(in...) + reset)
}
func Red(in ...any) {
	log.Println(red + fmt.Sprint(in...) + reset)
}

// Info level log (blue)
func Info(op string, msg string) {
	log.Println(blue + fmt.Sprintf("%s:%s", op, msg) + reset)
}

// Success level log (green)
func Success(op string, msg string) {
	log.Println(green + fmt.Sprintf("%s:%s", op, msg) + reset)
}

// Warn level log (yellow)
func Warn(op string, msg string, err error) {
	if err == nil {
		log.Println(yellow + fmt.Sprintf("%s:%s", op, msg) + reset)
	} else {
		log.Println(yellow + fmt.Sprintf("%s:%s:%s", op, msg, err.Error()) + reset)
	}
}

// Error level log (red)
func Error(op string, msg string, err error) {
	if err == nil {
		log.Println(red + fmt.Sprintf("%s:%s", op, msg) + reset)
	} else {
		log.Println(red + fmt.Sprintf("%s:%s:%s", op, msg, err.Error()) + reset)
	}
}
