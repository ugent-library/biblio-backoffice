// https://github.com/ugent-library/biblio-backoffice/issues/1278

describe("Issue #1278: [Plato imports] As a researcher or supervisor, I can see how I am related to a record (supervisor, author, editor, ...) in the card", () => {
  beforeEach(() => {
    cy.loginAsResearcher();
  });

  describe("for publications", () => {
    it("should display my role when I'm an author", () => {
      cy.setUpPublication();
      verifyMyRoles("publication", "registrar");

      // Add myself as author
      cy.visitPublication();
      cy.updateFields(
        "Authors",
        () => {
          cy.setFieldByLabel("First name", "Biblio");
          cy.setFieldByLabel("Last name", "Researcher");
          cy.contains(".btn", "Add author").click();
        },
        true,
      );
      verifyMyRoles("publication", "author");

      cy.loginAsLibrarian();
      cy.switchMode("Librarian");
      verifyMyRoles("publication");
    });

    it("should display my role when I'm a supervisor", () => {
      cy.setUpPublication("Dissertation");
      verifyMyRoles("publication", "registrar");

      // Add myself as supervisor
      cy.visitPublication();
      cy.updateFields(
        "Supervisors",
        () => {
          cy.setFieldByLabel("First name", "Biblio");
          cy.setFieldByLabel("Last name", "Researcher");
          cy.contains(".btn", "Add supervisor").click();
        },
        true,
      );
      verifyMyRoles("publication", "supervisor");

      cy.loginAsLibrarian();
      cy.switchMode("Librarian");
      verifyMyRoles("publication");
    });

    it("should display my roles when I'm both an author and a supervisor", () => {
      cy.setUpPublication("Dissertation");
      verifyMyRoles("publication", "registrar");

      // Add myself as supervisor
      cy.visitPublication();
      cy.updateFields(
        "Supervisors",
        () => {
          cy.setFieldByLabel("First name", "Biblio");
          cy.setFieldByLabel("Last name", "Researcher");
          cy.contains(".btn", "Add supervisor").click();
        },
        true,
      );
      verifyMyRoles("publication", "supervisor");

      // Add myself as author
      cy.visitPublication();
      cy.updateFields(
        "Authors",
        () => {
          cy.setFieldByLabel("First name", "Biblio");
          cy.setFieldByLabel("Last name", "Researcher");
          cy.contains(".btn", "Add author").click();
        },
        true,
      );
      verifyMyRoles("publication", "author", "supervisor");

      cy.loginAsLibrarian();
      cy.switchMode("Librarian");
      verifyMyRoles("publication");
    });

    it("should display my role when I'm only a registrar", function () {
      cy.setUpPublication();
      verifyMyRoles("publication", "registrar");

      // Add other author
      cy.visitPublication();
      cy.updateFields(
        "Authors",
        () => {
          cy.setFieldByLabel("First name", "John");
          cy.setFieldByLabel("Last name", "Doe");
          cy.contains(".btn", "Add author").click();
        },
        true,
      );
      verifyMyRoles("publication", "registrar");

      // Add myself as editor
      cy.visitPublication();
      cy.updateFields(
        "Editors",
        () => {
          cy.setFieldByLabel("First name", "Biblio");
          cy.setFieldByLabel("Last name", "Researcher");
          cy.contains(".btn", "Add editor").click();
        },
        true,
      );
      verifyMyRoles("publication", "registrar");

      cy.loginAsLibrarian();
      cy.switchMode("Librarian");
      verifyMyRoles("publication");
    });
  });

  describe("for datasets", () => {
    it("should display my role when I'm a creator", () => {
      cy.setUpDataset();
      verifyMyRoles("dataset", "registrar");

      // Add myself as creator
      cy.visitDataset();
      cy.updateFields(
        "Creators",
        () => {
          cy.setFieldByLabel("First name", "Biblio");
          cy.setFieldByLabel("Last name", "Researcher");
          cy.contains(".btn", "Add creator").click();
        },
        true,
      );
      verifyMyRoles("dataset", "creator");

      cy.loginAsLibrarian();
      cy.switchMode("Librarian");
      verifyMyRoles("dataset");
    });

    it("should display my role when I'm only a registrar", () => {
      cy.setUpDataset();
      verifyMyRoles("dataset", "registrar");

      // Add other creator
      cy.visitDataset();
      cy.updateFields(
        "Creators",
        () => {
          cy.setFieldByLabel("First name", "John");
          cy.setFieldByLabel("Last name", "Doe");
          cy.contains(".btn", "Add creator").click();
        },
        true,
      );
      verifyMyRoles("dataset", "registrar");

      cy.loginAsLibrarian();
      cy.switchMode("Librarian");
      verifyMyRoles("dataset");
    });
  });

  type UserRole = "author" | "supervisor" | "creator" | "registrar";

  function verifyMyRoles(
    scope: "publication" | "dataset",
    ...expectedRoles: UserRole[]
  ) {
    cy.get<string>("@biblioId").then((biblioId) => {
      cy.visit(`/${scope}`, { qs: { q: biblioId } });
    });

    cy.get(".card .card-header").should("contain", `Showing 1 ${scope}s`);
    cy.get(".card .card-body .list-group .list-group-item").should(
      "have.length",
      1,
    );

    const selector = cy
      .get(".card .card-body .list-group .list-group-item")
      .should("have.length", 1)
      .find(".c-author-list .c-author .badge");
    if (expectedRoles.length) {
      selector.should("contain", `Your role: ${expectedRoles.join(", ")}`);
    } else {
      selector.should("not.exist");
    }
  }
});
