// https://github.com/ugent-library/biblio-backoffice/issues/1370

import * as dayjs from "dayjs";
import * as CustomParseFormat from "dayjs/plugin/customParseFormat";

import { getRandomText } from "support/util";

dayjs.extend(CustomParseFormat);

describe("Issue #1370: Make created, edited and system update timestamp more informative.", () => {
  const RANDOM_TEXT = getRandomText();
  const CREATED_REGEX =
    /^Created (?<timestamp>\d{4}-\d{2}-\d{2} \d{2}:\d{2}) by Biblio Researcher1.$/;
  const EDITED_REGEX =
    /^Edited (?<timestamp>\d{4}-\d{2}-\d{2} \d{2}:\d{2}) by Biblio Librarian1.$/;

  describe("for publications", () => {
    before(() => {
      cy.login("researcher1");
      cy.setUpPublication("Miscellaneous", {
        title: `Publication ${RANDOM_TEXT}`,
      });

      cy.login("librarian1");
      cy.visitPublication();
      cy.updateFields(
        "Publication details",
        () => {
          cy.setFieldByLabel("Publication year", "2000");
        },
        true,
      );

      // Give elastic some extra time to index
      cy.wait(1000);
    });

    beforeEach(() => {
      cy.login("researcher1");

      cy.visit("/publication", { qs: { q: RANDOM_TEXT } });
      cy.extractBiblioId();
    });

    it("should display the date created and edited in the publications list", () => {
      cy.visit("/publication", { qs: { q: RANDOM_TEXT } });

      assertTimestamp(
        ".list-group-item .c-meta-item",
        "Created",
        CREATED_REGEX,
      );
      assertTimestamp(".list-group-item .c-meta-item", "Edited", EDITED_REGEX);
    });

    it("should display the date created and edited in the publication detail page", () => {
      cy.visitPublication();

      assertTimestamp("#summary .c-subline", "Created", CREATED_REGEX);
      assertTimestamp("#summary .c-subline", "Edited", EDITED_REGEX);
    });
  });

  describe("for datasets", () => {
    before(() => {
      cy.login("researcher1");
      cy.setUpDataset({
        title: `Dataset ${RANDOM_TEXT}`,
        shouldWaitForIndex: true,
      });

      cy.login("librarian1");
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
      cy.login("researcher1");

      cy.visit("/dataset", { qs: { q: RANDOM_TEXT } });
      cy.extractBiblioId();
    });

    it("should display the date created and edited in the datasets list", () => {
      cy.visit("/dataset", { qs: { q: RANDOM_TEXT } });

      assertTimestamp(
        ".list-group-item .c-meta-item",
        "Created",
        CREATED_REGEX,
      );
      assertTimestamp(".list-group-item .c-meta-item", "Edited", EDITED_REGEX);
    });

    it("should display the date created and edited in the dataset detail page", () => {
      cy.visitDataset();

      assertTimestamp("#summary .c-subline", "Created", CREATED_REGEX);
      assertTimestamp("#summary .c-subline", "Edited", EDITED_REGEX);
    });
  });

  function assertTimestamp(
    selector: string,
    substring: "Created" | "Edited",
    REGEX: RegExp,
  ) {
    return cy
      .get(`${selector}:contains(${substring})`)
      .should("have.length", 1)
      .should("be.visible")
      .invoke("text")
      .invoke("trim")
      .should("match", REGEX)
      .then((text) => {
        const { timestamp } = text.match(REGEX).groups;
        const created = dayjs(timestamp, "YYYY-MM-DD HH:mm");

        // Allow a 2 minute margin of error
        const lower = dayjs().second(0).millisecond(0).subtract(1, "minute");
        const upper = lower.clone().add(2, "minutes");

        expect(created.unix()).to.be.within(lower.unix(), upper.unix());
      });
  }
});
