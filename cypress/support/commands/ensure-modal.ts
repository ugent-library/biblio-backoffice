import { logCommand } from './helpers'

export default function ensureModal(expectedTitle: string | RegExp): Cypress.Chainable<JQuery<HTMLElement>> {
  const log = logCommand('ensureModal', { 'Expected title': expectedTitle }, expectedTitle.toString())

  // Assertion "be.visible" doesn't work here because it is behind the dialog
  cy.get('#modal-backdrop', { log: false }).then(modalBackdrop => {
    if (!modalBackdrop.get(0).classList.contains('show')) {
      cy.wrap(modalBackdrop, { log: false }).should('have.class', 'show')
    }
  })

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
    .finishLog(log)
}

declare global {
  namespace Cypress {
    interface Chainable {
      ensureModal(expectedTitle: string | RegExp): Chainable<JQuery<HTMLElement>>
    }
  }
}
