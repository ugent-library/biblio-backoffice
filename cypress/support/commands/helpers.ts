export function logCommand(name, consoleProps = {}, message = '', $el = undefined) {
  return Cypress.log({
    $el,
    name,
    displayName: name
      .replace(/([A-Z])/g, ' $1')
      .trim()
      .toUpperCase(),
    message,
    consoleProps: () => consoleProps,
  })
}

export function updateLogMessage(log: Cypress.Log, append: unknown) {
  const message = log.get('message').split(', ')

  message.push(append)

  log.set('message', message.join(', '))
}

export function updateConsoleProps(log: Cypress.Log, callback: (ObjectLike) => void) {
  const consoleProps = log.get('consoleProps')()

  callback(consoleProps)

  log.set({ consoleProps: () => consoleProps })
}
