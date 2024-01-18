import { updateConsoleProps, updateLogMessage } from './helpers'

export default function finishLog(subject: unknown, log: Cypress.Log, appendToMessage = false) {
  if (log) {
    let theSubject = subject
    if (subject === null) {
      theSubject = '(null)'
    } else if (subject === '') {
      theSubject = '""'
    }

    updateConsoleProps(log, cp => (cp.yielded = theSubject))

    if (appendToMessage) {
      updateLogMessage(log, subject)
    }

    log.finish()
  }

  return subject
}

declare global {
  namespace Cypress {
    interface Chainable<Subject> {
      /**
       * Extends the log console props with the yielded result.
       *
       * @param Log log The log object to extend
       * @example
       * cy
       *   .validatedRequest(...)
       *   .finishLog(log)
       */
      finishLog(log: Log, appendToMessage?: boolean): Chainable<Subject>
    }
  }
}
