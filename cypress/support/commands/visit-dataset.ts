import { logCommand, updateLogMessage, updateConsoleProps } from './helpers'

export default function visitDataset(): void {
  const log = logCommand('visitDataset')

  cy.get('@biblioId', { log: false }).then(biblioId => {
    updateLogMessage(log, biblioId)
    updateConsoleProps(log, cp => (cp['Biblio ID'] = biblioId))

    cy.intercept(`/dataset/${biblioId}/description*`).as('editDataset')

    cy.visit(`/dataset/${biblioId}`, { log: false })

    cy.wait('@editDataset', { log: false })
  })
}

declare global {
  namespace Cypress {
    interface Chainable {
      visitDataset(): Chainable<void>
    }
  }
}
