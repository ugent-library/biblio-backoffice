export default function loginAsResearcher(): void {
  cy.login(
    Cypress.env("RESEARCHER_USER_NAME"),
    Cypress.env("RESEARCHER_USER_PASSWORD"),
  );
}

declare global {
  namespace Cypress {
    interface Chainable {
      loginAsResearcher(): Chainable<void>;
    }
  }
}
