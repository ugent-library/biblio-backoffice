import { logCommand, updateLogMessage, updateConsoleProps } from './helpers'

export default function visitPublication(alias = '@biblioId'): Cypress.Chainable<Cypress.AUTWindow> {
  const log = logCommand('visitPublication', { alias }, alias)

  return cy.get(alias, { log: false }).then(biblioId => {
    updateLogMessage(log, biblioId)
    updateConsoleProps(log, cp => (cp['Biblio ID'] = biblioId))

    return cy.visit(`/publication/${biblioId}`, { log: false })
  })
}

declare global {
  namespace Cypress {
    interface Chainable {
      visitPublication(alias?: string): Chainable<AUTWindow>
    }
  }
}
