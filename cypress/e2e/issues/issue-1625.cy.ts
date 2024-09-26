// https://github.com/ugent-library/biblio-backoffice/issues/1625

describe("Issue #1625: Editors do not show up on WoS import", () => {
  it("should import editors from Web of Science", () => {
    cy.loginAsResearcher("researcher1");

    cy.visit("/add-publication", {
      qs: {
        method: "wos",
      },
    });
    cy.get("input[name=file]").selectFile(
      "cypress/fixtures/wos-000500696400004.txt",
    );

    cy.contains(".list-group-item", "People & Affiliations").click();

    cy.get("#contributors-author-body tbody .bc-avatar-text")
      .should("contain.text", "Hanne Dubois")
      .should("contain.text", "Geert van Loo")
      .should("contain.text", "Andy Wullaert");

    cy.get("#contributors-editor-body tbody .bc-avatar-text")
      .should("contain.text", "Claire Vanpouille-Box")
      .should("contain.text", "Lorenzo Galluzzi")
      .should("contain.text", "John Doe")
      .should("contain.text", "Jane Dow");
  });

  it("should use Web of Science field BE as fallback if BF is missing", () => {
    cy.loginAsResearcher("researcher1");

    cy.visit("/add-publication", {
      qs: {
        method: "wos",
      },
    });
    cy.get("input[name=file]").selectFile(
      "cypress/fixtures/wos-000500696400004-2.txt",
    );

    cy.contains(".list-group-item", "People & Affiliations").click();

    cy.get("#contributors-author-body tbody .bc-avatar-text")
      .should("contain.text", "Hanne Dubois")
      .should("contain.text", "Geert van Loo")
      .should("contain.text", "Andy Wullaert");

    cy.get("#contributors-editor-body tbody .bc-avatar-text")
      .should("contain.text", "C VanpouilleBox")
      .should("contain.text", "L Galluzzi")
      .should("contain.text", "John Doe")
      .should("contain.text", "Jane Dow");
  });
});
