use nom::{character::complete::multispace0, combinator::map, IResult};

pub fn spaces(input: &str) -> IResult<&str, ()> {
    map(multispace0, |_| ())(input)
}

#[cfg(test)]
mod tests {
    use super::spaces;

    #[test]
    fn zero_spaces() {
        assert_eq!(spaces(""), Ok(("", ())));
    }

    #[test]
    fn one_space() {
        assert_eq!(spaces(" "), Ok(("", ())));
    }

    #[test]
    fn many_spaces() {
        assert_eq!(spaces("         "), Ok(("", ())));
    }

    #[test]
    fn many_spaces_before_word() {
        assert_eq!(spaces("         abc"), Ok(("abc", ())));
    }

    #[test]
    fn spaces_and_newlines_before_word() {
        assert_eq!(spaces("   \n\n   \n   \nabc"), Ok(("abc", ())));
    }

    #[test]
    fn no_spaces_err() {
        assert_eq!(
            spaces("abc"),
            Ok(("abc", ()))
        );
    }
}
