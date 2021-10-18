use nom::{
    IResult,
    error::ParseError,
    bytes::complete::tag,
    combinator::{opt, map, success, all_consuming},
    character::complete::{multispace0, char},
    multi::many0,
    branch::alt,
    sequence::{tuple, delimited},
};

mod string;
mod symbol;
mod numbers;
mod atoms;
use {
    atoms::atom,
    crate::{
        lisherr,
        list,
        form,
        symbol,
        types::{Atom, LishResult}
    }
};

fn space_delimited<'a, F, O, E>(inner: F) -> impl FnMut(&'a str) -> IResult<&'a str, O, E>
where
    F: 'a + FnMut(&'a str) -> IResult<&'a str, O, E>,
    E: ParseError<&'a str> {
    delimited(
        multispace0,
        inner,
        multispace0
    )
}

fn reader_macro<'a, F, E>(inner: F) -> impl FnMut(&'a str) -> IResult<&'a str, Atom, E>
where
    F: 'a + FnMut(&'a str) -> IResult<&'a str, Atom, E>,
    E: 'a + ParseError<&'a str> {
    map(space_delimited(
        tuple((
            opt(alt((
                tag("'"),
                tag("`"),
                tag(","),
                tag(",@"),
            ))),
            space_delimited(inner)
        ))
    ), |(reader_macro, atom)| match reader_macro {
        Some("'") => form![symbol!("quote"), atom],
        Some("`") => form![symbol!("quasiquote"), atom],
        Some(",") => form![symbol!("unquote"), atom],
        Some(",@") => form![symbol!("splice-unquote"), atom],
        None => atom,
        _  => unreachable!(),
    })
}

fn brackets_delimited<'a, O, F, E>(inner: F) -> impl FnMut(&'a str) -> IResult<&'a str, O, E>
where
    F: 'a + FnMut(&'a str) -> IResult<&'a str, O, E>,
    E: 'a + ParseError<&'a str> {
    delimited(
        char('('),
        inner,
        char(')')
    )
}

fn left_bracket_delimited<'a, O, F, E>(inner: F) -> impl FnMut(&'a str) -> IResult<&'a str, O, E>
where
    F: 'a + FnMut(&'a str) -> IResult<&'a str, O, E>,
    E: 'a + ParseError<&'a str> {
    delimited(
        char('('),
        inner,
        success(())
    )
}

fn right_bracket_delimited<'a, O, F, E>(inner: F) -> impl FnMut(&'a str) -> IResult<&'a str, O, E>
where
    F: 'a + FnMut(&'a str) -> IResult<&'a str, O, E>,
    E: 'a + ParseError<&'a str> {
    delimited(
        success(()),
        inner,
        char(')')
    )
}

fn left_outer_list_combine((left, mut inner): (Atom, Vec<Atom>)) -> Vec<Atom> {
    let mut res: Vec<Atom> = Vec::new();
    res.reserve(inner.len() + 1);
    res.push(left);
    res.append(&mut inner);
    res
}

fn right_outer_list_combine((mut inner, right): (Vec<Atom>, Atom)) -> Vec<Atom> {
    let mut res: Vec<Atom> = Vec::new();
    res.reserve(inner.len() + 1);
    res.append(&mut inner);
    res.push(right);
    res
}

fn left_right_outers_combine((left, mut inner, right): (Option<Atom>, Vec<Atom>, Option<Atom>)) -> Vec<Atom> {
    let mut res: Vec<Atom> = Vec::new();
    res.reserve(inner.len() + 2);
    left.map(|x| res.push(x));
    res.append(&mut inner);
    right.map(|x| res.push(x));
    res
}

fn inner_list(input: &str) -> IResult<&str, Atom> {
    reader_macro(map(brackets_delimited(
        space_delimited(many0(lish))
    ), |lst| list!(lst)))(input)
}

fn left_outer_list(input: &str) -> IResult<&str, Atom> {
    reader_macro(right_bracket_delimited(
        space_delimited(map(alt((
            many0(lish),
            map(tuple((
                left_outer_list,
                many0(lish)
            )), left_outer_list_combine)
        )), |lst| list!(lst)))
    ))(input)
}

fn right_outer_list(input: &str) -> IResult<&str, Atom> {
    reader_macro(left_bracket_delimited(
        space_delimited(map(tuple((
            many0(lish),
            opt(right_outer_list)
        )), |(inner, right)| list!(match right {
            None => inner,
            Some(x) => right_outer_list_combine((inner, x))
        })))
    ))(input)
}

fn outer_list(input: &str) -> IResult<&str, Atom> {
    reader_macro(map(
        alt((
            brackets_delimited(many0(lish)),
            left_bracket_delimited(map(space_delimited(tuple((
                many0(lish),
                opt(right_outer_list)
            ))), |(inner, right)| match right {
                Some(right_val) => right_outer_list_combine((inner, right_val)),
                None => inner
            })),
            all_consuming(right_bracket_delimited(space_delimited(many0(lish)))),
            right_bracket_delimited(
                map(space_delimited(tuple((
                    left_outer_list,
                    many0(lish),
                ))), left_outer_list_combine)
            ),
            map(space_delimited(tuple((
                opt(left_outer_list),
                many0(lish),
                opt(right_outer_list)
            ))), left_right_outers_combine)
        )), |lst| list!(lst)))(input)
}

fn lish(input: &str) -> IResult<&str, Atom> {
    reader_macro(space_delimited(
        alt((
            inner_list,
            atom
        ))
    ))(input)
}

pub fn read(cmd: String) -> LishResult {
    let result = outer_list(cmd.as_str());
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
    // #[test]
    // fn parse_nothing() {
    //     assert_eq!(read("".to_owned()), Ok(Atom::Nil))
    // }

    test_parse!(
        num, "1", form![1],
        num_spaces, "   7   ", form![7],
        negative_num, "-12", form![-12],
        r#true, "true", form![true],
        r#false, "false", form![false],
        plus, "+", form![symbol!("+")],
        minus, "-", form![symbol!("-")],
        dash_abc, "-abc", form![symbol!("-abc")],
        dash_arrow, "->>", form![symbol!("->>")],
        abc, "abc", form![symbol!("abc")],
        abc_spaces, "   abc   ", form![symbol!("abc")],
        abc5, "abc5", form![symbol!("abc5")],
        abc_dash_def, "abc-def", form![symbol!("abc-def")],
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
        string_spaces, r#"   "abc"   "#, form!["abc"],
        reader_macro, "'(a b c)", form![symbol!("quote"), form![symbol!("a"), symbol!("b"), symbol!("c")]],
        comment, "123 ; such number", form![123],
        string_arg_l, r#"(load-file "compose.lish""#, form![symbol!("load-file"), "compose.lish"],
        string_arg_r, r#"load-file "compose.lish")"#, form![symbol!("load-file"), "compose.lish"],
        right_outer_list_simple, "(+ 1 2", form![symbol!("+"), 1, 2],
        outer_list_simple, r#"echo 92"#, form![symbol!("echo"), 92],
        outer_plus, "+ 1 2", form![symbol!["+"], 1, 2],
        right_outer_twice, "(+ 1 2 (+ 3 4", form![symbol!["+"], 1, 2, form![symbol!["+"], 3, 4]],
        left_outer_twice, "+-curried 1) 3)", form![form![symbol!["+-curried"], 1], 3],
        outer_left_outer, "+-curried 1) 3", form![form![symbol!["+-curried"], 1], 3],
        outer_right_outer, "+ 1 2 (+ 3 4", form![symbol!["+"], 1, 2, form![symbol!["+"], 3, 4]],
    );
}
