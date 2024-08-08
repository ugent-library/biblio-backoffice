// https://github.com/ugent-library/biblio-backoffice/issues/1680

describe('Issue #1680: Make sure tabs are not reset when clicking "reset filters"', () => {
  beforeEach(() => {
    cy.loginAsResearcher("researcher1");
  });

  describe("for publications", () => {
    const tabs = [
      { label: "My publications", scope: "contributed" },
      { label: "Supervised by me", scope: "supervised" },
      { label: "Registered by me", scope: "created" },
    ];

    tabs.forEach((tab) => {
      it(`should not reset the tab "${tab.label}" when clicking "reset filters"`, () => {
        cy.visit("/publication");

        performTest(tab);
      });
    });
  });

  describe("for datasets", () => {
    const tabs = [
      { label: "My datasets", scope: "contributed" },
      { label: "Registered by me", scope: "created" },
    ];

    tabs.forEach((tab) => {
      it(`should not reset the tab "${tab.label}" when clicking "reset filters"`, () => {
        cy.visit("/dataset");

        performTest(tab);
      });
    });
  });

  function performTest(tab: { label: string; scope: string }) {
    cy.contains(".nav-link", tab.label).click();

    cy.get("input[name=q]:not([type=hidden])").type("my filter query");
    cy.contains(".btn", "Search").click();

    cy.contains(".dropdown", "Biblio status").within(() => {
      cy.contains("a.badge", "Biblio status").click();
      cy.contains("label", "Public").click();
      cy.contains(".btn", "Apply filter").click();
    });

    verifyFilter(tab, true);

    cy.contains(".btn", "Reset filters").click();

    verifyFilter(tab, false);
  }

  function verifyFilter(
    tab: { label: string; scope: string },
    beforeReset: boolean,
  ) {
    cy.getParams("f[scope]").should("eq", tab.scope);
    cy.contains(".nav-link", tab.label).should("have.class", "active");

    cy.getParams("q").should("eq", "my filter query");
    cy.get("input[name=q]:not([type=hidden])").should(
      "have.value",
      "my filter query",
    );

    if (beforeReset) {
      cy.getParams("f[status]").should("eq", "public");
      cy.contains(".dropdown a.badge", "Biblio status")
        .should("have.class", "bg-primary")
        .and("contain.text", "Public");
    } else {
      cy.getParams("f[status]").should("not.exist");
      cy.contains(".dropdown a.badge", "Biblio status")
        .should("not.have.class", "bg-primary")
        .and("not.contain.text", "Public");
    }
  }
});
