from LiSH.datatypes import Symbol, Hashmap, is_atom


def escape(s):
    return s.replace("\\", "\\\\").replace('"', '\\"').replace("\n", "\\n")


def pr_str(exp, escape=escape):
    if exp == []:
        return "nil"
    elif isinstance(exp, Symbol):
        return exp
    elif isinstance(exp, Hashmap):
        return '{' + ' '.join(map(lambda kv: pr_str(kv[0], escape) + ' ' + pr_str(kv[1], escape), exp.items())) + '}'
    elif is_atom(exp):
        if isinstance(exp, str):
            return exp if escape is None else f'"{escape(exp)}"'
        elif isinstance(exp, bool):
            return "true" if exp else "false"
        if isinstance(exp, int) or isinstance(exp, float):
            return str(exp)
        else:
            return f"(atom {exp})"
    elif isinstance(exp, list) and exp[0] == Symbol("quote") and len(exp) == 1:
        return "(quote)"
    elif isinstance(exp, list) and exp[0] == Symbol("quote"):
        assert len(exp) == 2, f"Quote has zero or more than one argument: {exp}"
        return "'" + pr_str(exp[1])
    elif isinstance(exp, list) and exp[0] == Symbol("quasiquote"):
        assert len(exp) == 2, f"Quasiquote has zero or more than one argument: {exp}"
        return "'" + pr_str(exp[1])
    elif isinstance(exp, list):
        return "(" + " ".join(map(lambda x: pr_str(x, escape), exp)) + ")"
    elif exp is None:
        print("[FEAR AND LOATHING IN NONE VEGAS]")
    elif callable(exp):
        return str(exp)
    else:
        print("WTF IS THIS:", exp)
        return str(exp)


def PRINT(exp):
    """Convert an expression into a Lisp-readable string

        Args:
            exp: expression to pretty print

        Returns:
            pretty printed expression"""
    return pr_str(exp)
