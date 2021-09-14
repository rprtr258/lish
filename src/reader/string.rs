use nom::{
    IResult,
    sequence::delimited,
    multi::fold_many0,
    combinator::map_res,
    character::complete::{char, one_of, anychar},
};

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

pub fn string(input: &str) -> IResult<&str, String> {
    delimited(
        char('"'),
        fold_many0(character, String::new, |mut string, c| {
            string.push(c);
            string
        }),
        char('"'),
    )(input)
}

#[cfg(test)]
mod tests {
    use super::string;

    macro_rules! test_string_formatted {
        ($format:expr, $($test_name:ident, $input:expr),* $(,)?) => {
            $(
                #[test]
                fn $test_name() {
                    assert_eq!(
                        string(&format!($format, $input)),
                        Ok(("", format!("{}", $input)))
                    )
                }
            )*
        }
    }

    test_string_formatted!("{:?}",
        mirror_doublequote, r#""1""#,
        slash_n, r#"\n"#,
        eight_backslashes, r#"\\\\"#,
        two_backslashes, r#"\"#,
        quote, r#"abc " def"#,
    );

    test_string_formatted!(r#""{}""#,
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
