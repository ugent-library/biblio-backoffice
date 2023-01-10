const NO_LOG = { log: false }

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

function logCommand(name, consoleProps = {}, message = '') {
  return Cypress.log({
    name,
    displayName: name
      .replace(/([A-Z])/g, ' $1')
      .trim()
      .toUpperCase(),
    message,
    consoleProps: () => consoleProps,
  })
}
