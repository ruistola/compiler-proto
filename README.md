# Tiny compiler front end prototype

A very limited scope compiler front end implemented as an exercise in Golang and as a self-study introduction to compilers. The language is essentially a statically compiled, strongly typed C- like language with LL(1) grammar, a Pratt parser for parsing operator heavy expressions, and a set of hand rolled statement parsing functions.

## Why tho

I am mainly publishing the source because this front end is very concise and light on abstractions, and should therefore make for reasonably good introductory reference material for programmers such as myself who may be experienced in general but just not familiar with the inner workings of compilers.

This project is entirely written in Go, so arguably any Grug Programmer should be able to understand the code, even with no prior Go experience. Give it a little elbow grease and a little spit-shine, and you might just make it mature enough to feed it into some compiler back end such as QBE or LLVM!

## Disclaimer

This is essentially abandonware because I got what I wanted out of it (a learning experience) and am working on a whole new language.

## Ok what am I supposed to do with it?

Run `go test .` in the repository root to execute all the tests that tokenize and parse the code samples in the `examples/` directory.

Run `go run .` to execute the main which loads `examples/program.jru` and fails because the source program contains an (intentional) error.
