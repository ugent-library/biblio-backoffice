type Breadcrumbs = {
  name?: string;
  url?: string;
};

describe("Site breadcrumbs", () => {
  beforeEach(() => {
    cy.login("researcher1");
  });

  it("should have no additional breadcrumbs on the anonymous homepage", () => {
    cy.clearAllCookies();

    cy.visit("/");

    verifyBreadcrumbs();
  });

  describe("The dashboards", () => {
    describe("As researcher", () => {
      it("should have breadcrumbs on the dashboard", () => {
        cy.visit("/dashboard");

        verifyBreadcrumbs({ name: "Dashboard" });
      });
    });

    describe("As Librarian", () => {
      it("should have breadcrumbs on the dashboard", () => {
        cy.login("librarian1");
        cy.switchMode("Librarian");

        cy.visit("/dashboard/publications/faculties");
        verifyBreadcrumbs({ name: "Dashboard" });

        cy.visit("/dashboard/publications/socs");
        verifyBreadcrumbs({ name: "Dashboard" });

        cy.visit("/dashboard/datasets/faculties");
        verifyBreadcrumbs({ name: "Dashboard" });

        cy.visit("/dashboard/datasets/socs");
        verifyBreadcrumbs({ name: "Dashboard" });
      });
    });
  });

  describe("Publications", () => {
    it("should have breadcrumbs on the publications page", () => {
      cy.visit("/publication");

      verifyBreadcrumbs({ name: "Publications" });
    });

    it("should have breadcrumbs on the publication detail page", () => {
      cy.setUpPublication();
      cy.visitPublication();

      verifyBreadcrumbs(
        { name: "Publications", url: "/publication" },
        { name: "Publication detail" },
      );
    });

    describe("Publication import flow", () => {
      const MULTIPLE_IMPORT_PUBLICATION_DETAIL_REGEX =
        /^\/add-publication\/import\/multiple\/[A-Z0-9]+\/publication\/(?<biblioId>[A-Z0-9]+)$/;

      const assertions = [
        { name: "Publications", url: "/publication" },
        { name: "New publication" },
      ];

      it("should have breadcrumbs on the add publication page", () => {
        cy.visit("/add-publication");
        verifyBreadcrumbs(...assertions);
      });

      it("should have breadcrumbs in the Web of Science flow", () => {
        cy.visit("/add-publication?method=wos");
        verifyBreadcrumbs(...assertions);
        cy.get("input[name=file]").selectFile(
          "cypress/fixtures/import-from-wos-single.txt",
        );
        cy.location("pathname").should(
          "match",
          /^\/add-publication\/import\/multiple\/[A-Z0-9]+\/confirm$/,
        );
        verifyBreadcrumbs(...assertions);

        cy.contains(".list-group-item", "Description").click();
        cy.location("pathname")
          .should("match", MULTIPLE_IMPORT_PUBLICATION_DETAIL_REGEX)
          .then(
            (pathname) =>
              pathname.match(MULTIPLE_IMPORT_PUBLICATION_DETAIL_REGEX)!.groups
                .biblioId!,
          )
          .as("biblioId");
        verifyBreadcrumbs(...assertions);
        cy.contains('Back to "Review and publish" overview').click();

        cy.contains(".list-group-item", "People & Affiliations").click();
        cy.location("pathname").should(
          "match",
          MULTIPLE_IMPORT_PUBLICATION_DETAIL_REGEX,
        );
        verifyBreadcrumbs(...assertions);
        cy.contains('Back to "Review and publish" overview').click();

        cy.contains(".list-group-item", "Full text & Files").click();
        cy.location("pathname").should(
          "match",
          MULTIPLE_IMPORT_PUBLICATION_DETAIL_REGEX,
        );
        verifyBreadcrumbs(...assertions);
        cy.contains('Back to "Review and publish" overview').click();

        cy.addAuthor("John", "Doe");
        cy.contains(".btn", "Publish all to Biblio").click();
        cy.location("pathname").should(
          "match",
          /^\/add-publication\/import\/multiple\/[A-Z0-9]+\/finish$/,
        );
        verifyBreadcrumbs(...assertions);
      });

      it("should have breadcrumbs in the import via identifier flow", () => {
        const DOI = "10.7202/1041023ar";

        cy.deletePublications(DOI);

        cy.visit("/add-publication?method=identifier");
        verifyBreadcrumbs(...assertions);

        cy.get('input[name="identifier"]').type(DOI);
        cy.contains(".btn", "Add publication(s)").click();
        cy.location("pathname").should(
          "eq",
          "/add-publication/import/single/confirm",
        );
        verifyBreadcrumbs(...assertions);

        cy.contains(".btn", "Complete Description").click();
        cy.location("pathname")
          .should("match", /^\/publication\/[A-Z0-9]+\/add\/confirm$/)
          .then(
            (pathname) =>
              pathname.match(
                /^\/publication\/(?<biblioId>[A-Z0-9]+)\/add\/confirm$/,
              )!.groups.biblioId!,
          )
          .as("biblioId");
        verifyBreadcrumbs(...assertions);

        cy.addAuthor("John", "Doe");
        cy.reload(); // Get latest snapshot ID
        cy.contains(".btn", "Publish to Biblio").click();
        cy.location("pathname").should(
          "match",
          /^\/publication\/[A-Z0-9]+\/add\/finish$/,
        );
        verifyBreadcrumbs(...assertions);
      });

      it("should have breadcrumbs in the manual flow", () => {
        cy.visit("/add-publication?method=manual");
        verifyBreadcrumbs(...assertions);

        cy.contains("Miscellaneous").click();
        cy.contains(".btn", "Add publication(s)").click();
        cy.location("pathname").should(
          "eq",
          "/add-publication/import/single/confirm",
        );
        verifyBreadcrumbs(...assertions);

        cy.updateFields(
          "Publication details",
          () => {
            cy.setFieldByLabel("Title", "The publication title");
            cy.setFieldByLabel(
              "Publication year",
              new Date().getFullYear().toString(),
            );
          },
          true,
        );

        cy.contains(".btn", "Complete Description").click();
        cy.location("pathname")
          .should("match", /^\/publication\/[A-Z0-9]+\/add\/confirm$/)
          .then(
            (pathname) =>
              pathname.match(
                /^\/publication\/(?<biblioId>[A-Z0-9]+)\/add\/confirm$/,
              )!.groups.biblioId,
          )
          .as("biblioId");
        verifyBreadcrumbs(...assertions);

        cy.addAuthor("John", "Doe");
        cy.reload(); // Get latest snapshot ID
        cy.contains(".btn", "Publish to Biblio").click();
        cy.location("pathname").should(
          "match",
          /^\/publication\/[A-Z0-9]+\/add\/finish$/,
        );
        verifyBreadcrumbs(...assertions);
      });

      it("should have breadcrumbs in the import from BibTeX flow", () => {
        cy.visit("/add-publication?method=bibtex");
        verifyBreadcrumbs(...assertions);

        cy.get("input[name=file]").selectFile(
          "cypress/fixtures/import-from-bibtex.bib",
        );
        cy.location("pathname").should(
          "match",
          /^\/add-publication\/import\/multiple\/[A-Z0-9]+\/confirm$/,
        );
        verifyBreadcrumbs(...assertions);

        cy.contains(".list-group-item", "Description").click();
        cy.location("pathname")
          .should("match", MULTIPLE_IMPORT_PUBLICATION_DETAIL_REGEX)
          .then(
            (pathname) =>
              pathname.match(MULTIPLE_IMPORT_PUBLICATION_DETAIL_REGEX)!.groups
                .biblioId!,
          )
          .as("biblioId");
        verifyBreadcrumbs(...assertions);
        cy.contains('Back to "Review and publish" overview').click();

        cy.contains(".list-group-item", "People & Affiliations").click();
        cy.location("pathname").should(
          "match",
          MULTIPLE_IMPORT_PUBLICATION_DETAIL_REGEX,
        );
        verifyBreadcrumbs(...assertions);
        cy.contains('Back to "Review and publish" overview').click();

        cy.contains(".list-group-item", "Full text & Files").click();
        cy.location("pathname").should(
          "match",
          MULTIPLE_IMPORT_PUBLICATION_DETAIL_REGEX,
        );
        verifyBreadcrumbs(...assertions);
        cy.contains('Back to "Review and publish" overview').click();

        cy.addAuthor("John", "Doe");
        cy.contains(".btn", "Publish all to Biblio").click();
        cy.location("pathname").should(
          "match",
          /^\/add-publication\/import\/multiple\/[A-Z0-9]+\/finish$/,
        );
        verifyBreadcrumbs(...assertions);
      });
    });
  });

  describe("Datasets", () => {
    it("should have breadcrumbs on the datasets page", () => {
      cy.visit("/dataset");

      verifyBreadcrumbs({ name: "Datasets" });
    });

    it("should have breadcrumbs on the dataset detail page", () => {
      cy.setUpDataset();
      cy.visitDataset();

      verifyBreadcrumbs(
        { name: "Datasets", url: "/dataset" },
        { name: "Dataset detail" },
      );
    });

    describe("Dataset import flow", () => {
      const assertions = [
        { name: "Datasets", url: "/dataset" },
        { name: "New dataset" },
      ];

      it("should have breadcrumbs on the add dataset page", () => {
        cy.visit("/add-dataset");
        verifyBreadcrumbs(...assertions);
      });

      it("should have breadcrumbs in the import via identifier flow", () => {
        const DOI = "10.48804/A76XM9";

        cy.deleteDatasets(DOI);

        cy.visit("/add-dataset");
        cy.contains("h6", "Register your dataset via a DOI").click();
        cy.contains(".btn", "Add dataset").click();
        verifyBreadcrumbs(...assertions);

        cy.get("input[name=identifier]").type(DOI);
        cy.contains(".btn", "Add dataset").click();
        cy.location("pathname").should("eq", "/add-dataset/import/confirm");
        verifyBreadcrumbs(...assertions);

        cy.updateFields(
          "Dataset details",
          () => {
            cy.setFieldByLabel("Data format", "text/csv");
            cy.contains(".badge", "text/csv").click();
            cy.setFieldByLabel("License", "CC BY (4.0)");
            cy.setFieldByLabel("Access level", "Closed access");
          },
          true,
        );

        cy.contains(".btn", "Complete Description").click();
        cy.location("pathname")
          .should("match", /^\/dataset\/[A-Z0-9]+\/add\/confirm$/)
          .then(
            (pathname) =>
              pathname.match(
                /^\/dataset\/(?<biblioId>[A-Z0-9]+)\/add\/confirm$/,
              )!.groups.biblioId!,
          )
          .as("biblioId");
        verifyBreadcrumbs(...assertions);

        cy.addCreator("John", "Doe");
        cy.reload(); // Get latest snapshot ID
        cy.contains(".btn", "Publish to Biblio").click();
        cy.location("pathname").should(
          "match",
          /^\/dataset\/[A-Z0-9]+\/add\/finish$/,
        );
        verifyBreadcrumbs(...assertions);
      });

      it("should have breadcrumbs in the manual flow", () => {
        cy.visit("/add-dataset");

        cy.contains("h6", "Register a dataset manually").click();
        cy.contains(".btn", "Add dataset").click();
        verifyBreadcrumbs(...assertions);

        cy.updateFields(
          "Dataset details",
          () => {
            cy.setFieldByLabel("Title", "The dataset title");
            cy.setFieldByLabel("Persistent identifier type", "DOI");
            cy.setFieldByLabel("Identifier", "ABC-123");
            cy.setFieldByLabel(
              "Publication year",
              new Date().getFullYear().toString(),
            );
            cy.setFieldByLabel("Publisher", "UGent");
            cy.setFieldByLabel("Data format", "text/csv");
            cy.contains(".badge", "text/csv").click();
            cy.setFieldByLabel("License", "CC BY (4.0)");
            cy.setFieldByLabel("Access level", "Closed access");
          },
          true,
        );

        cy.contains(".btn", "Complete Description").click();
        cy.location("pathname")
          .should("match", /^\/dataset\/[A-Z0-9]+\/add\/confirm$/)
          .then(
            (pathname) =>
              pathname.match(
                /^\/dataset\/(?<biblioId>[A-Z0-9]+)\/add\/confirm$/,
              )!.groups.biblioId,
          )
          .as("biblioId");
        verifyBreadcrumbs(...assertions);

        cy.addCreator("John", "Doe");
        cy.reload(); // Get latest snapshot ID
        cy.contains(".btn", "Publish to Biblio").click();
        cy.location("pathname").should(
          "match",
          /^\/dataset\/[A-Z0-9]+\/add\/finish$/,
        );
        verifyBreadcrumbs(...assertions);
      });
    });
  });

  describe("Curator features", () => {
    beforeEach(() => {
      cy.login("librarian1");
      cy.switchMode("Librarian");
    });

    it("should have breadcrumbs on the batch page", () => {
      cy.visit("/publication/batch");

      verifyBreadcrumbs({ name: "Batch import" });
    });

    it("should have breadcrumbs on the suggestions page", () => {
      cy.visit("/candidate-records");

      verifyBreadcrumbs({ name: "Suggestions" });
    });
  });

  function verifyBreadcrumbs(...breadcrumbs: Breadcrumbs[]) {
    breadcrumbs.unshift(
      { url: "/" }, // The logo
      { name: "Biblio", url: "/" }, // The Biblio link
    );

    cy.get(".breadcrumb .breadcrumb-item")
      .should("have.length", breadcrumbs.length)
      .each(($breadcrumb, i) => {
        if (breadcrumbs[i].name) {
          expect($breadcrumb).to.contain(breadcrumbs[i].name);
        }

        if (breadcrumbs[i].url) {
          expect($breadcrumb).to.have.descendants("a");
          expect($breadcrumb.find("a")).to.have.attr(
            "href",
            breadcrumbs[i].url,
          );
        }
      });
  }
});
