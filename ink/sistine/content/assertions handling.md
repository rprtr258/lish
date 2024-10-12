#someday_maybe #project

obviously should be turned on while developing

should they be turned on for release build? now i see 3 solutions:

- turn them off not to waste cpu time on them (gross)
- turn them on and waste cpu time (very gross)
- turn them off and handle carefully app fails (auto-restart?) (log program state and traceback into error log file)
    
    seems (perfect?) but might be difficult to implement