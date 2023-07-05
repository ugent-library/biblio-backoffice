const NO_LOG = { log: false }

declare namespace Cypress {
  interface Chainable<Subject> {
    login(username: string, password: string): Chainable<void>

    loginAsResearcher(): Chainable<void>

    loginAsLibrarian(): Chainable<void>

    switchMode(mode: 'Researcher' | 'Librarian'): Chainable<void>

    ensureModal(expectedTitle: string, strict?: boolean): Chainable<JQuery<HTMLElement>>

    ensureNoModal(): Chainable<void>

    extractBiblioId(alias?: string): Chainable<string> | never

    /**
     * Extends the log console props with the yielded result.
     *
     * @param Log log The log object to extend
     * @example
     * cy
     *   .validatedRequest(...)
     *   .finishLog(log)
     */
    finishLog(log: Log, appendToMessage?: boolean): Chainable<Subject>
  }
}

// Parent commands
Cypress.Commands.addAll({
  login(username, password) {
    // WARNING: Whenever you change the code of the session setup, Cypress will throw an error:
    //   This session already exists. You may not create a new session with a previously used identifier.
    //   If you want to create a new session with a different setup function, please call cy.session() with
    //   a unique identifier other than...
    //
    // Temporarily uncomment the following line to clear the sessions if this happens
    // Cypress.session.clearAllSavedSessions()

    logCommand('login', { username }, username)

    cy.session(
      username,
      () => {
        cy.request('/login', NO_LOG)
          .then(response => {
            const action = response.body.match(/action\=\"(.*)\" /)[1]

            return action.replace(/&amp;/g, '&')
          })
          .then(actionUrl =>
            cy.request({
              method: 'POST',
              url: actionUrl,
              form: true,

              body: {
                username,
                password,
              },

              // Make sure we redirect back and get the cookie we need
              followRedirect: true,

              // Make sure we don't leak passwords in the Cypress log
              log: false,
            })
          )
      },
      {
        cacheAcrossSpecs: true,
      }
    )
  },

  loginAsResearcher() {
    cy.login(Cypress.env('RESEARCHER_USER_NAME'), Cypress.env('RESEARCHER_USER_PASSWORD'))
  },

  loginAsLibrarian() {
    cy.login(Cypress.env('LIBRARIAN_USER_NAME'), Cypress.env('LIBRARIAN_USER_PASSWORD'))
  },

  switchMode(mode: 'Researcher' | 'Librarian') {
    const currentMode = mode === 'Researcher' ? 'Librarian' : 'Researcher'

    let log: Cypress.Log

    cy.visit('/', NO_LOG)

    cy.intercept({ method: 'PUT', url: '/role/*' }, NO_LOG).as('switch-role')

    cy.contains(`.c-sidebar > .dropdown > button`, currentMode, NO_LOG)
      .click(NO_LOG)
      .next('.dropdown-menu', NO_LOG)
      .contains(mode, NO_LOG)
      .then($el => {
        log = logCommand('switchMode', { 'Current mode': currentMode, 'New mode': mode }, mode, $el)
        log.snapshot('before')
      })
      .click(NO_LOG)

    cy.wait('@switch-role', NO_LOG)

    cy.contains(`.c-sidebar > .dropdown > button`, mode, NO_LOG).then($el => {
      log.set('$el', $el).snapshot('after').end()
    })
  },

  ensureModal(expectedTitle, strict = true) {
    logCommand('ensureModal', { expectedTitle, strict }, expectedTitle)

    cy.get('#modals', NO_LOG)
      .should('not.be.empty', NO_LOG)

      .within(NO_LOG, () => {
        // Assertion "be.visible" doesn't work here because it is behind the dialog
        cy.get('#modal-backdrop', NO_LOG).should('have.class', 'show')

        cy.get('#modal', NO_LOG)
          .should('be.visible')
          .within(NO_LOG, () => {
            cy.get('.modal-title', NO_LOG).should(strict ? 'have.text' : 'contain.text', expectedTitle)
          })
      })

    // Yield the #modal dialog element
    cy.get('#modal .modal-dialog', NO_LOG)
  },

  ensureNoModal() {
    logCommand('ensureNoModal')

    cy.get('#modals', NO_LOG).children(NO_LOG).should('have.length', 0)

    cy.get('#modal-backdrop, #modal', NO_LOG).should('not.exist')
  },
})

// Child commands
Cypress.Commands.addAll(
  { prevSubject: true },
  {
    extractBiblioId(subject, alias = 'biblioId') {
      const log = logCommand('extractBiblioId', { subject, alias }, `@${alias}`)

      if (subject.length !== 1) {
        expect(subject).to.have.length(1, `Expected subject to have length 1, but it has length ${subject.length}`)
      }

      cy.wrap(subject, NO_LOG)
        .contains('Biblio ID:', NO_LOG)
        .find('.c-code', NO_LOG)
        .invoke(NO_LOG, 'text')
        .as(alias, { type: 'static' })
        .finishLog(log, true)
    },

    finishLog(subject, log, appendToMessage = false) {
      let theSubject = subject
      if (subject === null) {
        theSubject = '(null)'
      } else if (subject === '') {
        theSubject = '""'
      }

      updateConsoleProps(log, cp => (cp.yielded = theSubject))

      if (appendToMessage) {
        updateLogMessage(log, subject)
      }

      log.end()

      return subject
    },
  }
)

function logCommand(name, consoleProps = {}, message = '', $el = undefined) {
  return Cypress.log({
    $el,
    name,
    displayName: name
      .replace(/([A-Z])/g, ' $1')
      .trim()
      .toUpperCase(),
    message,
    consoleProps: () => consoleProps,
  })
}

function updateLogMessage(log: Cypress.Log, subject: unknown) {
  const message = log.get('message').split(', ')

  message.push(subject)

  log.set('message', message.join(', '))
}


function updateConsoleProps(log: Cypress.Log, callback: (ObjectLike) => void) {
  const consoleProps = log.get('consoleProps')()

  callback(consoleProps)

  log.set({ consoleProps: () => consoleProps })
}
