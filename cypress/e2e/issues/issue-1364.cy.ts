// https://github.com/ugent-library/biblio-backoffice/issues/1364

describe('Issue #1364: Add "Updated (oldest first)" to sorting options', () => {
  const randomTitle = crypto.randomUUID();

  before(() => {
    cy.loginAsResearcher();

    for (let i = 1; i <= 3; i++) {
      cy.setUpPublication("Miscellaneous", false, `Title ${i} ${randomTitle}`);
      cy.extractBiblioId().as(`publication_${i}`);
      cy.visitPublication(`@publication_${i}`);
      cy.updateFields(
        "Publication details",
        () => {
          cy.setFieldByLabel(
            "Publication year",
            (new Date().getFullYear() - 10 + i).toString(),
          );
        },
        true,
      );

      cy.setUpDataset(false, `Title ${i} ${randomTitle}`);
      cy.extractBiblioId().as(`dataset_${i}`);
      cy.visitDataset(`@dataset_${i}`);
      cy.updateFields(
        "Dataset details",
        () => {
          cy.setFieldByLabel(
            "Publication year",
            (new Date().getFullYear() - 10 + i).toString(),
          );
        },
        true,
      );
    }

    // Update 2nd publication/dataset again to give it a newer updated date
    cy.visitPublication("@publication_2");
    cy.updateFields(
      "Publication details",
      () => {
        cy.setFieldByLabel("Publisher", "An updated publisher title");
      },
      true,
    );

    cy.visitDataset("@dataset_2");
    cy.updateFields(
      "Dataset details",
      () => {
        cy.setFieldByLabel("Publisher", "An updated publisher title");
      },
      true,
    );
  });

  beforeEach(() => cy.loginAsResearcher());

  describe("for publications", () => {
    it("should be possible to sort by oldest updated first", () => {
      cy.visit("/publication", { qs: { q: randomTitle } });

      executeTest();
    });
  });

  describe("for datasets", () => {
    it("should be possible to sort by oldest updated first", () => {
      cy.visit("/dataset", { qs: { q: randomTitle } });

      executeTest();
    });
  });

  function executeTest() {
    cy.get("select[name=sort] option")
      .map(({ value, textContent }: HTMLOptionElement) => ({
        [value]: textContent,
      }))
      .should("have.length", 5)
      .should("eql", [
        { "date-updated-desc": "Updated (newest first)" },
        { "date-updated-asc": "Updated (oldest first)" },
        { "date-created-desc": "Added (newest first)" },
        { "date-created-asc": "Added (oldest first)" },
        { "year-desc": "Publication year (newest first)" },
      ]);

    getTitles().should("eql", ["Title 2", "Title 3", "Title 1"]);

    cy.setFieldByLabel("Sort by", "date-updated-asc");
    cy.wait(100);
    getTitles().should("eql", ["Title 1", "Title 3", "Title 2"]);

    cy.setFieldByLabel("Sort by", "date-created-desc");
    cy.wait(100);
    getTitles().should("eql", ["Title 3", "Title 2", "Title 1"]);

    cy.setFieldByLabel("Sort by", "date-created-asc");
    cy.wait(100);
    getTitles().should("eql", ["Title 1", "Title 2", "Title 3"]);

    cy.setFieldByLabel("Sort by", "year-desc");
    cy.wait(100);
    getTitles().should("eql", ["Title 3", "Title 2", "Title 1"]);
  }

  function getTitles() {
    return cy
      .get(".list-group-item h4")
      .map<HTMLElement, string>((t) =>
        t.textContent.replace(new RegExp(randomTitle + ".*", "g"), "").trim(),
      );
  }
});
