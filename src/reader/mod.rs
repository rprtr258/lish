use nom::IResult;
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
        types::{Atom, LishResult}
    }
};

pub fn lish(input: &str) -> IResult<&str, Atom> {
    // let list_content = interleave_parsers(
    //     spaces,
    //     lish
    // );
    // let list = left(
    //     right(
    //         option(
    //             match_str("(")
    //         ),
    //         list_content
    //     ),
    //     option(
    //         match_str(")")
    //     ),
    // ).map(|lst| list_vec!(lst));
    // TODO: fix error transform
    // alt((list, |input| atom(input).map_err(|e| "Error parsing atom")))
    // |_| Ok(("", Atom::Nil))
    unimplemented!()
}

// fn read_list(tokens: &mut Reader) -> LishResult {
//     let mut res = Vec::new();
//     loop {
//         match &tokens.peek()?[..] {
//         ")" => {
//             tokens.next().unwrap();
//             break
//         }
//         _ => res.push(read_form(tokens)?),
//         }
//     }
//     Ok(match res.len() {
//     0 => Atom::Nil,
//     _ => list_vec!(res),
//     })
// }

// TODO: reader macro
// fn read_form(tokens: &mut Reader) -> LishResult {
//     match &tokens.peek()?[..] {
//     "(" => {
//         tokens.next().unwrap();
//         read_list(tokens)
//     },
//     "'" => {
//         tokens.next().unwrap();
//         Ok(list_vec!(vec![symbol!("quote"), read_form(tokens)?]))
//     },
//     "`" => {
//         tokens.next().unwrap();
//         Ok(list_vec!(vec![symbol!("quasiquote"), read_form(tokens)?]))
//     },
//     "," => {
//         tokens.next().unwrap();
//         Ok(list_vec!(vec![symbol!("unquote"), read_form(tokens)?]))
//     },
//     ",@" => {
//         tokens.next().unwrap();
//         Ok(list_vec!(vec![symbol!("splice-unquote"), read_form(tokens)?]))
//     },
//     _ => Ok(read_atom(tokens.next()?)),
//     }
// }

// TODO: add braces implicitly
pub fn read(cmd: String) -> LishResult {
    // TODO:
    // lazy_static! {
        // static ref lish_parser = lish();
    // }
    let result = lish(cmd.as_str());
    match result {
    Ok((_, res)) => Ok(res),
    Err(s) => lisherr!(s)
    }
}

#[cfg(test)]
mod reader_tests {
    // use crate::{
    //     form,
    //     symbol,
    //     types::{Atom, Atom::Nil},
    // };
    // use super::read;

    // macro_rules! test_parse {
    //     ($($test_name:ident, $input:expr, $res:expr),* $(,)?) => {
    //         $(
    //             #[test]
    //             fn $test_name() {
    //                 assert_eq!(read($input.to_string()), Ok(Atom::from($res)))
    //             }
    //         )*
    //     }
    // }

    // TODO: parse_nothing, "", None,
    // test_parse!(
    //     num, "1", 1,
    //     num_spaces, "   7   ", 7,
    //     negative_num, "-12", -12,
    //     r#true, "true", true,
    //     r#false, "false", false,
    //     plus, "+", symbol!("+"),
    //     minus, "-", symbol!("-"),
    //     dash_abc, "-abc", symbol!("-abc"),
    //     dash_arrow, "->>", symbol!("->>"),
    //     abc, "abc", symbol!("abc"),
    //     abc_spaces, "   abc   ", symbol!("abc"),
    //     abc5, "abc5", symbol!("abc5"),
    //     abc_dash_def, "abc-def", symbol!("abc-def"),
    //     nil, "()", Nil,
    //     nil_spaces, "(   )", Nil,
    //     set, "(set a 2)", form![symbol!("set"), symbol!("a"), 2],
    //     list_nil, "(())", form![Nil],
    //     list_nil_2, "(()())", form![Nil, Nil],
    //     list_list, "((3 4))", form![form![3, 4]],
    //     list_inner, "(+ 1 (+ 3 4))", form![symbol!("+"), 1, form![symbol!("+"), 3, 4]],
    //     list_inner_spaces, "  ( +   1   (+   2 3   )   )  ", form![symbol!("+"), 1, form![symbol!("+"), 2, 3]],
    //     plus_expr, "(+ 1 2)", form![symbol!("+"), 1, 2],
    //     star_expr, "(* 1 2)", form![symbol!("*"), 1, 2],
    //     pow_expr, "(** 1 2)", form![symbol!("**"), 1, 2],
    //     star_negnum_expr, "(* -1 2)", form![symbol!("*"), -1, 2],
    //     string_spaces, r#"   "abc"   "#, "abc",
    //     reader_macro, "'(a b c)", form![symbol!("quote"), form![symbol!("a"), symbol!("b"), symbol!("c")]],
    //     comment, "123 ; such number", 123,
    // );
}
