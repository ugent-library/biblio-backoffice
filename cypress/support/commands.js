Cypress.Commands.addAll(
  {
    cacheAcrossSpecs: true,
  },
  {
    login(username, password) {
      cy.session(username, () => {
        // WARNING: Whenever you change the code of the session setup, Cypress will throw an error:
        //

        cy.request('/login')
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
      })
    },
  }
)

Cypress.Commands.add('loginAsResearcher', () => {
  cy.login(Cypress.env('RESEARCHER_USER_NAME'), Cypress.env('RESEARCHER_USER_PASSWORD'))
})

Cypress.Commands.add('loginAsLibrarian', () => {
  cy.login(Cypress.env('LIBRARIAN_USER_NAME'), Cypress.env('LIBRARIAN_USER_PASSWORD'))
})
