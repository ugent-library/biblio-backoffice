// https://github.com/ugent-library/biblio-backoffice/issues/1125

describe('Issue #1125: Add "locked" message when record is locked', () => {
  it('should display "locked" message when record is locked', () => {
    cy.loginAsResearcher("researcher1");

    cy.setUpPublication("Book");

    // Lock the publication
    cy.loginAsLibrarian("librarian1");

    cy.visitPublication();

    cy.contains(".btn", "Lock record").click();

    // Verify the locked message
    cy.loginAsResearcher("researcher1");

    cy.visitPublication();

    // Give tooltip mechanism some time to load
    cy.wait(100);

    // Ensure publication is locked
    cy.get("#summary .bc-toolbar")
      .contains("Locked")
      .should("be.visible")
      .prev("i")
      .should("have.class", "if-lock")
      .click(); // Cypress doesn't have a hover command but click does the trick (as long as it doesn't navigate elsewhere)

    cy.contains(".tooltip", "Locked for editing").should("be.visible");

    // Assert the alert message
    // TODO: extract this to an ensureAlert command
    cy.get(".alert.alert-info")
      .should("be.visible")
      .within(() => {
        cy.get(".alert-title")
          .should("have.text", "This record has been reviewed and locked.")
          .next()
          .should(
            "have.text",
            "For any change requests or questions, get in touch via biblio@ugent.be. Thank you for your contribution!",
          );
      });
  });
});
