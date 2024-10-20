# basic key-value storage library built on composite values

{
  new: () => (
    store := {}

    {
      store: store
      get: key => store.(key)
      set: (key, val) => store.(key) := val
      delete: key => store.(key) := ()
    }
  )
}