// https://github.com/ugent-library/biblio-backoffice/issues/1309

describe('Issue #1309: Cannot return to profile after using "view as"', () => {
  it("should be possible for a researcher to return to their own profile when viewing as another user", () => {
    const LIBRARIAN_NAME = Cypress.env("LIBRARIAN_NAME");

    cy.loginAsLibrarian();

    cy.visit("/");

    cy.contains(".bc-avatar-and-text:visible", LIBRARIAN_NAME).click();

    cy.contains("View as").click();

    cy.ensureModal("View as other user").within(() => {
      cy.setFieldByLabel("First name", "John");
      cy.setFieldByLabel("Last name", "Doe");

      cy.contains("800000000001")
        .should("be.visible")
        .closest(".list-group-item")
        .contains(".btn", "Change user")
        .click();
    });

    cy.ensureNoModal();

    cy.contains("Viewing the perspective of John Doe.").should("be.visible");
    cy.get(".bc-avatar-and-text:visible").should("contain.text", "John Doe");

    cy.contains(".btn", `return to ${LIBRARIAN_NAME}`).click();

    cy.contains("Viewing the perspective of").should("not.exist");
    cy.get(".bc-avatar-and-text:visible").should(
      "contain.text",
      LIBRARIAN_NAME,
    );
  });
});
