import { logCommand } from './helpers'

export default function ensureToast(expectedTitle: string | RegExp): Cypress.Chainable<JQuery<HTMLElement>> {
  const log = logCommand('ensureToast', { 'Expected title': expectedTitle }, expectedTitle.toString())

  return cy
    .contains('.toast', expectedTitle, { log: false })
    .should('be.visible')

    .within({ log: false }, () => {})
    .finishLog(log)
}

declare global {
  namespace Cypress {
    interface Chainable {
      ensureToast(expectedTitle: string | RegExp): Chainable<JQuery<HTMLElement>>
    }
  }
}
