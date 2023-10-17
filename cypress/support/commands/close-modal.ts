import { logCommand, updateConsoleProps } from './helpers'

const NO_LOG = { log: false }

export default function closeModal(subject: undefined | JQuery<HTMLElement>, save: boolean | string | RegExp = false) {
  const dismissButtonText = typeof save === 'boolean' ? (save ? 'Save' : 'Cancel') : save

  const log = logCommand('closeModal', { subject, 'Dismiss button text': dismissButtonText }, dismissButtonText)
  log.set('type', !subject ? 'parent' : 'child')

  const doCloseModal = () => {
    cy.contains('.modal-footer .btn', dismissButtonText, NO_LOG)
      .then($el => {
        log.set('$el', $el)
        updateConsoleProps(log, cp => {
          cp['Button element'] = $el
        })

        return $el
      })
      .click(NO_LOG)
  }

  if (subject) {
    cy.wrap(subject, NO_LOG).then(doCloseModal)
  } else {
    doCloseModal()
  }
}

declare global {
  namespace Cypress {
    interface Chainable {
      closeModal(save: boolean): Chainable<void>
      closeModal(dismissButtonText?: string | RegExp): Chainable<void>
    }
  }
}
