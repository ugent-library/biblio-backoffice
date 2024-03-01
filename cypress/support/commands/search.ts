import { logCommand, updateLogMessage } from './helpers'

const NO_LOG = { log: false }

export default function search(query: string): Cypress.Chainable<number> {
  const log = logCommand('search', { query }, query)

  cy.get('input[placeholder="Search..."]', NO_LOG)
    .clear()
    .type(query, NO_LOG)
    .closest('.input-group', NO_LOG)
    .contains('.btn', 'Search', NO_LOG)
    .click(NO_LOG)

  return cy
    .contains('.card-header', '.pagination', NO_LOG)
    .find('.c-body-small', NO_LOG)
    .then($el => {
      const text = $el.text().trim()
      const regex = /^Showing (?<count>\d+)/

      expect(text).to.match(regex)

      const count = parseInt(text.match(regex).groups.count)

      updateLogMessage(log, count)

      return count
    })
    .finishLog(log)
}

declare global {
  namespace Cypress {
    interface Chainable {
      search(query: string): Cypress.Chainable<number>
    }
  }
}
