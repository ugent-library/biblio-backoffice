const PAGE_SIZE = 500

describe('Clean-up', { redirectionLimit: PAGE_SIZE }, () => {
  ;['publication', 'dataset'].forEach(type => {
    it(`should clean up ${type}s`, () => {
      cy.loginAsLibrarian()

      cy.switchMode('Librarian')

      const selector = 'button.dropdown-item[hx-get*="/confirm-delete"]:first'

      deleteFirstItem()

      function deleteFirstItem() {
        cy.visit(`/${type}`, {
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
                const id = button.attr('hx-get').match(/^\/(publication|dataset)\/(?<id>.*)\/confirm-delete/).groups.id

                cy.intercept({
                  method: 'DELETE',
                  url: `/${type}/${id}*`,
                }).as('delete-route')
              })

            //   Force is necessary because button is invisible at this point
            cy.get('@confirm-delete').click({ force: true })

            cy.contains('.modal-dialog .btn', 'Delete')
              .click()
              .then(() => {
                cy.wait('@delete-route')

                // Recursive call to delete other items
                deleteFirstItem()
              })
          } else {
            cy.get('.card-header').should('be.visible').should('contain.text', 'Showing 0')
          }
        })
      }
    })
  })
})
