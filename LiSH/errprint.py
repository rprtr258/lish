from sys import stderr


def errprint(*args, **kwargs):
    print(*args, file=stderr, **kwargs)
