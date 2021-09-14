use nom::{character::complete::multispace1, combinator::map, IResult};

pub fn spaces(i: &str) -> IResult<&str, ()> {
    map(multispace1, |_| ())(i)
}

#[cfg(test)]
mod tests {
    use nom::{
        error::{Error, ErrorKind},
        Err,
    };
    use super::spaces;

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
            Err(Err::Error(Error::new("abc", ErrorKind::MultiSpace)))
        );
    }
}
