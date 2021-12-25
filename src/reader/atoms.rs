// use {
//     crate::types::{Atom, Atom::{Bool, Symbol, String, Int, Float}},
//     super::{numbers::{int, float}, string::string, symbol::symbol},
// };

// fn parse_bool(input: &str) -> IResult<&str, bool> {
//     map_res(
//         alt((tag("false"), tag("true"))),
//         |b: &str| b.parse::<bool>()
//     )(input)
// }

// pub fn atom(input: &str) -> IResult<&str, Atom> {
//     alt((
//         map(int, Int),
//         map(float, Float),
//         map(parse_bool, Bool),
//         map(string, String),
//         map(symbol, Symbol),
//     ))(input)
// }

// #[cfg(test)]
// mod tests {
//     use super::parse_bool;

//     #[test]
//     fn parse_true() {
//         assert_eq!(parse_bool("true"), Ok(("", true)));
//     }

//     #[test]
//     fn parse_false() {
//         assert_eq!(parse_bool("false"), Ok(("", false)));
//     }
// }
