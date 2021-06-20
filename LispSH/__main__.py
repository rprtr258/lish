from sys import argv

from LispSH import load_file, repl

if __name__ == "__main__":
    if len(argv) == 2:
        file_to_load = argv[1]
        load_file(file_to_load)
    elif len(argv) > 2:
        print("Too many arguments provided")
        exit(1)
    repl()

