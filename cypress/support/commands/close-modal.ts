import { logCommand, updateConsoleProps } from './helpers'

const NO_LOG = { log: false }

type CloseModalOptions = {
  log?: boolean
}

export default function closeModal(
  subject: undefined | JQuery<HTMLElement>,
  save: boolean | string | RegExp = false,
  options: CloseModalOptions = {}
): void {
  const dismissButtonText = typeof save === 'boolean' ? (save ? 'Save' : 'Cancel') : save

  let log: Cypress.Log | null = null
  if (options.log === true) {
    log = logCommand('closeModal', { subject, 'Dismiss button text': dismissButtonText }, dismissButtonText)
    log.set('type', !subject ? 'parent' : 'child')
  }

  const doCloseModal = () => {
    cy.contains('.modal-footer .btn', dismissButtonText, NO_LOG)
      .then($el => {
        if (options.log === true) {
          log.set('$el', $el)
          updateConsoleProps(log, cp => {
            cp['Button element'] = $el
          })
        }

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
      closeModal(save: boolean | string | RegExp, options?: CloseModalOptions): Chainable<void>
    }
  }
}
