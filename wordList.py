#!/usr/bin/env python3

# Usage: run on a file with newline-separated words to produce a Go file that
# defines the constant wordList, an array of strings, in the package wordenc.

from __future__ import print_function

if __name__ == "__main__":
    import argparse

    parser = argparse.ArgumentParser()
    parser.add_argument("file", type=argparse.FileType())
    args = parser.parse_args()

    print("package wordenc")
    print()
    print("var wordList = [...]string{")

    for line in args.file:
        word = line.strip()
        print('\t"{}",'.format(word))

    print("}")
