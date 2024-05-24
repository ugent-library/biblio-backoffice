// https://github.com/ugent-library/biblio-backoffice/issues/1108

describe("Issue #1108: Cannot add an external author without first name", () => {
  beforeEach(() => {
    cy.loginAsResearcher();
  });

  describe("for publications", () => {
    beforeEach(() => {
      cy.setUpPublication("Book");
      cy.visitPublication();
    });

    it("should be possible to add author without first name", () => {
      cy.updateFields(
        "Authors",
        () => {
          cy.setFieldByLabel("Last name", "Dow");

          cy.contains("Dow External, non-UGent")
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
        .should("contain", "[missing] Dow");
    });

    it("should be possible to add author without last name", () => {
      cy.updateFields(
        "Authors",
        () => {
          cy.setFieldByLabel("First name", "Jane");

          cy.contains("Jane External, non-UGent")
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
        .should("contain", "Jane [missing]");
    });
  });

  describe("for datasets", () => {
    beforeEach(() => {
      cy.setUpDataset();
      cy.visitDataset();
    });

    it("should be possible to add creator without first name", () => {
      cy.updateFields(
        "Creators",
        () => {
          cy.setFieldByLabel("Last name", "Dow");

          cy.contains("Dow External, non-UGent")
            .closest(".list-group-item")
            .contains(".btn", "Add external creator")
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
        .should("contain", "[missing] Dow");
    });

    it("should be possible to add creator without last name", () => {
      cy.updateFields(
        "Creators",
        () => {
          cy.setFieldByLabel("First name", "Jane");

          cy.contains("Jane External, non-UGent")
            .closest(".list-group-item")
            .contains(".btn", "Add external creator")
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
        .should("contain", "Jane [missing]");
    });
  });
});
