import { logCommand } from './helpers'

export default function ensureNoModal(): void {
  logCommand('ensureNoModal')

  cy.get('#modals', { log: false }).children({ log: false }).should('have.length', 0)

  cy.get('#modal-backdrop, #modal', { log: false }).should('not.exist')
}

declare global {
  namespace Cypress {
    interface Chainable {
      ensureNoModal(): Chainable<void>
    }
  }
}
