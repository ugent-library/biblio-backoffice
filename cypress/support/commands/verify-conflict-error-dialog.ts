export default function verifyConflictErrorDialog(scope: Biblio.Scope) {
  return cy.ensureModal(null).within(() => {
    cy.contains(
      `${Cypress._.capitalize(scope)} has been modified by another user. Please reload the page.`,
    ).should("be.visible");
  });
}

declare global {
  namespace Cypress {
    interface Chainable {
      verifyConflictErrorDialog(
        scope: Biblio.Scope,
      ): Chainable<JQuery<HTMLElement>>;
    }
  }
}
