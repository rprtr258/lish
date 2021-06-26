from LispSH.datatypes import Symbol, Atom

def pr_str(exp):
    def _escape(s):
        return s.replace("\\", "\\\\").replace('"', '\\"').replace("\n", "\\n")

    if exp == []:
        return "nil"
    elif isinstance(exp, Symbol):
        return exp
    elif isinstance(exp, Atom):
        if isinstance(exp.value, str):
            return f'"{_escape(exp.value)}"'
        if isinstance(exp.value, int) or isinstance(exp.value, float):
            return str(exp.value)
        elif isinstance(exp.value, bool):
            return "true" if exp.value else "false"
        else:
            return f"(atom {exp.value})"
    elif isinstance(exp, list) and exp[0] == Symbol("quote") and len(exp) == 1:
        return "(quote)"
    elif isinstance(exp, list) and exp[0] == Symbol("quote"):
        assert len(exp) == 2, f"Quote has zero or more than one argument: {exp}"
        return "'" + PRINT(exp[1])
    elif isinstance(exp, list):
        return "(" + " ".join(map(PRINT, exp)) + ")"
    elif exp is None:
        print("[FEAR AND LOATHING IN NONE VEGAS]")
    else:
        print("WTF IS THIS:", exp)
        return str(exp)

def PRINT(exp):
    "Convert an expression into a Lisp-readable string"
    return pr_str(exp)
