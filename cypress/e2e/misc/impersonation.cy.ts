describe("User impersonation", () => {
  it("should be possible for a librarian to use the app as another user and to return to their own profile afterwards", () => {
    const LIBRARIAN_NAME = "Biblio Librarian1";

    cy.login("librarian1");

    cy.visit("/");

    cy.contains(".bc-avatar-and-text:visible", LIBRARIAN_NAME).click();
    cy.contains("View as").click();

    cy.ensureModal("View as other user").within(() => {
      cy.get("input[name=first_name]").should("be.focused");

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
