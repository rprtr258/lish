format := import('str.ink').format

Level := {
    DEBUG:    0
    INFO:     1
    WARNING:  2
    ERROR:    3
    CRITICAL: 4
}

`
  The Logger type returns a Logger object that can be configured to the user's needs

  List of variables that the formatter understands and can replace at runtime:

  "%name%": Logger name
  %message%: Message
  %time%: Current time in ms
  %datetime%: Current date time
  %level%: Log level string
  %levelnum%: Log level number
  %baselevel%: Base log level string
  %baselevelnum%: Base log level number
  %format%: Provided format string
  %hh%: Current hour
  %mm%: Current minute
  %ss%: Current second
  %time%: Current time
  %day%: Current day
  %monthname%: Current month name
  %year%: Current year
  %dayname%: Current day name
`

`
(string) => Logger
Logger := (name) => (
  name: string # name of the logger
  logLevel: number # current log level of the logger
  format: string # log message format
  fields: {[string]: any} # object to store custom log message formats
  setName: (name: string) => ()                 # Sets the name of the logger
  setBaseLogLevel: (baseLogLevel: number) => () # Sets the base log level of the logger
  setFilePath: (filePath: string) => ()         # Sets the file path for logging
  setFormat: (format: string) => ()             # Sets the log message format
  with: (key: string, value: any) => () # Adds a format key-value pair to the logger
  withFields: ({[string]: any}) => ()
  debug: (message: string) => ()    # Logs a message at debug level
  info: (message: string) => ()     # Logs a message at info level
  warn: (message: string) => ()     # Logs a message at warning level
  error: (message: string) => ()    # Logs a message at error level
  critical: (message: string) => () # Logs a message at critical level
)
`
Logger := (name) => (
  this := {
    name: name
    logLevel: Level.INFO
    format: '[{{level}}] {{message}}\n' # TODO: format func
    fields: {}
  }
  lg := (level, message) => (level < this.logLevel) :: {
    false -> out(format(this.format, {
      level: level :: {
        Level.DEBUG    -> 'DEBUG'
        Level.INFO     -> 'INFO'
        Level.WARNING  -> 'WARNING'
        Level.ERROR    -> 'ERROR'
        Level.CRITICAL -> 'CRITICAL'
      }
      message: message
      fields: this.fields
    }))
  }
  this.setName := name => this.name := name
  this.setLevel := level => this.logLevel := level
  this.setFormat := format => this.format := format
  this.with := (key, value) => this.fields.(key) := value
  this.withFields := fields => (import('functional.ink').each)(keys(fields), (key) => this.fields.(key) := dict.(value))
  this.debug    := message => lg(Level.DEBUG, message)
  this.info     := message => lg(Level.INFO, message)
  this.warn     := message => lg(Level.WARNING, message)
  this.error    := message => lg(Level.ERROR, message)
  this.critical := message => (
    lg(Level.CRITICAL, message)
    exit(1)
  )
  this
)

logger := Logger('root')

{
  Level: Level
  Logger: Logger
  logger: logger
  debug: logger.debug
  info: logger.info
  warn: logger.warn
  error: logger.error
  critical: logger.critical
  log: val => (logger.info)(string(val))
}
