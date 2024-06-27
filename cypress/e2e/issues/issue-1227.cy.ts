import { getRandomText } from "support/util";

describe("Issue #1127: Cannot search any longer on book title, journal title, short journal title nor conference title", () => {
  beforeEach(() => {
    cy.loginAsResearcher();
  });

  it("should be possible to search by publisher", () => {
    const randomText = getRandomText();

    cy.setUpPublication("Miscellaneous", {
      otherFields: {
        publisher: `Publisher: ${randomText}`,
      },
      shouldWaitForIndex: true,
    });

    cy.wait(1000);

    cy.visit("/publication");

    cy.search(randomText).should("eq", 1);
  });

  it("should be possible to search by alternative title", () => {
    const randomText1 = getRandomText();
    const randomText2 = getRandomText();

    cy.setUpPublication("Miscellaneous", {
      otherFields: {
        alternative_title: [
          `Alternative title: ${randomText1}`,
          `Alternative title: ${randomText2}`,
        ],
      },
      shouldWaitForIndex: true,
    });

    cy.wait(1000);

    cy.visit("/publication");

    cy.search(randomText1).should("eq", 1);
    cy.search(randomText2).should("eq", 1);
  });

  it("should be possible to search by conference title", () => {
    const randomText = getRandomText();

    cy.setUpPublication("Conference contribution", {
      shouldWaitForIndex: true,
    });

    cy.visitPublication();

    cy.updateFields(
      "Conference details",
      () => {
        cy.setFieldByLabel("Conference", `The conference name: ${randomText}`);
      },
      true,
    );

    cy.wait(1000);

    cy.visit("/publication");

    cy.search(randomText).should("eq", 1);
  });

  it("should be possible to search by journal title", () => {
    const randomText = getRandomText();

    cy.setUpPublication("Journal Article", {
      otherFields: {
        publication: `The journal title: ${randomText}`,
      },
      shouldWaitForIndex: true,
    });

    cy.wait(1000);

    cy.visit("/publication");

    cy.search(randomText).should("eq", 1);
  });

  it("should be possible to search by short journal title", () => {
    const randomText = getRandomText();

    cy.setUpPublication("Journal Article", {
      otherFields: {
        publication_abbreviation: `The short journal title: ${randomText}`,
      },
      shouldWaitForIndex: true,
    });

    cy.wait(1000);

    cy.visit("/publication");

    cy.search(randomText).should("eq", 1);
  });

  it("should be possible to search by book title", () => {
    const randomText = getRandomText();

    cy.setUpPublication("Book Chapter", {
      otherFields: {
        publication: `The book title: ${randomText}`,
      },
      shouldWaitForIndex: true,
    });

    cy.wait(1000);

    cy.visit("/publication");

    cy.search(randomText).should("eq", 1);
  });
});
