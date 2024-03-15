// https://github.com/ugent-library/biblio-backoffice/issues/1370

import * as dayjs from "dayjs";

import { getRandomText } from "support/util";

describe("Issue #1370: Make created, edited and system update timestamp more informative.", () => {
  const randomText = getRandomText();
  let creationTime: string;
  let editedTime: string;

  describe("for publications", () => {
    before(() => {
      cy.loginAsResearcher();
      cy.setUpPublication("Miscellaneous", {
        title: `Publication ${randomText}`,
      }).then(() => {
        creationTime = dayjs().format("DD/MM/YYYY HH:mm");
      });

      cy.loginAsLibrarian();
      cy.visitPublication();
      cy.updateFields(
        "Publication details",
        () => {
          cy.setFieldByLabel("Publication year", "2000");
        },
        true,
      ).then(() => {
        editedTime = dayjs().format("DD/MM/YYYY HH:mm");
      });
    });

    beforeEach(() => {
      cy.loginAsResearcher();

      cy.visit("/publication", { qs: { q: randomText } });
      cy.extractBiblioId();
    });

    it("should display the date created and edited in the publications list", () => {
      cy.visit("/publication", { qs: { q: randomText } });

      getCreatedText(".list-group-item .c-meta-item").should(
        "eq",
        `Created ${creationTime} by Biblio Researcher.`,
      );

      getEditedText(".list-group-item .c-meta-item").should(
        "eq",
        `Edited ${editedTime} by Biblio Librarian.`,
      );
    });

    it("should display the date created and edited in the publication detail page", () => {
      cy.visitPublication();

      getCreatedText("#summary .c-subline").should(
        "eq",
        `Created ${creationTime} by Biblio Researcher.`,
      );

      getEditedText("#summary .c-subline").should(
        "eq",
        `Edited ${editedTime} by Biblio Librarian.`,
      );
    });
  });

  describe("for datasets", () => {
    before(() => {
      cy.loginAsResearcher();
      cy.setUpDataset({ title: `Dataset ${randomText}` });

      cy.loginAsLibrarian();
      cy.visitDataset();
      cy.updateFields(
        "Dataset details",
        () => {
          cy.setFieldByLabel("Publication year", "2000");
        },
        true,
      );
    });

    beforeEach(() => {
      cy.loginAsResearcher();

      cy.visit("/dataset", { qs: { q: randomText } });
      cy.extractBiblioId();
    });

    it("should display the date created and edited in the datasets list", () => {
      cy.visit("/dataset", { qs: { q: randomText } });

      getCreatedText(".list-group-item .c-meta-item").should(
        "eq",
        `Created ${creationTime} by Biblio Researcher.`,
      );

      getEditedText(".list-group-item .c-meta-item").should(
        "eq",
        `Edited ${editedTime} by Biblio Librarian.`,
      );
    });

    it("should display the date created and edited in the dataset detail page", () => {
      cy.visitDataset();

      getCreatedText("#summary .c-subline").should(
        "eq",
        `Created ${creationTime} by Biblio Researcher.`,
      );

      getEditedText("#summary .c-subline").should(
        "eq",
        `Edited ${editedTime} by Biblio Librarian.`,
      );
    });
  });

  function getEditedText(selector: string) {
    return cy
      .contains(selector, "Edited")
      .should("be.visible")
      .invoke("text")
      .invoke("trim");
  }

  function getCreatedText(selector: string) {
    return cy
      .contains(selector, "Created")
      .should("be.visible")
      .invoke("text")
      .invoke("trim");
  }
});
