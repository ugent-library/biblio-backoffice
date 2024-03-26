// https://github.com/ugent-library/biblio-backoffice/issues/1108

describe("Issue #1108: Cannot add author without first name", () => {
  beforeEach(() => {
    cy.loginAsResearcher();

    cy.setUpPublication("Book");

    cy.visitPublication();
  });

  it("should be possible to add author without first name", () => {
    cy.updateFields(
      "Authors",
      () => {
        cy.setFieldByLabel("Last name", "Doe");

        cy.contains("Doe External, non-UGent")
          .closest(".list-group-item")
          .contains(".btn", "Add external author")
          .click();
      },
      true,
    );

    cy.get(".card#authors")
      .find("#contributors-author-table tr")
      .last()
      .find("td")
      .first()
      .find(".bc-avatar-text")
      .should("contain", "[missing] Doe");
  });

  it("should be possible to add author without last name", () => {
    cy.updateFields(
      "Authors",
      () => {
        cy.setFieldByLabel("First name", "John");

        cy.contains("John External, non-UGent")
          .closest(".list-group-item")
          .contains(".btn", "Add external author")
          .click();
      },
      true,
    );

    cy.get(".card#authors")
      .find("#contributors-author-table tr")
      .last()
      .find("td")
      .first()
      .find(".bc-avatar-text")
      .should("contain", "John [missing]");
  });
});
