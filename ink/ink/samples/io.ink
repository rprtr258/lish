# utility for reading an entire file
readFile := (path, cb) => (
  BufSize := 4096 # bytes
  (sub := (offset, acc) => read(path, offset, BufSize, evt => evt.type :: {
    'error' -> cb(())
    'data' -> (
      dataLen := len(evt.data)
      true :: {
        dataLen == BufSize -> sub(offset + dataLen, acc.len(acc) := evt.data)
        _ -> cb(acc.len(acc) := evt.data)
      }
    )
  }))(0, '')
)

# utility for writing an entire file
# it's not buffered, because it's simpler, but may cause jank later
# we'll address that if/when it becomes a performance issue
writeFile := (path, data, cb) => delete(path, evt => evt.type :: {
  # write() by itself will not truncate files that are too long,
  # so we delete the file and re-write. Not efficient, but writeFile
  # is not meant for large files
  'end' -> write(path, 0, data, evt => evt.type :: {
    'error' -> cb(())
    'end' -> cb(true)
  })
  _ -> cb(())
})

{readFile, writeFile}
