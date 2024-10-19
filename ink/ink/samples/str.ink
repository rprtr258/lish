`` standard string library

std := load('functional')
map := std.map
reduce := std.reduce
slice := load('std').slice

`` checking if a given character is of a type
checkRange := (lo, hi) => c => (
  p := point(c)
  point(lo)-1 < p & p < point(hi)+1
)
upper? := checkRange('A', 'Z')
lower? := checkRange('a', 'z')
digit? := checkRange('0', '9')
letter? := c => upper?(c) | lower?(c)

`` is the char a whitespace?
ws? := c => point(c) :: {
  32 -> true `` space
  10 -> true `` newline
  9 -> true `` hard tab
  13 -> true `` carriage return
  _ -> false
}

`` hasPrefix? checks if a string begins with the given prefix substring
hasPrefix? := (s, prefix) => reduce(prefix, (acc, c, i) => acc & (s.(i) = c), true)

`` hasSuffix? checks if a string ends with the given suffix substring
hasSuffix? := (s, suffix) => (
  diff := len(s) - len(suffix)
  reduce(suffix, (acc, c, i) => acc & (s.(i + diff) = c), true)
)

`` mostly used for internal bookkeeping, matchesAt? reports if a string contains
`` the given substring at the given index idx.
matchesAt? := (s, substring, idx) => (
  max := len(substring)
  (sub := i => i :: {
    max -> true
    _ -> s.(idx + i) :: {
      substring.(i) -> sub(i + 1)
      _ -> false
    }
  })(0)
)

`` index is indexOf() for ink strings
index := (s, substring) => (
  max := len(s) - 1
  (sub := i => true :: {
    matchesAt?(s, substring, i) -> i
    i < max -> sub(i + 1)
    _ -> ~1
  })(0)
)

`` contains? checks if a string contains the given substring
contains? := (s, substring) => index(s, substring) > ~1

`` transforms given string to lowercase
lower := s => reduce(s, (acc, c, i) => upper?(c) :: {
  true -> acc.(i) := char(point(c) + 32)
  false -> acc.(i) := c
}, '')

`` transforms given string to uppercase
upper := s => reduce(s, (acc, c, i) => lower?(c) :: {
  true -> acc.(i) := char(point(c) - 32)
  false -> acc.(i) := c
}, '')

`` primitive "title-case" transformation, uppercases first letter and lowercases the rest.
title := s => (
  lowered := lower(s)
  lowered.0 := upper(lowered.0)
)

`` replace all occurrences of old substring with new substring in a string
replace := (s, old, new) => old :: {
  '' -> s
  _ -> (
    lold := len(old)
    lnew := len(new)
    (sub := (acc, i) => matchesAt?(acc, old, i) :: {
      true -> sub(
        slice(acc, 0, i) + new + slice(acc, i + lold, len(acc))
        i + lnew
      )
      false -> i < len(acc) :: {
        true -> sub(acc, i + 1)
        false -> acc
      }
    })(s, 0)
  )
}

`` convert string into list of characters
chars := s => map(s, c => c)

`` split given string into a list of substrings, splitting by the delimiter
split := (s, delim) => delim :: {
  '' -> chars(s)
  _ -> (
    coll := []
    ldelim := len(delim)
    (sub := (acc, i, last) => true :: {
      matchesAt?(acc, delim, i) -> (
        coll.len(coll) := slice(acc, last, i)
        sub(acc, i + ldelim, i + ldelim)
      )
      i < len(acc) -> sub(acc, i + 1, last)
      _ -> coll.len(coll) := slice(acc, last, len(acc))
    })(s, 0, 0)
  )
}

trimPrefixNonEmpty := (s, prefix) => (
  max := len(s)
  lpref := len(prefix)
  idx := (sub := i => i < max & matchesAt?(s, prefix, i) :: {
    true -> sub(i + lpref)
    _ -> i
  })(0)
  slice(s, idx, len(s))
)

