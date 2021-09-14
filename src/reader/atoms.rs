use nom::{
    IResult,
    bytes::complete::tag,
    branch::alt,
    combinator::{map, map_res},
};
use {
    crate::types::{Atom, Atom::{Bool, Symbol, String}},
    super::{numbers::{int, float}, string::string, symbol::symbol},
};

fn bool_parse(input: &str) -> IResult<&str, Atom> {
    map_res(
        alt((tag("false"), tag("true"))),
        |b: &str| b.parse::<bool>().map(Bool)
    )(input)
}

pub fn atom(input: &str) -> IResult<&str, Atom> {
    alt((
        int,
        float,
        bool_parse,
        map(string, String),
        map(symbol, Symbol),
    ))(input)
}

#[cfg(test)]
mod tests {
    // use crate::{types::{Atom, Atom::{Bool, String}}};
    // use super::{bool_parse};

    // #[test]
    // fn parse_true() {
    //     assert_eq!(bool_parse("true"), Ok(("", Bool(true))));
    // }

    // #[test]
    // fn parse_false() {
    //     assert_eq!(bool_parse("false"), Ok(("", Bool(false))));
    // }
}
