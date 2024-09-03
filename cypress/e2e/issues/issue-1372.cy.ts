// https://github.com/ugent-library/biblio-backoffice/issues/1372

import { getRandomText } from "support/util";

describe("Issue #1372: Administrative support staff cannot see linked datasets/publications in a record if they are not assigned", () => {
  const publicationTitle = getRandomText();
  const datasetTitle = getRandomText();

  before(() => {
    cy.loginAsLibrarian("librarian1");

    cy.setUpPublication(undefined, {
      title: publicationTitle,
      biblioIDAlias: "publication",
      publish: true,
      shouldWaitForIndex: true,
    });
    cy.addAuthor("Extra", "Author1", {
      biblioIDAlias: "@publication",
      external: true,
    });
    cy.addAuthor("Extra", "Author2", {
      biblioIDAlias: "@publication",
      external: true,
    });
    cy.addAuthor("Extra", "Author3", {
      biblioIDAlias: "@publication",
      external: true,
    });
    cy.addAuthor("Extra", "Author4", {
      biblioIDAlias: "@publication",
      external: true,
    });
    cy.addAuthor("Biblio", "Researcher1", { biblioIDAlias: "@publication" });

    cy.setUpDataset({
      title: datasetTitle,
      biblioIDAlias: "dataset",
      publish: true,
      shouldWaitForIndex: true,
    });
    cy.addCreator("Extra", "Creator1", {
      biblioIDAlias: "@dataset",
      external: true,
    });
    cy.addCreator("Extra", "Creator2", {
      biblioIDAlias: "@dataset",
      external: true,
    });
    cy.addCreator("Extra", "Creator3", {
      biblioIDAlias: "@dataset",
      external: true,
    });
    cy.addCreator("Extra", "Creator4", {
      biblioIDAlias: "@dataset",
      external: true,
    });
    cy.addCreator("Biblio", "Researcher2", { biblioIDAlias: "@dataset" });

    // Link the two together
    cy.visitPublication("@publication");
    cy.contains(".nav .nav-item", "Datasets").click();

    cy.get("#datasets-body").should("contain.text", "No datasets");

    cy.contains(".card", "Related datasets")
      .contains(".btn", "Add dataset")
      .click();

    cy.ensureModal("Select datasets").within(() => {
      cy.intercept("/publication/*/datasets/suggestions?*").as(
        "suggestDataset",
      );

      cy.getLabel("Search datasets")
        .next("input")
        .should("be.focused")
        .type(datasetTitle);
      cy.wait("@suggestDataset");

      cy.contains(".list-group-item", datasetTitle)
        .contains(".btn", "Add dataset")
        .click();
    });
    cy.ensureNoModal();
  });

  beforeEach(() => {
    cy.loginAsLibrarian("librarian1", "Researcher");

    // Refresh biblio ID aliases set in before-hook for subsequent tests
    cy.visit("/publication", { qs: { q: publicationTitle } });
    cy.extractBiblioId("publication");

    cy.visit("/dataset", { qs: { q: datasetTitle } });
    cy.extractBiblioId("dataset");
  });

  it("should be possible for researchers to see linked datasets of a publication, even if they cannot access them directly", () => {
    cy.visitPublication("@publication");
    cy.contains(".nav .nav-item", "Datasets").click();

    cy.get("#datasets-body").should("not.contain.text", "No datasets");
    cy.get("#datasets-body .list-group-item")
      .as("item")
      .should("have.length", 1)
      .should("contain.text", datasetTitle);

    verifyActionability("publication", true, true);

    // Now view this as researcher
    cy.loginAsResearcher("researcher1");
    cy.get<string>("@dataset").then((datasetId) => {
      cy.request({
        url: `/dataset/${datasetId}`,
        failOnStatusCode: false,
      }).then((response) => {
        expect(response.status).to.eq(403);
        expect(response.statusText).to.eq("Forbidden");
      });
    });

    cy.visitPublication("@publication");
    cy.contains(".nav .nav-item", "Datasets").click();

    cy.contains(".btn", "Add dataset").should("be.visible");
    cy.get("#datasets-body").should("not.contain.text", "No datasets");
    cy.get("#datasets-body .list-group-item")
      .should("have.length", 1)
      .should("contain.text", datasetTitle);

    verifyActionability("publication", false, true);

    // Now lock the publication and make sure the researcher cannot add and remove datasets anymore
    cy.loginAsLibrarian("librarian1");
    cy.visitPublication("@publication");
    cy.contains(".btn", "Lock record").click();
    cy.closeToast();

    cy.loginAsResearcher("researcher1");
    cy.visitPublication("@publication");
    cy.contains(".nav .nav-item", "Datasets").click();
    cy.contains(".btn", "Add dataset").should("not.exist");

    verifyActionability("publication", false, false);
  });

  it("should be possible for researchers to see linked publications of a dataset, even if they cannot access them directly", () => {
    cy.visitDataset("@dataset");
    cy.contains(".nav .nav-item", "Publications").click();

    cy.get("#publications-body").should("not.contain.text", "No publications");
    cy.get("#publications-body .list-group-item")
      .as("item")
      .should("have.length", 1)
      .should("contain.text", publicationTitle);

    verifyActionability("dataset", true, true);

    // Now view this as researcher
    cy.loginAsResearcher("researcher2");
    cy.get<string>("@publication").then((publicationId) => {
      cy.request({
        url: `/publication/${publicationId}`,
        failOnStatusCode: false,
      }).then((response) => {
        expect(response.status).to.eq(403);
        expect(response.statusText).to.eq("Forbidden");
      });
    });

    cy.visitDataset("@dataset");
    cy.contains(".nav .nav-item", "Publications").click();

    cy.contains(".btn", "Add publication").should("be.visible");
    cy.get("#publications-body").should("not.contain.text", "No publications");
    cy.get("#publications-body .list-group-item")
      .should("have.length", 1)
      .should("contain.text", publicationTitle);

    verifyActionability("dataset", false, true);

    // Now lock the dataset and make sure the researcher cannot add and remove publications anymore
    cy.loginAsLibrarian("librarian1");
    cy.visitDataset("@dataset");
    cy.contains(".btn", "Lock record").click();
    cy.closeToast();

    cy.loginAsResearcher("researcher2");
    cy.visitDataset("@dataset");
    cy.contains(".nav .nav-item", "Publications").click();
    cy.contains(".btn", "Add publication").should("not.exist");

    verifyActionability("dataset", false, false);
  });

  function verifyActionability(
    context: "publication" | "dataset",
    isActionable: boolean,
    canAddAndRemoveRelatedItems: boolean,
  ) {
    const prefix = isActionable ? "" : "not.";
    const relatedContext =
      context === "publication" ? "dataset" : "publication";
    const relatedContributorType =
      relatedContext === "publication" ? "author" : "creator";

    cy.get("@item").within(() => {
      // Verify the title is (not) a link to the detail page
      const $title = cy
        .get(".list-group-item-title")
        .parent()
        .should(prefix + "have.prop", "tagName", "A");

      if (isActionable) {
        $title
          .should("have.attr", "href")
          .should("match", new RegExp(`^/${relatedContext}/[A-Z0-9]+$`));
      }

      // Verify there is an/no "more authors/creators" link
      cy.get(".c-author-list .c-author:not(:contains('Your role'))")
        .as("contributors")
        .map<HTMLElement, string>((i) => i.textContent.trim())
        .should("have.ordered.members", [
          "John Doe",
          `Extra ${Cypress._.capitalize(relatedContributorType)}1`,
          `Extra ${Cypress._.capitalize(relatedContributorType)}2`,
          `3 more ${relatedContributorType}s`,
        ]);

      const $moreContributors = cy
        .get("@contributors")
        .contains(`3 more ${relatedContributorType}s`)
        .should(prefix + "have.prop", "tagName", "A");

      if (isActionable) {
        $moreContributors
          .should(prefix + "have.attr", "href")
          .should(
            "match",
            new RegExp(`^/${relatedContext}/[A-Z0-9]+\\?show=contributors$`),
          );
      }

      // Verify there is an/no "Add department" link
      const $addDepartment = cy
        .contains("a", "Add department")
        .should(isActionable ? "be.visible" : "not.exist");

      if (isActionable) {
        $addDepartment
          .should("have.attr", "href")
          .should(
            "match",
            new RegExp(`^/${relatedContext}/[A-Z0-9]+\\?show=contributors$`),
          );
      }

      // Verify button toolbar
      cy.get(".c-button-toolbar").should(
        isActionable || canAddAndRemoveRelatedItems
          ? "be.visible"
          : "not.exist",
      );

      if (isActionable || canAddAndRemoveRelatedItems) {
        // Verify there is an/no link to "View publication/dataset" in the button toolbar
        cy.contains(".c-button-toolbar .btn", `View`).should(
          isActionable ? "be.visible" : "not.exist",
        );

        // In any case, there should be a "Remove link" button in the button toolbar
        cy.contains(".c-button-toolbar .btn", "Remove link").should(
          "be.visible",
        );
      }
    });
  }
});
