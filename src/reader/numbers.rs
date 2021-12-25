// #[cfg(test)]
// mod tests {
//     #[test]
//     fn zero() {
//         assert_eq!(int("0"), Ok(("", 0)));
//     }

//     #[test]
//     fn parse_positive_int() {
//         assert_eq!(int("123456"), Ok(("", 123456)));
//     }

//     #[test]
//     fn parse_negative_int() {
//         assert_eq!(int("-123456"), Ok(("", -123456)));
//     }

//     // #[test]
//     // fn parse_not_an_int() {
//     //     assert_eq!(
//     //         int("abc"),
//     //         Err(Err::Error(Error::new("abc", ErrorKind::Digit)))
//     //     );
//     // }

//     #[test]
//     fn parse_float() {
//         let (s, x) = float("0.1").unwrap();
//         assert_eq!(s, "");
//         assert!((x - 0.1).abs() < 1e-9);
//     }

//     #[test]
//     fn negative_float() {
//         let (s, x) = float("-0.1").unwrap();
//         assert_eq!(s, "");
//         assert!((x - -0.1).abs() < 1e-9);
//     }

//     #[test]
//     fn scientific_float() {
//         let (s, x) = float("3.2e-3").unwrap();
//         assert_eq!(s, "");
//         assert!((x - 3.2e-3).abs() < 1e-9);
//     }

//     #[test]
//     fn parse_not_a_float() {
//         assert_eq!(
//             float("a"),
//             Err(Err::Error(Error::new("a", ErrorKind::Float)))
//         );
//     }
// }
