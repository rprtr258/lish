use crate::types::{LishRet};

pub fn print(val: &LishRet) -> String {
    match val {
        Ok(x) => format!("{:?}", x),
        Err(e) => format!("ERROR: {:?}", e),
    }
}
