import { logCommand, updateConsoleProps } from './helpers'

export default function extractBiblioId(subject: undefined | JQuery<HTMLElement>, alias = 'biblioId') {
  const log = logCommand('extractBiblioId', { alias }, `@${alias}`)

  let cySubject: Cypress.Chainable
  if (subject) {
    if (subject.length !== 1) {
      expect(subject).to.have.length(1, `Expected subject to have length 1, but it has length ${subject.length}`)
    }

    cySubject = cy.wrap(subject, { log: false })
  } else {
    cySubject = cy.get('.list-group-item-main', { log: false })
  }

  cySubject
    .then(el => {
      updateConsoleProps(log, cp => (cp.subject = el))
    })
    .contains('Biblio ID:', { log: false })
    .find('.c-code', { log: false })
    .invoke({ log: false }, 'text')
    .as(alias, { type: 'static' })
    .finishLog(log, true)
}

declare global {
  namespace Cypress {
    interface Chainable {
      extractBiblioId(alias?: string): Chainable<string> | never
    }
  }
}
