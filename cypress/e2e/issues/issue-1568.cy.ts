// https://github.com/ugent-library/biblio-backoffice/issues/1568

import "../../support/commands/drag";

/**
 * The issue describes a concurrency case testable in multiple browser tabs. Since Cypress can't do that, we're gonna trigger some changes via AJAX calls
 * without HTMX swapping out the results as if this happened by another user or in another tab.
 */
describe("Issue #1568: Missing conflict handling when editing / adding / removing contributors", () => {
  beforeEach(() => {
    cy.loginAsResearcher("researcher1");
  });

  describe("Publication authors", () => {
    beforeEach(() => {
      cy.setUpPublication();
      cy.addAuthor("Author", "1", { external: true });
      cy.addAuthor("Author", "2", { external: true });
      cy.addAuthor("Author", "3", { external: true });

      cy.visitPublication();
      cy.contains(".nav-item", "People & Affiliations").click();

      verifyContributors("#contributors-author-table", [
        "Author 1",
        "Author 2",
        "Author 3",
      ]);
    });

    it("should show a conflict error when reordering publication authors after reordering them in another session", () => {
      // Perform the reorder via AJAX
      cy.get("#contributors-author-table tbody")
        .invoke("attr", "hx-headers")
        .then(JSON.parse)
        .then((htmxValues) => {
          cy.get<string>("@biblioId").then((biblioId) => {
            // Reverse the authors
            const form = new URLSearchParams();
            form.append("position", "2");
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
        "Author 3",
      ]);

      cy.intercept("POST", "/publication/*/contributors/author/order").as(
        "reorderAuthors",
      );

      // Perform the reorder
      cy.get("#author-1 .sortable-handle").drag("#author-2 .sortable-handle");

      verifyContributors("#contributors-author-table", [
        "Author 1",
        "Author 3",
        "Author 2",
      ]);

      cy.wait("@reorderAuthors").should(
        "have.nested.property",
        "response.statusCode",
        200,
      );

      cy.verifyConflictErrorDialog("publication").closeModal("Close");

      // Verify that nothing was reordered
      cy.reload();

      verifyContributors("#contributors-author-table", [
        "Author 3",
        "Author 2",
        "Author 1",
      ]);
    });

    it("should show a conflict error when reordering publication authors after removing one in another session", () => {
      // Delete the second author by AJAX call and cancel the confirm delete dialog
      cy.contains("tr", "Author 2").find(".btn .if-delete").click();
      cy.ensureModal("Confirm deletion")
        .within(() => {
          cy.contains(".btn", "Delete").triggerHtmx("hx-delete");
        })
        .closeModal("Cancel");
      cy.ensureNoModal();

      // Page is not refreshed so we still think there are 3 authors
      verifyContributors("#contributors-author-table", [
        "Author 1",
        "Author 2",
        "Author 3",
      ]);

      cy.intercept("POST", "/publication/*/contributors/author/order").as(
        "reorderAuthors",
      );

      // Perform the reorder
      cy.get("#author-1 .sortable-handle").drag("#author-2 .sortable-handle");

      verifyContributors("#contributors-author-table", [
        "Author 1",
        "Author 3",
        "Author 2",
      ]);

      cy.wait("@reorderAuthors").should(
        "have.nested.property",
        "response.statusCode",
        200,
      );

      cy.verifyConflictErrorDialog("publication").closeModal("Close");

      // Verify that nothing was reordered
      cy.reload();

      verifyContributors("#contributors-author-table", [
        "Author 1",
        "Author 3",
      ]);
    });

    it("should show a conflict error when reordering publication authors after adding one in another session", () => {
      // Add a fourth author by AJAX call
      cy.addAuthor("Author", "4", { external: true });

      cy.intercept("POST", "/publication/*/contributors/author/order").as(
        "reorderAuthors",
      );

      // Perform the reorder
      cy.get("#author-0 .sortable-handle").drag("#author-1 .sortable-handle");

      verifyContributors("#contributors-author-table", [
        "Author 2",
        "Author 1",
        "Author 3",
      ]);

      cy.wait("@reorderAuthors").should(
        "have.nested.property",
        "response.statusCode",
        200,
      );

      cy.verifyConflictErrorDialog("publication").closeModal("Close");

      // Verify that nothing was reordered
      cy.reload();

      verifyContributors("#contributors-author-table", [
        "Author 1",
        "Author 2",
        "Author 3",
        "Author 4",
      ]);
    });
  });

  describe("Publication editors", () => {
    beforeEach(() => {
      cy.setUpPublication();
      cy.addEditor("Editor", "1", { external: true });
      cy.addEditor("Editor", "2", { external: true });
      cy.addEditor("Editor", "3", { external: true });

      cy.visitPublication();
      cy.contains(".nav-item", "People & Affiliations").click();

      verifyContributors("#contributors-editor-table", [
        "Editor 1",
        "Editor 2",
        "Editor 3",
      ]);
    });

    it("should show a conflict error when reordering publication editors after reordering them in another session", () => {
      // Perform the reorder via AJAX
      cy.get("#contributors-editor-table tbody")
        .invoke("attr", "hx-headers")
        .then(JSON.parse)
        .then((htmxValues) => {
          cy.get<string>("@biblioId").then((biblioId) => {
            // Reverse the editors
            const form = new URLSearchParams();
            form.append("position", "2");
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
        "Editor 3",
      ]);

      cy.intercept("POST", "/publication/*/contributors/editor/order").as(
        "reorderEditors",
      );

      // Perform the reorder
      cy.get("#editor-1 .sortable-handle").drag("#editor-2 .sortable-handle");

      verifyContributors("#contributors-editor-table", [
        "Editor 1",
        "Editor 3",
        "Editor 2",
      ]);

      cy.wait("@reorderEditors").should(
        "have.nested.property",
        "response.statusCode",
        200,
      );

      cy.verifyConflictErrorDialog("publication").closeModal("Close");

      // Verify that nothing was reordered
      cy.reload();

      verifyContributors("#contributors-editor-table", [
        "Editor 3",
        "Editor 2",
        "Editor 1",
      ]);
    });

    it("should show a conflict error when reordering publication editors after removing one in another session", () => {
      // Delete the second editor by AJAX call and cancel the confirm delete dialog
      cy.contains("tr", "Editor 2").find(".btn .if-delete").click();
      cy.ensureModal("Confirm deletion")
        .within(() => {
          cy.contains(".btn", "Delete").triggerHtmx("hx-delete");
        })
        .closeModal("Cancel");
      cy.ensureNoModal();

      // Page is not refreshed so we still think there are 3 editors
      verifyContributors("#contributors-editor-table", [
        "Editor 1",
        "Editor 2",
        "Editor 3",
      ]);

      cy.intercept("POST", "/publication/*/contributors/editor/order").as(
        "reorderEditors",
      );

      // Perform the reorder
      cy.get("#editor-1 .sortable-handle").drag("#editor-2 .sortable-handle");

      verifyContributors("#contributors-editor-table", [
        "Editor 1",
        "Editor 3",
        "Editor 2",
      ]);

      cy.wait("@reorderEditors").should(
        "have.nested.property",
        "response.statusCode",
        200,
      );

      cy.verifyConflictErrorDialog("publication").closeModal("Close");

      // Verify that nothing was reordered
      cy.reload();

      verifyContributors("#contributors-editor-table", [
        "Editor 1",
        "Editor 3",
      ]);
    });

    it("should show a conflict error when reordering publication editors after adding one in another session", () => {
      // Add a fourth editor by AJAX call
      cy.addEditor("Editor", "4", { external: true });

      cy.intercept("POST", "/publication/*/contributors/editor/order").as(
        "reorderEditors",
      );

      // Perform the reorder
      cy.get("#editor-0 .sortable-handle").drag("#editor-1 .sortable-handle");

      verifyContributors("#contributors-editor-table", [
        "Editor 2",
        "Editor 1",
        "Editor 3",
      ]);

      cy.wait("@reorderEditors").should(
        "have.nested.property",
        "response.statusCode",
        200,
      );

      cy.verifyConflictErrorDialog("publication").closeModal("Close");

      // Verify that nothing was reordered
      cy.reload();

      verifyContributors("#contributors-editor-table", [
        "Editor 1",
        "Editor 2",
        "Editor 3",
        "Editor 4",
      ]);
    });
  });

  describe("Publication supervisors", () => {
    beforeEach(() => {
      cy.setUpPublication("Dissertation");
      cy.addSupervisor("Supervisor", "1", { external: true });
      cy.addSupervisor("Supervisor", "2", { external: true });
      cy.addSupervisor("Supervisor", "3", { external: true });

      cy.visitPublication();
      cy.contains(".nav-item", "People & Affiliations").click();

      verifyContributors("#contributors-supervisor-table", [
        "Supervisor 1",
        "Supervisor 2",
        "Supervisor 3",
      ]);
    });

    it("should show a conflict error when reordering publication supervisors after reordering them in another session", () => {
      // Perform the reorder via AJAX
      cy.get("#contributors-supervisor-table tbody")
        .invoke("attr", "hx-headers")
        .then(JSON.parse)
        .then((htmxValues) => {
          cy.get<string>("@biblioId").then((biblioId) => {
            // Reverse the supervisors
            const form = new URLSearchParams();
            form.append("position", "2");
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
        "Supervisor 3",
      ]);

      cy.intercept("POST", "/publication/*/contributors/supervisor/order").as(
        "reorderSupervisors",
      );

      // Perform the reorder
      cy.get("#supervisor-1 .sortable-handle").drag(
        "#supervisor-2 .sortable-handle",
      );

      verifyContributors("#contributors-supervisor-table", [
        "Supervisor 1",
        "Supervisor 3",
        "Supervisor 2",
      ]);

      cy.wait("@reorderSupervisors").should(
        "have.nested.property",
        "response.statusCode",
        200,
      );

      cy.verifyConflictErrorDialog("publication").closeModal("Close");

      // Verify that nothing was reordered
      cy.reload();

      verifyContributors("#contributors-supervisor-table", [
        "Supervisor 3",
        "Supervisor 2",
        "Supervisor 1",
      ]);
    });

    it("should show a conflict error when reordering publication supervisors after removing one in another session", () => {
      // Delete the second supervisor by AJAX call and cancel the confirm delete dialog
      cy.contains("tr", "Supervisor 2").find(".btn .if-delete").click();
      cy.ensureModal("Confirm deletion")
        .within(() => {
          cy.contains(".btn", "Delete").triggerHtmx("hx-delete");
        })
        .closeModal("Cancel");
      cy.ensureNoModal();

      // Page is not refreshed so we still think there are 3 supervisors
      verifyContributors("#contributors-supervisor-table", [
        "Supervisor 1",
        "Supervisor 2",
        "Supervisor 3",
      ]);

      cy.intercept("POST", "/publication/*/contributors/supervisor/order").as(
        "reorderSupervisors",
      );

      // Perform the reorder
      cy.get("#supervisor-1 .sortable-handle").drag(
        "#supervisor-2 .sortable-handle",
      );

      verifyContributors("#contributors-supervisor-table", [
        "Supervisor 1",
        "Supervisor 3",
        "Supervisor 2",
      ]);

      cy.wait("@reorderSupervisors").should(
        "have.nested.property",
        "response.statusCode",
        200,
      );

      cy.verifyConflictErrorDialog("publication").closeModal("Close");

      // Verify that nothing was reordered
      cy.reload();

      verifyContributors("#contributors-supervisor-table", [
        "Supervisor 1",
        "Supervisor 3",
      ]);
    });

    it("should show a conflict error when reordering publication supervisors after adding one in another session", () => {
      // Add a fourth supervisor by AJAX call
      cy.addSupervisor("Supervisor", "4", { external: true });

      cy.intercept("POST", "/publication/*/contributors/supervisor/order").as(
        "reorderSupervisors",
      );

      // Perform the reorder
      cy.get("#supervisor-0 .sortable-handle").drag(
        "#supervisor-1 .sortable-handle",
      );

      verifyContributors("#contributors-supervisor-table", [
        "Supervisor 2",
        "Supervisor 1",
        "Supervisor 3",
      ]);

      cy.wait("@reorderSupervisors").should(
        "have.nested.property",
        "response.statusCode",
        200,
      );

      cy.verifyConflictErrorDialog("publication").closeModal("Close");

      // Verify that nothing was reordered
      cy.reload();

      verifyContributors("#contributors-supervisor-table", [
        "Supervisor 1",
        "Supervisor 2",
        "Supervisor 3",
        "Supervisor 4",
      ]);
    });
  });

  describe("Dataset creators", () => {
    beforeEach(() => {
      cy.setUpDataset();
      cy.addCreator("Creator", "1", { external: true });
      cy.addCreator("Creator", "2", { external: true });
      cy.addCreator("Creator", "3", { external: true });

      cy.visitDataset();

      cy.contains(".nav-item", "People & Affiliations").click();

      verifyContributors("#contributors-author-table", [
        "Creator 1",
        "Creator 2",
        "Creator 3",
      ]);
    });

    it("should show a conflict error when reordering dataset creators after reordering them in another session", () => {
      // Perform the reorder via AJAX
      cy.get("#contributors-author-table tbody")
        .invoke("attr", "hx-headers")
        .then(JSON.parse)
        .then((htmxValues) => {
          cy.get<string>("@biblioId").then((biblioId) => {
            // Reverse the creators
            const form = new URLSearchParams();
            form.append("position", "2");
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
        "Creator 3",
      ]);

      cy.intercept("POST", "/dataset/*/contributors/author/order").as(
        "reorderCreators",
      );

      // Perform the reorder
      cy.get("#author-1 .sortable-handle").drag("#author-2 .sortable-handle");

      verifyContributors("#contributors-author-table", [
        "Creator 1",
        "Creator 3",
        "Creator 2",
      ]);

      cy.wait("@reorderCreators").should(
        "have.nested.property",
        "response.statusCode",
        200,
      );

      cy.verifyConflictErrorDialog("dataset").closeModal("Close");

      // Verify that nothing was reordered
      cy.reload();

      verifyContributors("#contributors-author-table", [
        "Creator 3",
        "Creator 2",
        "Creator 1",
      ]);
    });

    it("should show a conflict error when reordering dataset creators after removing one in another session", () => {
      // Delete the second creator by AJAX call and cancel the confirm delete dialog
      cy.contains("tr", "Creator 2").find(".btn .if-delete").click();
      cy.ensureModal("Confirm deletion")
        .within(() => {
          cy.contains(".btn", "Delete").triggerHtmx("hx-delete");
        })
        .closeModal("Cancel");
      cy.ensureNoModal();

      // Page is not refreshed so we still think there are 3 creators
      verifyContributors("#contributors-author-table", [
        "Creator 1",
        "Creator 2",
        "Creator 3",
      ]);

      cy.intercept("POST", "/dataset/*/contributors/author/order").as(
        "reorderCreators",
      );

      // Perform the reorder
      cy.get("#author-1 .sortable-handle").drag("#author-2 .sortable-handle");

      verifyContributors("#contributors-author-table", [
        "Creator 1",
        "Creator 3",
        "Creator 2",
      ]);

      cy.wait("@reorderCreators").should(
        "have.nested.property",
        "response.statusCode",
        200,
      );

      cy.verifyConflictErrorDialog("dataset").closeModal("Close");

      // Verify that nothing was reordered
      cy.reload();

      verifyContributors("#contributors-author-table", [
        "Creator 1",
        "Creator 3",
      ]);
    });

    it("should show a conflict error when reordering dataset creators after adding one in another session", () => {
      // Add a fourth creator by AJAX call
      cy.addCreator("Creator", "4", { external: true });

      cy.intercept("POST", "/dataset/*/contributors/author/order").as(
        "reorderCreators",
      );

      // Perform the reorder
      cy.get("#author-0 .sortable-handle").drag("#author-1 .sortable-handle");

      verifyContributors("#contributors-author-table", [
        "Creator 2",
        "Creator 1",
        "Creator 3",
      ]);

      cy.wait("@reorderCreators").should(
        "have.nested.property",
        "response.statusCode",
        200,
      );

      cy.verifyConflictErrorDialog("dataset").closeModal("Close");

      // Verify that nothing was reordered
      cy.reload();

      verifyContributors("#contributors-author-table", [
        "Creator 1",
        "Creator 2",
        "Creator 3",
        "Creator 4",
      ]);
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
