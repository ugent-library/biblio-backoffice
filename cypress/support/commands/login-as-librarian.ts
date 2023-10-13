export default function loginAsLibrarian(): void {
  cy.login(Cypress.env('LIBRARIAN_USER_NAME'), Cypress.env('LIBRARIAN_USER_PASSWORD'))
}

declare global {
  namespace Cypress {
    interface Chainable {
      loginAsLibrarian(): Chainable<void>
    }
  }
}
