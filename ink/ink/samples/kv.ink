# basic key-value storage library built on composite values

create := () => (
  store := {}

  {
    store: store
    get: key => store.(key)
    set: (key, val) => store.(key) := val
    delete: key => store.(key) := ()
  }
)
