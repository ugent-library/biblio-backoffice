// https://github.com/ugent-library/biblio-backoffice/issues/0844

describe('Issue #0844: [filters] Apply filters when clicking on "apply" and when clicking next to dialogue box', () => {
  const DUMMY_WINDOW_PROPERTY = "DUMMY_WINDOW_PROPERTY";

  const ACTION_CLICK_APPLY_BUTTON = () =>
    cy.contains(".btn", "Apply filter").click();
  const ACTION_CLICK_OUTSIDE_FACET = () =>
    cy.document().find("body").contains("h4", "Overview").click();
  const ACTION_HIT_ESCAPE_KEY = () => cy.document().find("body").type("{esc}");

  before(() => {
    cy.loginAsLibrarian();

    cy.setUpPublication();
    cy.visitPublication();
    cy.updateFields(
      "Publication details",
      () => {
        cy.setFieldByLabel(
          "Publication year",
          new Date().getFullYear().toString(),
        );
        cy.setFieldByLabel("Web of Science type", "dummy WoS type");
      },
      true,
    );
    cy.updateFields(
      "Librarian tags",
      () => {
        cy.setFieldByLabel("Librarian tags", "dummy librarian tag");
      },
      true,
    );

    cy.setUpDataset();
    cy.visitDataset();
    cy.updateFields(
      "Librarian tags",
      () => {
        cy.getLabel("Librarian tags")
          .next()
          .find("tags")
          .type("dummy librarian tag{enter}");
      },
      true,
    );
  });

  describe("as researcher", () => {
    beforeEach(() => {
      cy.loginAsResearcher();
    });

    describe("for publications", () => {
      const FACETS = [
        "status",
        "classification",
        "faculty_id",
        "year",
        "type",
        "has_message",
        "locked",
        "has_files",
        "vabb_type",
        "created_since",
        "updated_since",
      ];

      generateTests("/publication", FACETS);
    });

    describe("for datasets", () => {
      const FACETS = [
        "status",
        "faculty_id",
        "locked",
        "has_message",
        "created_since",
        "updated_since",
      ];

      generateTests("/dataset", FACETS);
    });
  });

  describe("as librarian", () => {
    beforeEach(() => {
      cy.loginAsLibrarian();

      cy.switchMode("Librarian");
    });

    describe("for publications", () => {
      const FACETS = [
        "status",
        "classification",
        "faculty_id",
        "year",
        "type",
        "publication_status",
        "reviewer_tags",
        "has_message",
        "locked",
        "extern",
        "wos_type",
        "vabb_type",
        "has_files",
        "file_relation",
        "created_since",
        "updated_since",
        "legacy",
      ];

      generateTests("/publication", FACETS);
    });

    describe("for datasets", () => {
      const FACETS = [
        "status",
        "faculty_id",
        "locked",
        "identifier_type",
        "reviewer_tags",
        "has_message",
        "created_since",
        "updated_since",
      ];

      generateTests("/dataset", FACETS);
    });
  });

  function generateTests(route: string, facets: string[]) {
    beforeEach(() => {
      cy.visit(route);

      cy.then(() => {
        // Display all filters if necessary
        // Using jQuery here as the Cypress way would fail if the label is missing
        Cypress.$('label:contains("Show all filters")').trigger("click");
      });

      cy.get("[data-facet-dropdown]").should("have.length", facets.length);

      cy.window().then((w) => {
        return (w[DUMMY_WINDOW_PROPERTY] = DUMMY_WINDOW_PROPERTY);
      });
    });

    // Generate the tests for 3 random facets
    Cypress._.chain(facets)
      .shuffle()
      .take(3)
      .value()
      .forEach((facet) => {
        describe(`for facet "${facet}"`, () => {
          beforeEach(() => {
            // Verify there was no filter query param to start with
            cy.getParams(`f[${facet}]`).should("not.exist");

            // Open the facet filter
            cy.get(`[data-facet-dropdown="${facet}"]`)
              .as("dropdown")
              .find("a.badge")
              .should("not.have.class", "bg-primary");
          });

          it('should apply filters when clicking on "apply"', () => {
            executeTest(facet, ACTION_CLICK_APPLY_BUTTON);
          });

          it("should apply filters when clicking on next to dialogue box", () => {
            executeTest(facet, ACTION_CLICK_OUTSIDE_FACET);
          });

          it("should apply filters when hitting the ESCAPE key", () => {
            executeTest(facet, ACTION_HIT_ESCAPE_KEY);
          });

          it("should not apply filters when selection has not changed", () => {
            cy.get("@dropdown")
              .click()
              .find(".dropdown-menu")
              .within(() => {
                cy.contains(".btn", "Apply filter").should("be.visible");

                // Set the filter field
                fillOutFilterFields(facet);

                // Unset the filter field
                cy.get("@facetField").then(($input) => {
                  if ($input.prop("type") === "text") {
                    cy.wrap($input).clear();
                  } else if ($input.prop("type") === "checkbox") {
                    cy.wrap($input).uncheck();
                  }

                  cy.window().should("have.property", DUMMY_WINDOW_PROPERTY);

                  ACTION_CLICK_OUTSIDE_FACET();

                  cy.contains(".btn", "Apply filter").should("not.be.visible");
                });
              });

            // Verify that page was NOT reloaded
            cy.getParams(`f[${facet}]`).should("not.exist");
            cy.window().should("have.property", DUMMY_WINDOW_PROPERTY);
            cy.get("@dropdown")
              .find("a.badge")
              .should("not.have.class", "bg-primary");
          });
        });
      });
  }

  function executeTest(facet: string, triggerFacetFilterCallback: () => void) {
    cy.get("@dropdown").click();

    cy.get("@dropdown")
      .find(".dropdown-menu")
      .within(() => {
        cy.contains(".btn", "Apply filter").should("be.visible");

        fillOutFilterFields(facet);

        cy.window().should("have.property", DUMMY_WINDOW_PROPERTY);

        triggerFacetFilterCallback();

        cy.contains(".btn", "Apply filter").should("not.be.visible");
      });

    // Verify that page was reloaded
    cy.get("@facetValue").then((facetValue) => {
      cy.getParams(`f[${facet}]`).should("eq", facetValue);
    });
    cy.window().should("not.have.property", DUMMY_WINDOW_PROPERTY);
    cy.get("@dropdown").find("a.badge").should("have.class", "bg-primary");
  }

  function fillOutFilterFields(facet: string) {
    cy.get<HTMLInputElement>(`input[name="f[${facet}]"]`)
      // In case of checkboxes, we pick a random one to check
      .random()
      .as("facetField")
      .then(($input) => {
        if ($input.prop("type") === "text") {
          cy.wrap($input).type(getRandomTextFieldValue());
        } else if ($input.prop("type") === "checkbox") {
          cy.wrap($input).check();
        } else {
          throw new Error(`Unsupported input type "${$input.prop("type")}".`);
        }
      });

    cy.get("@facetField").prop("value").as("facetValue", { type: "static" });
  }

  function getRandomTextFieldValue() {
    return Cypress._.chain([
      "2024-03-02",
      "2023",
      "yesterday",
      "today",
      "tomorrow",
    ])
      .shuffle()
      .first()
      .value();
  }
});
