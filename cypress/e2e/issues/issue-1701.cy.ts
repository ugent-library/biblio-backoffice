// https://github.com/ugent-library/biblio-backoffice/issues/1701

/**
 * The issue describes a concurrency case testable in multiple browser tabs. Since Cypress can't do that, we're gonna trigger some changes via AJAX calls
 * without HTMX swapping out the results as if this happened by another user or in another tab.
 */
describe("Issue #1701: Fix contributor removal concurrency bug", () => {
  beforeEach(() => {
    cy.loginAsResearcher("researcher1");
  });

  describe("Publication authors", () => {
    beforeEach(() => {
      cy.setUpPublication();
      cy.addAuthor("Author", "1", { external: true });
      cy.addAuthor("Author", "2", { external: true });

      cy.visitPublication();
      cy.contains(".nav-item", "People & Affiliations").click();

      verifyContributors("#contributors-author-table", [
        "Author 1",
        "Author 2",
      ]);
    });

    it("should remove the correct author when reordered by another user", () => {
      // Perform the reorder via AJAX
      cy.get("#contributors-author-table tbody")
        .invoke("attr", "hx-headers")
        .then(JSON.parse)
        .then((htmxValues) => {
          cy.get<string>("@biblioId").then((biblioId) => {
            // Reverse the authors
            const form = new URLSearchParams();
            form.append("position", "1");
            form.append("position", "0");

            cy.htmxRequest({
              method: "POST",
              url: `/publication/${biblioId}/contributors/author/order`,
              headers: htmxValues,
              form: true,
              body: form.toString(),
              log: true,
            });
          });
        });

      // Page is not refreshed so we still see the original order
      verifyContributors("#contributors-author-table", [
        "Author 1",
        "Author 2",
      ]);

      // Now let's try to delete Author 1
      cy.contains("tr", "Author 1").find(".btn .if-delete").click();
      cy.ensureModal("Confirm deletion").closeModal("Delete");

      cy.verifyConflictErrorDialog("publication").closeModal("Close");

      cy.reload();

      verifyContributors("#contributors-author-table", [
        "Author 2",
        "Author 1",
      ]);
    });

    it("should not remove two authors when two users simultaneously remove the same author", () => {
      // First let's remove Author 1 by AJAX call (as if another user did this)
      cy.contains("tr", "Author 1").find(".btn .if-delete").click();
      cy.ensureModal("Confirm deletion")
        .within(() => {
          cy.contains(".btn", "Delete").triggerHtmx("hx-delete");
        })
        .closeModal("Cancel");

      // Then remove Author 1 again by using the GUI
      cy.contains("tr", "Author 1").find(".btn .if-delete").click();
      cy.ensureModal("Confirm deletion").closeModal("Delete");

      cy.verifyConflictErrorDialog("publication").closeModal("Close");

      cy.reload();

      verifyContributors("#contributors-author-table", ["Author 2"]);
    });
  });

  describe("Publication editors", () => {
    beforeEach(() => {
      cy.setUpPublication();
      cy.addEditor("Editor", "1", { external: true });
      cy.addEditor("Editor", "2", { external: true });

      cy.visitPublication();
      cy.contains(".nav-item", "People & Affiliations").click();

      verifyContributors("#contributors-editor-table", [
        "Editor 1",
        "Editor 2",
      ]);
    });

    it("should remove the correct editor when reordered by another user", () => {
      // Perform the reorder via AJAX
      cy.get("#contributors-editor-table tbody")
        .invoke("attr", "hx-headers")
        .then(JSON.parse)
        .then((htmxValues) => {
          cy.get<string>("@biblioId").then((biblioId) => {
            // Reverse the editors
            const form = new URLSearchParams();
            form.append("position", "1");
            form.append("position", "0");

            cy.htmxRequest({
              method: "POST",
              url: `/publication/${biblioId}/contributors/editor/order`,
              headers: htmxValues,
              form: true,
              body: form.toString(),
              log: true,
            });
          });
        });

      // Page is not refreshed so we still see the original order
      verifyContributors("#contributors-editor-table", [
        "Editor 1",
        "Editor 2",
      ]);

      // Now let's try to delete Editor 1
      cy.contains("tr", "Editor 1").find(".btn .if-delete").click();
      cy.ensureModal("Confirm deletion").closeModal("Delete");

      cy.verifyConflictErrorDialog("publication").closeModal("Close");

      cy.reload();

      verifyContributors("#contributors-editor-table", [
        "Editor 2",
        "Editor 1",
      ]);
    });

    it("should not remove two editors when two users simultaneously remove the same editor", () => {
      // First let's remove Editor 1 by AJAX call (as if another user did this)
      cy.contains("tr", "Editor 1").find(".btn .if-delete").click();
      cy.ensureModal("Confirm deletion")
        .within(() => {
          cy.contains(".btn", "Delete").triggerHtmx("hx-delete");
        })
        .closeModal("Cancel");

      // Then remove Editor 1 again by using the GUI
      cy.contains("tr", "Editor 1").find(".btn .if-delete").click();
      cy.ensureModal("Confirm deletion").closeModal("Delete");

      verifyContributors("#contributors-editor-table", ["Editor 2"]);
    });
  });

  describe("Publication supervisors", () => {
    beforeEach(() => {
      cy.setUpPublication("Dissertation");
      cy.addSupervisor("Supervisor", "1", { external: true });
      cy.addSupervisor("Supervisor", "2", { external: true });

      cy.visitPublication();
      cy.contains(".nav-item", "People & Affiliations").click();

      verifyContributors("#contributors-supervisor-table", [
        "Supervisor 1",
        "Supervisor 2",
      ]);
    });

    it("should remove the correct supervisor when reordered by another user", () => {
      // Perform the reorder via AJAX
      cy.get("#contributors-supervisor-table tbody")
        .invoke("attr", "hx-headers")
        .then(JSON.parse)
        .then((htmxValues) => {
          cy.get<string>("@biblioId").then((biblioId) => {
            // Reverse the supervisors
            const form = new URLSearchParams();
            form.append("position", "1");
            form.append("position", "0");

            cy.htmxRequest({
              method: "POST",
              url: `/publication/${biblioId}/contributors/supervisor/order`,
              headers: htmxValues,
              form: true,
              body: form.toString(),
              log: true,
            });
          });
        });

      // Page is not refreshed so we still see the original order
      verifyContributors("#contributors-supervisor-table", [
        "Supervisor 1",
        "Supervisor 2",
      ]);

      // Now let's try to delete Supervisor 1
      cy.contains("tr", "Supervisor 1").find(".btn .if-delete").click();
      cy.ensureModal("Confirm deletion").closeModal("Delete");

      cy.verifyConflictErrorDialog("publication").closeModal("Close");

      cy.reload();

      verifyContributors("#contributors-supervisor-table", [
        "Supervisor 2",
        "Supervisor 1",
      ]);
    });

    it("should not remove two supervisors when two users simultaneously remove the same supervisor", () => {
      // First let's remove Supervisor 1 by AJAX call (as if another user did this)
      cy.contains("tr", "Supervisor 1").find(".btn .if-delete").click();
      cy.ensureModal("Confirm deletion")
        .within(() => {
          cy.contains(".btn", "Delete").triggerHtmx("hx-delete");
        })
        .closeModal("Cancel");

      // Then remove Supervisor 1 again by using the GUI
      cy.contains("tr", "Supervisor 1").find(".btn .if-delete").click();
      cy.ensureModal("Confirm deletion").closeModal("Delete");

      verifyContributors("#contributors-supervisor-table", ["Supervisor 2"]);
    });
  });

  describe("Dataset creators", () => {
    beforeEach(() => {
      cy.setUpDataset();
      cy.addCreator("Creator", "1", { external: true });
      cy.addCreator("Creator", "2", { external: true });

      cy.visitDataset();
      cy.contains(".nav-item", "People & Affiliations").click();

      verifyContributors("#contributors-author-table", [
        "Creator 1",
        "Creator 2",
      ]);
    });

    it("should remove the correct creator when reordered by another user", () => {
      // Perform the reorder via AJAX
      cy.get("#contributors-author-table tbody")
        .invoke("attr", "hx-headers")
        .then(JSON.parse)
        .then((htmxValues) => {
          cy.get<string>("@biblioId").then((biblioId) => {
            // Reverse the authors
            const form = new URLSearchParams();
            form.append("position", "1");
            form.append("position", "0");

            cy.htmxRequest({
              method: "POST",
              url: `/dataset/${biblioId}/contributors/author/order`,
              headers: htmxValues,
              form: true,
              body: form.toString(),
              log: true,
            });
          });
        });

      // Page is not refreshed so we still see the original order
      verifyContributors("#contributors-author-table", [
        "Creator 1",
        "Creator 2",
      ]);

      // Now let's try to delete Creator 1
      cy.contains("tr", "Creator 1").find(".btn .if-delete").click();
      cy.ensureModal("Confirm deletion").closeModal("Delete");

      cy.verifyConflictErrorDialog("publication").closeModal("Close");

      cy.reload();

      verifyContributors("#contributors-author-table", [
        "Creator 2",
        "Creator 1",
      ]);
    });

    it("should not remove two creators when two users simultaneously remove the same creator", () => {
      // First let's remove Creator 1 by AJAX call (as if another user did this)
      cy.contains("tr", "Creator 1").find(".btn .if-delete").click();
      cy.ensureModal("Confirm deletion")
        .within(() => {
          cy.contains(".btn", "Delete").triggerHtmx("hx-delete");
        })
        .closeModal("Cancel");

      // Then remove Creator 1 again by using the GUI
      cy.contains("tr", "Creator 1").find(".btn .if-delete").click();
      cy.ensureModal("Confirm deletion").closeModal("Delete");

      verifyContributors("#contributors-author-table", ["Creator 2"]);
    });
  });

  function verifyContributors(tableSelector: string, contributors: string[]) {
    cy.get(tableSelector)
      .find("tbody tr")
      .should("have.length", contributors.length)
      .map((tr: HTMLTableRowElement) => tr.querySelector("td").textContent)
      .should("eql", contributors);
  }
});