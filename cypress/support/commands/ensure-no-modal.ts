import { logCommand } from './helpers'

type EnsureNoModalOptions = {
  log?: boolean
}

export default function ensureNoModal(options: EnsureNoModalOptions = {}): void {
  if (options.log === true) {
    logCommand('ensureNoModal')
  }

  cy.get('#modals > *', { log: false })
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
      ensureNoModal(options?: EnsureNoModalOptions): Chainable<void>
    }
  }
}
