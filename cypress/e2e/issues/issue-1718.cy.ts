// https://github.com/ugent-library/biblio-backoffice/issues/1718

describe("Issue #1718: Proxies: HTMX target error when adding a new proxy", () => {
  it("should be possible to select a proxy researcher after searching too specific", () => {
    cy.loginAsLibrarian("librarian1");
    cy.visit("/proxies");

    cy.contains(".btn", "Add proxy").click();

    cy.ensureModal("Choose a proxy").within(() => {
      cy.setFieldByLabel("Search a proxy", "Researcher");

      cy.contains(".list-group-item", "Biblio Researcher1")
        .contains(".btn", "Choose proxy")
        .click();
    });

    cy.ensureModal("Select researchers").within(() => {
      cy.get(".list-group-item").should("have.length.greaterThan", 0);

      cy.setFieldByLabel("Search researchers", "abc");
      cy.get(".list-group-item").should("have.length", 0);

      cy.setFieldByLabel("Search researchers", "");
      cy.get(".list-group-item").should("have.length.greaterThan", 0);
    });
  });
});
