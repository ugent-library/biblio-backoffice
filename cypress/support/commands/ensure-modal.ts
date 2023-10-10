import { logCommand } from './helpers'

export default function ensureModal(expectedTitle: string | RegExp): Cypress.Chainable<JQuery<HTMLElement>> {
  const log = logCommand('ensureModal', { expectedTitle }, expectedTitle.toString())

  return cy
    .get('#modals', { log: false })
    .should('not.be.empty', { log: false })

    .within({ log: false }, () => {
      // Assertion "be.visible" doesn't work here because it is behind the dialog
      cy.get('#modal-backdrop', { log: false }).should('have.class', 'show')

      return cy
        .get('#modal', { log: false })
        .should('be.visible')
        .within({ log: false }, () => {
          if (expectedTitle instanceof RegExp) {
            cy.get('.modal-title', { log: false }).invoke({ log: false }, 'text').should('match', expectedTitle)
          } else {
            cy.get('.modal-title', { log: false }).should('have.text', expectedTitle)
          }
        })
    })
    .finishLog(log)
}

declare global {
  namespace Cypress {
    interface Chainable {
      ensureModal(expectedTitle: string | RegExp): Chainable<JQuery<HTMLElement>>
    }
  }
}
