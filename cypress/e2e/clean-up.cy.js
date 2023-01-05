const PAGE_SIZE = 500

describe('Clean-up', { redirectionLimit: PAGE_SIZE }, () => {
  it('clean up publications', () => {
    cy.loginAsResearcher()

    const selector = 'button.dropdown-item[hx-get*="/confirm-delete"]:first'

    deleteFirstPublication()

    function deleteFirstPublication() {
      cy.visit(`/publication?q=&f[scope]=created&f[status]=private&page-size=${PAGE_SIZE}`).then(() => {
        // Make sure the button exists, otherwise test will fail when directly calling cy.get(selector)
        const $deleteButton = Cypress.$(selector)

        if ($deleteButton.length > 1) {
          throw new Error(`More than one delete button selected. Invalid selector "${selector}"`)
        }

        if ($deleteButton.length === 1) {
          cy.get(selector)
            .as('confirmDelete')
            .then(button => {
              const id = button.attr('hx-get').match(/^\/publication\/(?<id>.*)\/confirm-delete/).groups.id

              cy.intercept({
                method: 'DELETE',
                url: `/publication/${id}*`,
              }).as('deleteRoute')
            })

          //   Force is necessary because button is invisible at this point
          cy.get('@confirmDelete').click({ force: true })

          cy.contains('.modal-dialog button', 'Delete')
            .click()
            .then(() => {
              cy.wait('@deleteRoute')

              // Recursive call to delete other publications
              deleteFirstPublication()
            })
        }
      })
    }
  })
})
