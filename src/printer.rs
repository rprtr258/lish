use crate::types::{Atom};

pub fn print(cmd: &Atom) -> String {
    format!("{:?}", cmd)
}
