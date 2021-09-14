use nom::{
    IResult,
    sequence::delimited,
    multi::fold_many0,
    combinator::{map, map_res},
    character::complete::{char, one_of, anychar},
};
use crate::types::{Atom, Atom::String};

fn character(input: &str) -> IResult<&str, char> {
    let (input, c) = one_of("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890 ()-_+=&'*,/:;<>?@^`{}~!\\")(input)?;
    if c == '\\' {
        map_res(anychar, |c| {
            Ok(match c {
            '"' | '\\' => c,
            'n' => '\n',
            'r' => '\r',
            't' => '\t',
            _ => return Err(()),
            })
        })(input)
    } else {
        Ok((input, c))
    }
}

pub fn string(input: &str) -> IResult<&str, Atom> {
    map(
        delimited(
            char('"'),
            fold_many0(character, std::string::String::new, |mut string, c| {
                string.push(c);
                string
            }),
            char('"'),
        ),
        String
    )(input)
}

#[cfg(test)]
mod tests {
    use crate::types::Atom;
    use super::string;

    mod string {
        use crate::types::Atom::String;
        use super::string;

        macro_rules! test_parse_string {
            ($($test_name:ident, $input:expr),* $(,)?) => {
                $(
                    #[test]
                    fn $test_name() {
                        assert_eq!(string(&format!(r#""{}""#, $input).to_owned()), Ok(("", String(format!("{}", $input)))))
                    }
                )*
            }
        }

        macro_rules! test_mirror_parse_string {
            ($($test_name:ident, $input:expr),* $(,)?) => {
                $(
                    #[test]
                    fn $test_name() {
                        assert_eq!(string(&format!("{:?}", $input).to_owned()), Ok(("", String(format!("{}", $input)))))
                    }
                )*
            }
        }

        test_mirror_parse_string!(
            mirror_doublequote, r#""1""#,
            slash_n, r#"\n"#,
            eight_backslashes, r#"\\\\"#,
            two_backslashes, r#"\"#,
            quote, r#"abc " def"#,
        );

        test_parse_string!(
            abc, "abc",
            with_parens, "abc (+ 1)",
            empty, "",
            ampersand, "&",
            singlequote, "'",
            openparen, "(",
            closeparen, ")",
            star, "*",
            plus, "+",
            comma, ",",
            minus, "-",
            slash, "/",
            colon, ":",
            semicolon, ";",
            less, "<",
            equal, "=",
            greater, ">",
            question, "?",
            dog, "@",
            caret, "^",
            underscore, "_",
            backquote, "`",
            opencurly, "{",
            closecurly, "}",
            tilde, "~",
            exclamation, "!",
        );
    }

    #[test]
    fn parse_string() {
        assert_eq!(string(r#""abc""#), Ok(("", Atom::from("abc"))));
    }
}
