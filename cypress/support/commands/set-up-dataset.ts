import { logCommand } from './helpers'

const NO_LOG = { log: false }

export default function setUpDataset(prepareForPublishing = false): void {
  logCommand('setUpDataset', {
    'Prepare for publishing': prepareForPublishing,
  })

  cy.visit('/dataset/add', NO_LOG)

  cy.intercept('/dataset/*/description*').as('completeDescription')

  cy.contains('Register a dataset manually', NO_LOG).find(':radio', NO_LOG).click(NO_LOG)
  cy.contains('.btn', 'Add dataset', NO_LOG).click(NO_LOG)

  // Extract biblioId at this point
  cy.get('#show-content', NO_LOG)
    .attr('hx-get')
    .then(hxGet => {
      const biblioId = hxGet.match(/\/dataset\/(?<biblioId>.*)\/description/)?.groups['biblioId']

      if (!biblioId) {
        throw new Error('Could not extract biblioId.')
      }

      return biblioId
    })
    .as('biblioId', { type: 'static' })

  cy.wait('@completeDescription', NO_LOG)

  cy.updateFields(
    'Dataset details',
    () => {
      cy.setFieldByLabel('Title', `The dataset title [CYPRESSTEST]`)

      if (prepareForPublishing) {
        cy.setFieldByLabel('Persistent identifier type', 'DOI')

        cy.setFieldByLabel('Identifier', '10.5072/test/t')
      }
    },
    true
  )

  cy.contains('.btn', 'Complete Description', NO_LOG).click(NO_LOG)
}

declare global {
  namespace Cypress {
    interface Chainable {
      setUpDataset(prepareForPublishing?: boolean): Chainable<void>
    }
  }
}
