from sys import argv

from LispSH import load_file, repl, default_env

if __name__ == "__main__":
    repl_env = default_env()
    if len(argv) == 2:
        file_to_load = argv[1]
        load_file(file_to_load, repl_env)
    elif len(argv) > 2:
        print("Too many arguments provided")
        exit(1)
    repl(repl_env)

