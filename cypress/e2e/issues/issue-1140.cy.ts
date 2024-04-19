// https://github.com/ugent-library/biblio-backoffice/issues/1140

describe("Issue #1140: External contributor info is empty in the suggest box", () => {
  it("should display the external contributor name in the suggestions", () => {
    cy.loginAsResearcher();

    cy.setUpPublication("Book");

    cy.visitPublication();

    cy.updateFields(
      "Authors",
      () => {
        cy.intercept({
          pathname: "/publication/*/contributors/author/suggestions",
          query: {
            first_name: "Jane",
            last_name: "Doe",
          },
        }).as("suggestions");

        cy.setFieldByLabel("First name", "Jane");
        cy.setFieldByLabel("Last name", "Doe");

        cy.wait("@suggestions");

        cy.contains("#person-suggestions .list-group-item", "Jane Doe")
          .contains(".btn", "Add external author")
          .click();

        cy.setFieldByLabel("Roles", "Validation");
      },
      true,
    );

    cy.contains("table#contributors-author-table tr", "Jane Doe")
      .find(".if.if-edit")
      .click();

    cy.get("#person-suggestions")
      .find(".list-group-item")
      .should("have.length", 1)
      .find(".bc-avatar-text")
      .should("contain", "Current selection")
      .should("contain", "External, non-UGent")
      .should("contain", "Jane Doe");
  });
});
