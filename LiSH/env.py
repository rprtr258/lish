from LiSH.reader import Expression
from LiSH.datatypes import Symbol
from LiSH.core import ns


class NamedFunction:
    def __init__(self, name, body):
        self.name = name
        self.body = body

    def __call__(self, *args, **kwargs):
        return self.body(*args, **kwargs)

    def __repr__(self):
        return f"#{self.name}"


class Env(dict):
    "An environment: a dict of {'var':val} pairs, with an outer Env."
    def __init__(self, binds=[], exprs=[], outer=None):
        for var_name, var_value in zip(binds, exprs):
            self[var_name] = var_value
        self.outer = outer

    def set(self, var, value):
        env = self.find(var)
        if env is None:
            env = self
        env[var] = value

    def find(self, var: Symbol) -> Expression:
        """Find the innermost Env where var appears.

            Args:
                var: variable name to find

            Returns:
                innermost environment with such variable if found, None otherwise"""
        if var in self:
            return self
        if self.outer is None:
            return None  # nil
        return self.outer.find(var)

    def get(self, var):
        env = self.find(var)
        if env is None:
            return []
        return env[var]

    def __repr__(self):
        INDENT = "    "
        res = "{\n" + \
            f"{INDENT}" + \
            f"\n{INDENT}".join(f"{k}: {v}" for k, v in self.items() if not isinstance(v, NamedFunction)) + \
            "\n}"
        if self.outer is not None:
            res += " < " + repr(self.outer) + "\n"
        return res


def default_env():
    """An environment with some Scheme standard procedures.

        Returns:
            environment with default functions"""
    env = Env()
    # env.update(vars(math)) # sin, cos, sqrt, pi, ...
    env.update({
        Symbol(k): NamedFunction(k, v)
        for k, v in ns.items()})
    return env
