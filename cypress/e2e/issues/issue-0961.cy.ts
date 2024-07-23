// https://github.com/ugent-library/biblio-backoffice/issues/961

describe("Issue #961: [filters] Prioritise filter sequence and visibility", () => {
  describe("for researchers", () => {
    beforeEach(() => {
      cy.login("librarian1");
    });

    it("should not show collapsible facet filters for publications", () => {
      cy.visit("/publication");

      assertNoCollapsibleFacetFilters();

      listAllFacetsByFacetLine().should("eql", [
        [
          "Biblio status",
          "Classification",
          "Faculty",
          "Publication year",
          "Publication type",
        ],
        [
          "Message",
          "Locked",
          "File",
          "VABB type",
          "Created since",
          "Updated since",
        ],
      ]);
    });

    it("should not show collapsible facet filters for datasets", () => {
      cy.visit("/dataset");

      assertNoCollapsibleFacetFilters();

      listAllFacetsByFacetLine().should("eql", [
        [
          "Biblio status",
          "Faculty",
          "Locked",
          "Message",
          "Created since",
          "Updated since",
        ],
      ]);
    });
  });

  describe("for librarians", () => {
    beforeEach(() => {
      cy.login("librarian1");
    });

    describe("in researcher mode", () => {
      beforeEach(() => {
        cy.visit("/");

        cy.get(".c-sidebar .dropdown button.dropdown-toggle .visually-hidden")
          .should("have.length", 1)
          .and("have.text", "Researcher");
      });

      it("should not show collapsible facet filters for publications", () => {
        cy.visit("/publication");

        assertNoCollapsibleFacetFilters();

        listAllFacetsByFacetLine().should("eql", [
          [
            "Biblio status",
            "Classification",
            "Faculty",
            "Publication year",
            "Publication type",
          ],
          [
            "Message",
            "Locked",
            "File",
            "VABB type",
            "Created since",
            "Updated since",
          ],
        ]);
      });

      it("should not show collapsible facet filters for datasets", () => {
        cy.visit("/dataset");

        assertNoCollapsibleFacetFilters();

        listAllFacetsByFacetLine().should("eql", [
          [
            "Biblio status",
            "Faculty",
            "Locked",
            "Message",
            "Created since",
            "Updated since",
          ],
        ]);
      });
    });

    describe("in librarian mode", () => {
      beforeEach(() => {
        cy.switchMode("Librarian");
      });

      it("should show collapsible facet filters for publications", () => {
        cy.visit("/publication");

        cy.get(".toggle-zone").should("be.visible");
        cy.get("#show-all-facet-filters-toggle").should("be.visible");

        cy.contains(".bc-toolbar", "Reset filters")
          .find(".bc-toolbar-left .bc-toolbar-item .badge-list")
          .find(".dropdown .badge")
          .as("facets");

        assertAssetsAreCollapsed();
        cy.get("@facets").filter(":visible").should("have.length", 10);

        cy.contains(".btn", "Show all filters").click();
        assertAssetsAreExpanded();
        cy.get("@facets").filter(":visible").should("have.length", 17);

        cy.contains(".btn", "Show less filters").click();
        assertAssetsAreCollapsed();
        cy.get("@facets").filter(":visible").should("have.length", 10);

        listAllFacetsByFacetLine().should("eql", [
          [
            "Biblio status",
            "Classification",
            "Faculty",
            "Publication year",
            "Publication type",
          ],
          [
            "Publication status",
            "Librarian tags",
            "Message",
            "Locked",
            "UGent",
          ],
          [
            "WOS type",
            "VABB type",
            "File",
            "File type",
            "Created since",
            "Updated since",
            "Legacy",
          ],
        ]);
      });

      it("should not collapse the third line if one of its filters is active", () => {
        cy.visit("/publication");
        assertAssetsAreCollapsed();

        cy.contains(".btn", "Show all filters").click();
        assertAssetsAreExpanded();

        cy.contains(".dropdown", "Legacy")
          .click()
          .within(() => {
            cy.contains("label", "legacy publication").click();
            cy.contains(".btn", "Apply filter").click();
          });
        cy.getParams("f[legacy]").should("eq", "true");

        cy.contains(".dropdown", "Legacy")
          .contains(".badge-value", "legacy publication")
          .should("be.visible");
        assertAssetsAreExpanded();
      });

      it("should collapse the third line if no active filters remain", () => {
        cy.visit("/publication", {
          qs: {
            "f[status]": "private",
            "f[has_files]": "false",
            "f[legacy]": "false",
          },
        });
        assertAssetsAreExpanded();

        cy.contains(".dropdown", "Legacy")
          .click()
          .within(() => {
            cy.contains("Deselect all").click();
            cy.contains(".btn", "Apply filter").click();
          });
        cy.getParams("f[legacy]").should("not.exist");
        assertAssetsAreExpanded();

        cy.contains(".dropdown", "File")
          .click()
          .within(() => {
            cy.contains("Deselect all").click();
            cy.contains(".btn", "Apply filter").click();
          });
        cy.getParams("f[has_files]").should("not.exist");
        assertAssetsAreCollapsed();
      });

      it("should collapse the third line if all filters are reset", () => {
        cy.visit("/publication");
        assertAssetsAreCollapsed();

        cy.contains(".dropdown", "Faculty")
          .click()
          .within(() => {
            cy.contains("No affiliation").click();
            cy.contains(".btn", "Apply filter").click();
          });
        cy.getParams().should("eql", { q: "", "f[faculty_id]": "missing" });
        assertAssetsAreCollapsed();

        cy.contains(".dropdown", "Locked")
          .click()
          .within(() => {
            cy.contains("unlocked").click();
            cy.contains(".btn", "Apply filter").click();
          });
        cy.getParams().should("eql", {
          q: "",
          "f[faculty_id]": "missing",
          "f[locked]": "false",
        });
        assertAssetsAreCollapsed();

        cy.contains(".btn", "Show all filters").click();
        assertAssetsAreExpanded();

        cy.contains(".dropdown", "File type")
          .click()
          .within(() => {
            cy.contains("Full text").click();
            cy.contains(".btn", "Apply filter").click();
          });
        cy.getParams().should("eql", {
          q: "",
          "f[faculty_id]": "missing",
          "f[locked]": "false",
          "f[file_relation]": "main_file",
        });
        assertAssetsAreExpanded();

        cy.contains(".dropdown", "Updated since")
          .click()
          .within(() => {
            cy.get("input[type=text]").type("2022-03-04");
            cy.contains(".btn", "Apply filter").click();
          });
        cy.getParams().should("eql", {
          q: "",
          "f[faculty_id]": "missing",
          "f[locked]": "false",
          "f[file_relation]": "main_file",
          "f[updated_since]": "2022-03-04",
        });
        assertAssetsAreExpanded();

        cy.contains(".btn", "Reset filters").click();
        cy.getParams().should("eql", { q: "" });
        assertAssetsAreCollapsed();
      });

      it("should not show collapsible facet filters for datasets", () => {
        cy.visit("/dataset");

        assertNoCollapsibleFacetFilters();

        listAllFacetsByFacetLine().should("eql", [
          ["Biblio status", "Faculty", "Locked", "Persistent identifier type"],
          ["Librarian tags", "Message", "Created since", "Updated since"],
        ]);
      });
    });
  });

  function assertNoCollapsibleFacetFilters() {
    cy.contains(".btn", "Show all filters").should("not.exist");
    cy.get(".toggle-zone").should("not.exist");
    cy.get("#show-all-facet-filters-toggle").should("not.exist");
  }

  function assertAssetsAreCollapsed() {
    cy.get(".toggle-zone .badge-list:visible").should("have.length", 2);
    cy.contains(".btn", "Show all filters").should("be.visible");
    cy.contains(".btn", "Show less filters").should("not.be.visible");
  }

  function assertAssetsAreExpanded() {
    cy.get(".toggle-zone .badge-list:visible").should("have.length", 3);
    cy.contains(".btn", "Show all filters").should("not.be.visible");
    cy.contains(".btn", "Show less filters").should("be.visible");
  }

  function listAllFacetsByFacetLine() {
    return cy
      .contains(".bc-toolbar", "Reset filters")
      .find(".bc-toolbar-left .bc-toolbar-item .badge-list")
      .map<HTMLElement, string[]>((badgeList) =>
        Array.from(
          badgeList.querySelectorAll(".dropdown .badge .badge-text"),
        ).map((e) => e.textContent),
      );
  }
});
