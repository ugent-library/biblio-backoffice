// https://github.com/ugent-library/biblio-backoffice/issues/1718

describe("Issue #1718: Proxies: HTMX target error when adding a new proxy", () => {
  beforeEach(() => {
    cy.loginAsLibrarian("librarian1");
    cy.visit("/proxies");

    // Delete existing proxies if they exist
    cy.then(() => {
      const deleteButtons = Cypress.$('.btn:contains("Delete")');
      if (deleteButtons.length > 0) {
        cy.intercept("DELETE", "/proxies/*/people/*").as("deleteProxy");

        deleteButtons.each((_, button) => {
          button.click();

          cy.wait("@deleteProxy");
        });
      }
    });

    cy.contains("No proxies to display.").should("be.visible");
  });

  it("should be possible to select a proxy researcher after searching too specific", () => {
    cy.contains(".btn", "Add proxy").click();

    cy.ensureModal("Choose a proxy").within(() => {
      cy.setFieldByLabel("Search a proxy", "Researcher");

      cy.contains(".list-group-item", "Biblio Researcher1")
        .contains(".btn", "Choose proxy")
        .click();
    });

    cy.ensureModal("Select researchers").within(() => {
      cy.get(".list-group-item").should("have.length.greaterThan", 0);
      cy.contains("No researchers found").should("not.exist");

      cy.setFieldByLabel("Search researchers", "abc");
      cy.get(".list-group-item").should("have.length", 0);
      cy.contains("No researchers found").should("be.visible");

      cy.setFieldByLabel("Search researchers", "");
      cy.get(".list-group-item").should("have.length.greaterThan", 0);
      cy.contains("No researchers found").should("not.exist");

      cy.contains(
        "Select researchers from the left panel search results.",
      ).should("be.visible");

      cy.contains(".list-group-item", "John Doe")
        .contains(".btn", "Select researcher")
        .click();
      cy.contains(
        "Select researchers from the left panel search results.",
      ).should("not.exist");
      cy.get("#people").should("contain.text", "John Doe");

      cy.contains("#people .list-group-item", "John Doe")
        .contains(".btn", "Deselect")
        .click();
      cy.contains(
        "Select researchers from the left panel search results.",
      ).should("be.visible");

      cy.contains(".btn", "Done").click();
    });

    cy.ensureNoModal();

    cy.contains("No proxies to display.").should("be.visible");
  });
});
