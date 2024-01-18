import { logCommand } from './helpers'

export default function ensureToast(expectedTitle?: string | RegExp): Cypress.Chainable<JQuery<HTMLElement>> {
  const log = logCommand('ensureToast', { 'Expected title': expectedTitle }, expectedTitle)

  return cy
    .get('.toast', { log: false })
    .should(toast => {
      if (toast.length !== 1) {
        expect(toast).to.have.length(1)
      }

      if (expectedTitle) {
        if (typeof expectedTitle === 'string') {
          expect(toast).to.contain(expectedTitle)
        } else {
          expect(toast).to.match(expectedTitle)
        }
      }

      return toast
    })
    .should('be.visible')

    .within({ log: false }, () => {})
    .finishLog(log)
}

declare global {
  namespace Cypress {
    interface Chainable {
      ensureToast(expectedTitle?: string | RegExp): Chainable<JQuery<HTMLElement>>
    }
  }
}
