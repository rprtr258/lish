use nom::{
    IResult,
    bytes::complete::tag,
    combinator::{opt, map},
    character::complete::char,
    multi::many0,
    branch::alt,
    sequence::{tuple, delimited},
};
mod string;
mod symbol;
mod numbers;
mod atoms;
mod utils;
use {
    atoms::atom,
    utils::spaces,
    crate::{
        lisherr,
        list_vec,
        form,
        symbol,
        types::{Atom, LishResult}
    }
};

fn lish(input: &str) -> IResult<&str, Atom> {
    let list = map(
        tuple((
            opt(alt((
                tag("'"),
                tag("`"),
                tag(","),
                tag(",@"),
            ))),
            delimited(
                delimited(spaces, char('('), spaces),
                many0(lish),
                opt(delimited(spaces, char(')'), spaces))
            )
        )),
        |(reader_macro, lst)| match reader_macro {
            Some("'") => form![symbol!("quote"), list_vec!(lst)],
            Some("`") => form![symbol!("quasiquote"), list_vec!(lst)],
            Some(",") => form![symbol!("unquote"), list_vec!(lst)],
            Some(",@") => form![symbol!("splice-unquote"), list_vec!(lst)],
            None => list_vec!(lst),
            Some(_)  => unreachable!(),
        }
    );
    delimited(
        spaces,
        alt((list, atom)),
        spaces
    )(input)
}

pub fn read(cmd: String) -> LishResult {
    let result = lish(cmd.as_str());
    match result {
        Ok((_, res)) => Ok(res),
        Err(s) => lisherr!(s)
    }
}

#[cfg(test)]
mod reader_tests {
    use crate::{
        form,
        symbol,
        types::Atom,
    };
    use super::read;

    macro_rules! test_parse {
        ($($test_name:ident, $input:expr, $res:expr),* $(,)?) => {
            $(
                #[test]
                fn $test_name() {
                    assert_eq!(read($input.to_owned()), Ok(Atom::from($res)))
                }
            )*
        }
    }

    // TODO: parse_nothing, "", Nil,
    test_parse!(
        num, "1", 1,
        num_spaces, "   7   ", 7,
        negative_num, "-12", -12,
        r#true, "true", true,
        r#false, "false", false,
        plus, "+", symbol!("+"),
        minus, "-", symbol!("-"),
        dash_abc, "-abc", symbol!("-abc"),
        dash_arrow, "->>", symbol!("->>"),
        abc, "abc", symbol!("abc"),
        abc_spaces, "   abc   ", symbol!("abc"),
        abc5, "abc5", symbol!("abc5"),
        abc_dash_def, "abc-def", symbol!("abc-def"),
        nil, "()", form![],
        nil_spaces, "(   )", form![],
        set, "(set a 2)", form![symbol!("set"), symbol!("a"), 2],
        list_nil, "(())", form![form![]],
        list_nil_2, "(()())", form![form![], form![]],
        list_list, "((3 4))", form![form![3, 4]],
        list_inner, "(+ 1 (+ 3 4))", form![symbol!("+"), 1, form![symbol!("+"), 3, 4]],
        list_inner_spaces, "  ( +   1   (+   2 3   )   )  ", form![symbol!("+"), 1, form![symbol!("+"), 2, 3]],
        plus_expr, "(+ 1 2)", form![symbol!("+"), 1, 2],
        star_expr, "(* 1 2)", form![symbol!("*"), 1, 2],
        pow_expr, "(** 1 2)", form![symbol!("**"), 1, 2],
        star_negnum_expr, "(* -1 2)", form![symbol!("*"), -1, 2],
        string_spaces, r#"   "abc"   "#, "abc",
        reader_macro, "'(a b c)", form![symbol!("quote"), form![symbol!("a"), symbol!("b"), symbol!("c")]],
        comment, "123 ; such number", 123,
        string_arg, r#"(load-file "compose.lish""#, form![symbol!("load-file"), "compose.lish"],
    );
}
