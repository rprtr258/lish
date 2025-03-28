string : type
boolean : type
number : type
any : type

### System interfaces
# argv of the currently running process
args : () => []string
# Read from given file descriptor from some offset for some bytes, returned as a list of bytes (numbers).
read : (string, number, number, string => ())
# Read from stdin. The callback function returns a boolean that determines whether to continue reading from input.
in : ((string => ()) => boolean)
# Write to given file descriptor at some offset, some given bytes.
write : (string, number, string, {type: "end"} => ())
# print to stdout
out : string => ()
fs :: {
  # List the contents of a directory. Effectively `stat()` for all files in the directory.
  dir : (string, []{name: string, len: number, dir: boolean} => ())
  # Make a new directory at the given path.
  make : (string, () => ())
  # `stat` a file at a path, returning its canonicalized filename, size, and whether it's a directory or a file.
  stat : (string, {type: 'data', data: () | {name: string, len: number, dir: boolean, mod: number}} => ())
  # Delete some given file.
  delete : (string, () => ())
}
http :: {
  # Bind to a local TCP port and start handling HTTP requests.
  listen : (string, () => ()) => (() => ())
  # Send an HTTP client request.
  #`url` is required, `method`, `headers`, `body` are optional and default to their sensible zero values.
  req : ({url: string, method: string | (), headers: any | (), body: any | ()}, () => ()) => (() => ())
}
# Call the callback function after at least the given number of seconds has elapsed.
wait : (number, () => ())
# a pseudorandom floating point number in interval `[0, 1)`.
rand : () => number
# a string of given length containing random bits, safe for cryptography work
urand : (length: number) => string
# number of seconds in floating point in UNIX epoch.
time : () => number
# Exec the command at a given path with given arguments, with a given stdin, call given callback with stdout when exited.
exec : (command: string, args: []string, stdin: string, callback: (stdout: string) => ()) => (() => ())
# Exit the current process with the given exit code.
exit : number => ()

### Type casts and utilities (implemented as native functions)
# Convert type to string
string : any => string
# Convert type to number
number : any => number

# length of a list, string, or list-like composite value (equal to the number of keys on the composite or list value)
len : composite => number
# list of keys of the given composite
keys : composite => []string

# Take the first byte (i.e. ASCII value) of the string and return its numerical value
point : string => number
# reverse of `point()`. Note that behavior for values above 255 (full Unicode values) is undefined (so far).
char : number => string
