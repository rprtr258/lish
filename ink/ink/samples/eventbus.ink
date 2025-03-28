{format} := import('str.ink')
{log: print} := import('logging.ink')
{filter, each} := import('functional.ink')

EventBus := () => (
  # event name to list of subscribers
  # {[string]: [any => ()]}
  subscribers := {}
  {
    # (string) => ()
    addEvent: (name) => subscribers.(name) :: {
      () -> subscribers.(name) := []
      _  -> print(format('Event {{.0}} is already defined', [name]))
    }
    removeEvent: (name) => subscribers.(name) := () # TODO: does it really delete by key?
    subscribe: (name, callback) => subscribers.(name) :: {
      () -> println(format('No such event "{{.0}}"', [name]))
      _ -> (
        evs := subscribers.(name)
        evs.(len(evs)) := callback
      )
    }
    unsubscribe: (name, callback) => subscribers.(name) :: {
      () -> println(format('No such event "{{.0}}"', [name]))
      _ -> subscribers.(name) == filter(subscribers.(name), (sub, _) => sub == callback)
    }
    emit: (name, payload) => subscribers.(name) :: {
      () -> () # no such event, skip
      _ -> each(subscribers.(name), (sub, _) => sub(payload))
    }
  }
)

# usage example:
Event := {
  Funds: {
    Added: 'funds_added'
    Removed: 'funds_removed'
  }
  Account: {Deleted: 'account_deleted'}
}

ebus := EventBus()

(ebus.addEvent)(Event.Funds.Added)
(ebus.addEvent)(Event.Funds.Removed)
(ebus.addEvent)(Event.Account.Deleted)
balance := {value: 55.5}

handleAccountDeleted := (_) => (
  print('OH NOOOOO')
)

handleFundsAdded := (payload) => (
  balance.value = balance.value + payload.amount
  print(balance.value)
)

handleFundsRemoved := (payload) => (
  balance.value = balance.value - payload.amount
  print(balance.value)
)

(ebus.subscribe)(Event.Funds.Added, handleFundsAdded)
(ebus.subscribe)(Event.Funds.Removed, handleFundsRemoved)
(ebus.subscribe)(Event.Account.Deleted, handleAccountDeleted)

(ebus.emit)(Event.Funds.Added, {amount: 20})
(ebus.emit)(Event.Funds.Removed, {amount: 13.46})
(ebus.emit)(Event.Account.Deleted, {})

# export
EventBus