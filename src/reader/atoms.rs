use nom::{
    IResult,
    bytes::complete::tag,
    branch::alt,
    combinator::{map, map_res},
};
use {
    crate::types::{Atom, Atom::{Bool, Symbol, String}},
    super::{numbers::{int, float}, string::string},
};

fn bool_parse(input: &str) -> IResult<&str, Atom> {
    map_res(
        alt((tag("false"), tag("true"))),
        |b: &str| b.parse::<bool>().map(Bool)
    )(input)
}

// fn symbol(input: &str) -> ParseResult<String> {
//     let mut next_index = 0;
//     let mut chars = input.chars();
//     match chars.next() {
//         Some(c) if c.is_alphabetic() => next_index += 1,
//         _ => return Err(input),
//     }
//     while let Some(c) = chars.next() {
//         if c.is_alphabetic()/*alphanumeric*/ || c == '-' {
//             next_index += 1;
//         } else {
//             break;
//         }
//     }
//     Ok((&input[next_index..], input[..next_index].to_string()))
// }
fn symbol(input: &str) -> IResult<&str, Atom> {
    map_res(
        tag("lol"),
        |s: &str| -> Result<Atom, &str> {Ok(Symbol(s.to_string()))}
    )(input)
}

pub fn atom(input: &str) -> IResult<&str, Atom> {
    alt((
        int,
        float,
        bool_parse,
        map(string, String),
        symbol
    ))(input)
}

#[cfg(test)]
mod tests {
    // use crate::{symbol, types::{Atom, Atom::{Bool, String}}};
    // use super::{bool_parse, symbol};

    // #[test]
    // fn parse_true() {
    //     assert_eq!(bool_parse("true"), Ok(("", Bool(true))));
    // }

    // #[test]
    // fn parse_false() {
    //     assert_eq!(bool_parse("false"), Ok(("", Bool(false))));
    // }

    // #[test]
    // fn parse_simple_symbol() {
    //     assert_eq!(symbol("abc"), Ok(("", symbol!("abc"))));
    // }

    // #[test]
    // fn test_symbol() {
    //     assert_eq!(
    //         symbol("i-am-an-symbol"),
    //         Ok(("", String("i-am-an-symbol".to_string()))),
    //     );
    //     assert_eq!(
    //         symbol("not entirely an symbol"),
    //         Ok((" entirely an symbol", String("not".to_string()))),
    //     );
    //     // assert_eq!(
    //     //     symbol("!not at all an symbol"),
    //     //     Err("!not at all an symbol"),
    //     // );
    // }
}
