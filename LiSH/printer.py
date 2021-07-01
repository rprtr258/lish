from LiSH.datatypes import Symbol, Hashmap, is_atom

# TODO: WHY THIS EXISTS?
def pr_str_no_escape(exp):
    if exp == []:
        return "nil"
    elif isinstance(exp, Symbol):
        return exp
    elif isinstance(exp, Hashmap):
        return '{' + ' '.join(map(lambda kv: pr_str_no_escape(kv[0]) + ' ' + pr_str_no_escape(kv[1]), exp.items())) + '}'
    elif is_atom(exp):
        if isinstance(exp, str):
            return exp
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
        return "'" + pr_str_no_escape(exp[1])
    elif isinstance(exp, list) and exp[0] == Symbol("quasiquote"):
        assert len(exp) == 2, f"Quasiquote has zero or more than one argument: {exp}"
        return "'" + pr_str_no_escape(exp[1])
    elif isinstance(exp, list):
        return "(" + " ".join(map(pr_str_no_escape, exp)) + ")"
    elif exp is None:
        print("[FEAR AND LOATHING IN NONE VEGAS]")
    elif callable(exp):
        return str(exp)
    else:
        print("WTF IS THIS:", exp)
        return str(exp)

def pr_str(exp):
    def _escape(s):
        return s.replace("\\", "\\\\").replace('"', '\\"').replace("\n", "\\n")

    if exp == []:
        return "nil"
    elif isinstance(exp, Symbol):
        return exp
    elif isinstance(exp, Hashmap):
        return '{' + ' '.join(map(lambda kv: pr_str(kv[0]) + ' ' + pr_str(kv[1]), exp.items())) + '}'
    elif is_atom(exp):
        if isinstance(exp, str):
            return f'"{_escape(exp)}"'
        elif isinstance(exp, bool):
            return "true" if exp else "false"
        if isinstance(exp, int) or isinstance(exp, float):
            return str(exp)
        else:
            return f"(atom {exp})"
    elif isinstance(exp, list) and exp[0] == Symbol("quote") and len(exp) == 1:
        return "(quote)"
    elif isinstance(exp, list) and exp[0] == Symbol("quote"):
        assert len(exp) == 2, f"quote has zero or more than one argument: {exp}"
        return "'" + pr_str(exp[1])
    elif isinstance(exp, list) and exp[0] == Symbol("quasiquote"):
        assert len(exp) == 2, f"quasiquote has zero or more than one argument: {exp}"
        return "`" + pr_str(exp[1])
    elif isinstance(exp, list) and exp[0] == Symbol("unquote"):
        assert len(exp) == 2, f"unquote has zero or more than one argument: {exp}"
        return "~" + pr_str(exp[1])
    elif isinstance(exp, list) and exp[0] == Symbol("splice-unquote"):
        assert len(exp) == 2, f"splice-unquote has zero or more than one argument: {exp}"
        return "~@" + pr_str(exp[1])
    elif isinstance(exp, list):
        return "(" + " ".join(map(pr_str, exp)) + ")"
    elif exp is None:
        print("[FEAR AND LOATHING IN NONE VEGAS]")
    elif callable(exp):
        return str(exp)
    else:
        print("WTF IS THIS:", exp)
        return str(exp)

def PRINT(exp):
    "Convert an expression into a Lisp-readable string"
    return pr_str(exp)
