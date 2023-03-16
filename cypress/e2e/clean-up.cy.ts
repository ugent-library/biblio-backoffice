const PAGE_SIZE = 500

describe('Clean-up', { redirectionLimit: PAGE_SIZE }, () => {
  it('clean up publications', () => {
    cy.loginAsLibrarian()

    cy.intercept('PUT', '/role/curator').as('curator')
    cy.visit('/')
    cy.contains('Researcher').click()
    cy.contains('Librarian').click()
    cy.wait('@curator')

    const selector = 'button.dropdown-item[hx-get*="/confirm-delete"]:first'

    deleteFirstPublication()

    function deleteFirstPublication() {
      cy.visit('/publication', {
        qs: {
          q: 'CYPRESSTEST',
          'page-size': PAGE_SIZE,
        },
      }).then(() => {
        // Make sure the button exists, otherwise test will fail when directly calling cy.get(selector)
        const $deleteButton = Cypress.$(selector)

        if ($deleteButton.length > 1) {
          throw new Error(`More than one delete button selected. Invalid selector "${selector}"`)
        }

        if ($deleteButton.length === 1) {
          cy.get(selector)
            .as('confirm-delete')
            .then(button => {
              const id = button.attr('hx-get').match(/^\/publication\/(?<id>.*)\/confirm-delete/).groups.id

              cy.intercept({
                method: 'DELETE',
                url: `/publication/${id}*`,
              }).as('delete-route')
            })

          //   Force is necessary because button is invisible at this point
          cy.get('@confirm-delete').click({ force: true })

          cy.contains('.modal-dialog .btn', 'Delete')
            .click()
            .then(() => {
              cy.wait('@delete-route')

              // Recursive call to delete other publications
              deleteFirstPublication()
            })
        } else {
          cy.get('.card-header')
            .should('be.visible')
            .should('contain.text', 'Publications')
            .should('contain.text', 'Showing 0')
        }
      })
    }
  })
})