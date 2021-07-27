#[derive(Clone)]
pub enum Atom {
    Bool(bool),
    Int(i64),
    Float(f64),
    String(String),
    Symbol(String),
}

impl std::fmt::Display for Atom {
    fn fmt(&self, fmt: &mut std::fmt::Formatter) -> std::fmt::Result {
        match self {
            Atom::Bool(b) => write!(fmt, "{}", b),
            Atom::Int(i) => write!(fmt, "{}", i),
            Atom::Float(f) => write!(fmt, "{}", f),
            Atom::String(s) => write!(fmt, "\"{}\"", s),
            Atom::Symbol(s) => write!(fmt, "{}", s),
        }
    }
}

#[derive(Debug, Clone)]
pub enum Form {
    List(Vec<Form>),
    Atom(Atom),
}
