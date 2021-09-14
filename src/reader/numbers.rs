use crate::types::{
    Atom,
    Atom::{Float, Int},
};
use nom::{
    character::complete::{digit1, char},
    combinator::{map_res, opt},
    number::complete::double,
    sequence::tuple,
    IResult,
};

pub fn int(input: &str) -> IResult<&str, Atom> {
    map_res(
        tuple((opt(char('-')), digit1)),
        |(sgn, s): (Option<char>, &str)| s
            .parse::<i64>()
            .map(|x| sgn.map_or(x, |_| -x))
            .map(Int)
    )(input)
}

pub fn float(input: &str) -> IResult<&str, Atom> {
    double(input).map(|(rest, x)| (rest, Float(x)))
}

#[cfg(test)]
mod tests {
    use nom::{
        error::{Error, ErrorKind},
        Err,
    };
    use super::{float, int};
    use crate::types::Atom::{Float, Int};

    #[test]
    fn zero() {
        assert_eq!(int("0"), Ok(("", Int(0))));
    }

    #[test]
    fn parse_positive_int() {
        assert_eq!(int("123456"), Ok(("", Int(123456))));
    }

    #[test]
    fn parse_negative_int() {
        assert_eq!(int("-123456"), Ok(("", Int(-123456))));
    }

    #[test]
    fn parse_not_an_int() {
        assert_eq!(
            int("abc"),
            Err(Err::Error(Error::new("abc", ErrorKind::Digit)))
        );
    }

    #[test]
    fn parse_float() {
        match float("0.1") {
            Ok(("", Float(x))) => assert!((x - 0.1).abs() < 1e-9),
            _ => unreachable!(),
        }
    }

    #[test]
    fn negative_float() {
        match float("-0.1") {
            Ok(("", Float(x))) => assert!((x - -0.1).abs() < 1e-9),
            _ => unreachable!(),
        }
    }

    #[test]
    fn scientific_float() {
        match float("3.2e-3") {
            Ok(("", Float(x))) => assert!((x - 3.2e-3).abs() < 1e-9),
            _ => unreachable!(),
        }
    }

    #[test]
    fn parse_not_a_float() {
        assert_eq!(
            float("a"),
            Err(Err::Error(Error::new("a", ErrorKind::Float)))
        );
    }
}
