// https://github.com/ugent-library/biblio-backoffice/issues/1370
// https://github.com/ugent-library/biblio-backoffice/issues/1258

import * as dayjs from "dayjs";
import * as CustomParseFormat from "dayjs/plugin/customParseFormat";

import { getRandomText } from "support/util";

dayjs.extend(CustomParseFormat);

describe("Record timestamps", () => {
  const RANDOM_TEXT = getRandomText();

  describe("for publications", () => {
    describe("normal flow", () => {
      before(() => {
        cy.loginAsResearcher("researcher1");
        cy.setUpPublication("Miscellaneous", {
          title: `Publication ${RANDOM_TEXT}`,
        });
        cy.addAuthor("Biblio", "Researcher2");

        cy.loginAsResearcher("researcher2");
        cy.addAuthor("New", "Author", { external: true });

        // Give elastic some extra time to index
        cy.wait(1000);
      });

      beforeEach(() => {
        cy.loginAsResearcher("researcher1");

        cy.visit("/publication", { qs: { q: RANDOM_TEXT } });
        cy.extractBiblioId();
      });

      it("should display the date created and edited in the publications list", () => {
        cy.visit("/publication", { qs: { q: RANDOM_TEXT } });

        assertTimestamp(
          ".list-group-item .c-meta-item",
          "Created",
          "Biblio Researcher1",
        );
        assertTimestamp(
          ".list-group-item .c-meta-item",
          "Edited",
          "Biblio Researcher2",
        );
      });

      it("should display the date created and edited in the publication detail page", () => {
        cy.visitPublication();

        assertTimestamp(
          "#summary .c-body-small",
          "Created",
          "Biblio Researcher1",
        );
        assertTimestamp(
          "#summary .c-body-small",
          "Edited",
          "Biblio Researcher2",
        );
      });
    });

    it("should not display curator names to non-curators", () => {
      cy.loginAsLibrarian("librarian1");

      cy.setUpPublication("Book");
      cy.addAuthor("Biblio", "Researcher1");
      cy.addAuthor("Biblio", "Librarian2");

      // Give elastic some extra time to index
      cy.wait(1000);

      cy.get<string>("@biblioId").then((biblioId) => {
        cy.visit("/publication", { qs: { q: biblioId } });

        assertTimestamp(
          ".list-group-item .c-meta-item",
          "Created",
          "Biblio Librarian1",
        );
        assertTimestamp(
          ".list-group-item .c-meta-item",
          "Edited",
          "Biblio Librarian1",
        );
      });

      cy.visitPublication();
      assertTimestamp("#summary .c-body-small", "Created", "Biblio Librarian1");
      assertTimestamp("#summary .c-body-small", "Edited", "Biblio Librarian1");

      // as a non-curator
      cy.loginAsResearcher("researcher1");

      cy.get<string>("@biblioId").then((biblioId) => {
        cy.visit("/publication", { qs: { q: biblioId } });

        assertTimestamp(
          ".list-group-item .c-meta-item",
          "Created",
          "a Biblio team member",
        );
        assertTimestamp(
          ".list-group-item .c-meta-item",
          "Edited",
          "a Biblio team member",
        );
      });

      cy.visitPublication();
      assertTimestamp(
        "#summary .c-body-small",
        "Created",
        "a Biblio team member",
      );
      assertTimestamp(
        "#summary .c-body-small",
        "Edited",
        "a Biblio team member",
      );

      // as a different curator
      cy.loginAsLibrarian("librarian2");

      cy.get<string>("@biblioId").then((biblioId) => {
        cy.visit("/publication", { qs: { q: biblioId } });

        assertTimestamp(
          ".list-group-item .c-meta-item",
          "Created",
          "Biblio Librarian1",
        );
        assertTimestamp(
          ".list-group-item .c-meta-item",
          "Edited",
          "Biblio Librarian1",
        );
      });

      cy.visitPublication();
      assertTimestamp("#summary .c-body-small", "Created", "Biblio Librarian1");
      assertTimestamp("#summary .c-body-small", "Edited", "Biblio Librarian1");
    });
  });

  describe("for datasets", () => {
    describe("normal flow", () => {
      before(() => {
        cy.loginAsResearcher("researcher1");
        cy.setUpDataset({
          title: `Dataset ${RANDOM_TEXT}`,
        });
        cy.addCreator("Biblio", "Researcher2");

        cy.loginAsResearcher("researcher2");
        cy.addCreator("New", "Creator", { external: true });

        // Give elastic some extra time to index
        cy.wait(1000);
      });

      beforeEach(() => {
        cy.loginAsResearcher("researcher1");

        cy.visit("/dataset", { qs: { q: RANDOM_TEXT } });
        cy.extractBiblioId();
      });

      it("should display the date created and edited in the datasets list", () => {
        cy.visit("/dataset", { qs: { q: RANDOM_TEXT } });

        assertTimestamp(
          ".list-group-item .c-meta-item",
          "Created",
          "Biblio Researcher1",
        );
        assertTimestamp(
          ".list-group-item .c-meta-item",
          "Edited",
          "Biblio Researcher2",
        );
      });

      it("should display the date created and edited in the dataset detail page", () => {
        cy.visitDataset();

        assertTimestamp(
          "#summary .c-body-small",
          "Created",
          "Biblio Researcher1",
        );
        assertTimestamp(
          "#summary .c-body-small",
          "Edited",
          "Biblio Researcher2",
        );
      });
    });

    it("should not display curator names to non-curators", () => {
      cy.loginAsLibrarian("librarian1");

      cy.setUpDataset();
      cy.addCreator("Biblio", "Researcher1");
      cy.addCreator("Biblio", "Librarian2");

      // Give elastic some extra time to index
      cy.wait(1000);

      cy.get<string>("@biblioId").then((biblioId) => {
        cy.visit("/dataset", { qs: { q: biblioId } });

        assertTimestamp(
          ".list-group-item .c-meta-item",
          "Created",
          "Biblio Librarian1",
        );
        assertTimestamp(
          ".list-group-item .c-meta-item",
          "Edited",
          "Biblio Librarian1",
        );
      });

      cy.visitDataset();
      assertTimestamp("#summary .c-body-small", "Created", "Biblio Librarian1");
      assertTimestamp("#summary .c-body-small", "Edited", "Biblio Librarian1");

      // as a non-curator
      cy.loginAsResearcher("researcher1");

      cy.get<string>("@biblioId").then((biblioId) => {
        cy.visit("/dataset", { qs: { q: biblioId } });

        assertTimestamp(
          ".list-group-item .c-meta-item",
          "Created",
          "a Biblio team member",
        );
        assertTimestamp(
          ".list-group-item .c-meta-item",
          "Edited",
          "a Biblio team member",
        );
      });

      cy.visitDataset();
      assertTimestamp(
        "#summary .c-body-small",
        "Created",
        "a Biblio team member",
      );
      assertTimestamp(
        "#summary .c-body-small",
        "Edited",
        "a Biblio team member",
      );

      // as a different curator
      cy.loginAsLibrarian("librarian2");

      cy.get<string>("@biblioId").then((biblioId) => {
        cy.visit("/dataset", { qs: { q: biblioId } });

        assertTimestamp(
          ".list-group-item .c-meta-item",
          "Created",
          "Biblio Librarian1",
        );
        assertTimestamp(
          ".list-group-item .c-meta-item",
          "Edited",
          "Biblio Librarian1",
        );
      });

      cy.visitDataset();
      assertTimestamp("#summary .c-body-small", "Created", "Biblio Librarian1");
      assertTimestamp("#summary .c-body-small", "Edited", "Biblio Librarian1");
    });
  });

  function assertTimestamp(
    selector: string,
    type: "Created" | "Edited",
    byUser: "a Biblio team member" | string,
  ) {
    const regex = new RegExp(
      `^${type} (?<timestamp>\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}) by ${byUser}.$`,
    );

    return cy
      .get(`${selector}:contains(${type})`)
      .should("have.length", 1)
      .should("be.visible")
      .invoke("text")
      .invoke("trim")
      .should("match", regex)
      .then((text) => text.match(regex).groups.timestamp)
      .should("be.oneOf", [
        dayjs().subtract(1, "minute").format("YYYY-MM-DD HH:mm"),
        dayjs().format("YYYY-MM-DD HH:mm"),
        dayjs().add(1, "minute").format("YYYY-MM-DD HH:mm"),
      ]);
  }
});
