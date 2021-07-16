from LiSH.printer import PRINT


def FunctionCallError(proc, args, e):
    form = " ".join(map(PRINT, [proc] + args))
    return RuntimeError(f"""Error evaluating ({form}).
Error is: {e}""")
