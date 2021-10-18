use nom::{IResult, character::complete::one_of, multi::fold_many1};

pub fn symbol(input: &str) -> IResult<&str, String> {
    fold_many1(
        one_of("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-_+=&*/:<>?@^{}~!."),
        String::new,
        |mut string, c| {
            string.push(c);
            string
        }
    )(input)
}

#[cfg(test)]
mod tests {
    use nom::{Err, error::{Error, ErrorKind}};
    use super::symbol;

    macro_rules! test_symbol {
        ($($test_name:ident, $input:expr, $rest:expr, $result:expr),* $(,)?) => {
            $(
                #[test]
                fn $test_name() {
                    assert_eq!(
                        symbol($input),
                        Ok(($rest, $result.to_owned()))
                    );
                }
            )*
        }
    }

    test_symbol!(
        simple, "abc", "", "abc",
        dashes, "i-am-an-symbol", "", "i-am-an-symbol",
        dot_identifier, "compose.lish", "", "compose.lish",
        with_rest, "not entirely an symbol", " entirely an symbol", "not",
        with_exclamation, "!not at all an symbol", " at all an symbol", "!not",
    );

    #[test]
    fn empty_symbol() {
        assert_eq!(
            symbol(""),
            Err(Err::Error(Error::new("", ErrorKind::Many1)))
        );
    }

    #[test]
    fn not_symbol() {
        assert_eq!(
            symbol("\""),
            Err(Err::Error(Error::new("\"", ErrorKind::Many1)))
        );
    }
}