`` trim string from start until it does not begin with prefix.
`` trimPrefix is more efficient than repeated application of
`` hasPrefix? because it minimizes copying.
trimPrefix := (s, prefix) => prefix :: {
  '' -> s
  _ -> trimPrefixNonEmpty(s, prefix)
}

trimSuffixNonEmpty := (s, suffix) => (
  lsuf := len(suffix)
  idx := (sub := i => i > 0 & matchesAt?(s, suffix, i - lsuf) :: {
    true -> sub(i - lsuf)
    _ -> i
  })(len(s))
  slice(s, 0, idx)
)

` trim string from end until it does not end with suffix.
  trimSuffix is more efficient than repeated application of
  hasSuffix? because it minimizes copying. `
trimSuffix := (s, suffix) => suffix :: {
  '' -> s
  _ -> trimSuffixNonEmpty(s, suffix)
}

`` trim string from both start and end with substring ss
trim := (s, ss) => trimPrefix(trimSuffix(s, ss), ss)

`` hexadecimal conversion utility functions
hToN := {0: 0, 1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 6, 7: 7, 8: 8, 9: 9, 'a': 10, 'b': 11, 'c': 12, 'd': 13, 'e': 14, 'f': 15}
nToH := '0123456789abcdef'

`` take number, return hex string
hex := n => (sub := (p, acc) => p < 16 :: {
  true -> nToH.(p) + acc
  false -> sub(floor(p / 16), nToH.(p % 16) + acc)
})(floor(n), '')

`` take hex string, return number
xeh := s => (
  `` i is the num of places from the left, 0-indexed
  max := len(s)
  (sub := (i, acc) => i :: {
    max -> acc
    _ -> sub(i + 1, acc * 16 + hToN.(s.(i)))
  })(0, 0)
)

`` tail recursive numeric list -> string converter
stringList := list => '[' + join(map(list, (x, _) => string(x)), ', ') + ']'

`` join a list of strings into a string
join := (list, joiner) => max := len(list) :: {
  0 -> ''
  _ -> (sub := (i, acc) => i :: {
    max -> acc
    _ -> sub(i + 1, acc.len(acc) := joiner + list.(i))
  })(1, list.0)
}

`` encode string buffer into a number list
encode := str => map(str, (c, _) => point(c))

`` decode number list into an ascii string
decode := data => reduce(data, (acc, cp) => acc.len(acc) := char(cp), '')

`` template formatting with {{ key }} constructs
format := (raw, values) => (
  `` parser state
  state := {
    `` current position in raw
    idx: 0
    `` parser internal state:
    `` 0 -> normal
    `` 1 -> seen one {
    `` 2 -> seen two {
    `` 3 -> seen a valid }
    which: 0
    `` buffer for currently reading key
    key: ''
    `` result build-up buffer
    buf: ''
  }

  `` helper function for appending to state.buf
  append := c => state.buf := state.buf + c

  `` read next token, update state
  readNext := () => (
    c := raw.(state.idx)

    state.which :: {
      0 -> c :: {
        '{' -> state.which := 1
        _ -> append(c)
      }
      1 -> c :: {
        '{' -> state.which := 2
        ` if it turns out that earlier brace was not
          a part of a format expansion, just backtrack `
        _ -> (
          append('{' + c)
          state.which := 0
        )
      }
      2 -> c :: {
        '}' -> (
          `` insert value
          append(string(state.key :: {
            '' -> values
            _ -> values.(state.key)
          }))
          state.key := ''
          state.which := 3
        )
        `` ignore spaces in keys -- not allowed
        ' ' -> ()
        _ -> state.key := state.key + c
      }
      3 -> c :: {
        '}' -> state.which := 0
        `` ignore invalid inputs -- treat them as nonexistent
        _ -> ()
      }
    }

    state.idx := state.idx + 1
  )

  `` main recursive sub-loop
  max := len(raw)
  (sub := () => state.idx < max :: {
    true -> (
      readNext()
      sub()
    )
    false -> state.buf
  })()
)