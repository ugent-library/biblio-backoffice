import { logCommand } from './helpers'

export default function ensureNoModal(): void {
  logCommand('ensureNoModal')

  cy.get('#modals *', { log: false })
    .should('have.length', 0)
    .then(() => {
      // Check before asserting to keep out of command log if ok
      if (Cypress.$('#modal, #modal-backdrop').length > 0) {
        cy.get('#modal').should('not.exist')
        cy.get('#modal-backdrop').should('not.exist')
      }
    })
}

declare global {
  namespace Cypress {
    interface Chainable {
      ensureNoModal(): Chainable<void>
    }
  }
}
