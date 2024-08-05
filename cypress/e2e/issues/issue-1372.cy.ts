// https://github.com/ugent-library/biblio-backoffice/issues/1372

import { getRandomText } from "support/util";

describe("Issue #1372: Administrative support staff cannot see linked datasets/publications in a record if they are not assigned", () => {
  const publicationTitle = getRandomText();
  const datasetTitle = getRandomText();

  before(() => {
    cy.login("librarian1");

    cy.setUpPublication(undefined, {
      title: publicationTitle,
      biblioIDAlias: "publication",
      publish: true,
      shouldWaitForIndex: true,
    });
    cy.addAuthor("Biblio", "Researcher1", { biblioIDAlias: "@publication" });

    cy.setUpDataset({
      title: datasetTitle,
      biblioIDAlias: "dataset",
      publish: true,
      shouldWaitForIndex: true,
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

      cy.getLabel("Search")
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

  it("should be possible for researchers to see linked datasets of a publication, even if they cannot access them directly", () => {
    cy.login("librarian1");

    cy.visit("/publication", { qs: { q: publicationTitle } });
    cy.extractBiblioId("publication");

    cy.visitPublication("@publication");
    cy.contains(".nav .nav-item", "Datasets").click();

    cy.get("#datasets-body")
      .should("not.contain.text", "No datasets")
      .find(".list-group-item")
      .should("have.length", 1)
      .should("contain.text", datasetTitle);

    // Now view this as researcher
    cy.login("researcher1");
    cy.visitPublication("@publication");
    cy.contains(".nav .nav-item", "Datasets").click();

    cy.get("#datasets-body")
      .should("not.contain.text", "No datasets")
      .find(".list-group-item")
      .should("have.length", 1)
      .should("contain.text", datasetTitle);
  });

  it("should be possible for researchers to see linked publications of a dataset, even if they cannot access them directly", () => {
    cy.login("librarian1");

    cy.visit("/dataset", { qs: { q: datasetTitle } });
    cy.extractBiblioId("dataset");

    cy.visitDataset("@dataset");
    cy.contains(".nav .nav-item", "Publications").click();

    cy.get("#publications-body")
      .should("not.contain.text", "No publications")
      .find(".list-group-item")
      .should("have.length", 1)
      .should("contain.text", publicationTitle);

    // Now view this as researcher
    cy.login("researcher2");
    cy.visitDataset("@dataset");
    cy.contains(".nav .nav-item", "Publications").click();

    cy.get("#publications-body")
      .should("not.contain.text", "No publications")
      .find(".list-group-item")
      .should("have.length", 1)
      .should("contain.text", publicationTitle);
  });
});
