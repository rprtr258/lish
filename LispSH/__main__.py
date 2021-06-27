from sys import argv

from LispSH import repl, default_env

if __name__ == "__main__":
    repl_env = default_env()
    repl(repl_env)

