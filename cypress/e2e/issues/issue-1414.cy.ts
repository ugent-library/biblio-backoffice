// https://github.com/ugent-library/biblio-backoffice/issues/1414

describe("Issue #1414: JS error when closing toast", () => {
  it("should not error when you close a toast manually before auto dismissal (gohtml)", () => {
    cy.loginAsLibrarian();

    cy.setUpPublication();
    cy.visitPublication();

    cy.contains(".btn", "Lock record").click();

    cy.wait(3000);

    cy.ensureToast("Publication was successfully locked.").closeToast();

    // The error occurred after 5000ms, so we wait another 3000ms to make sure the test hasn't succeeded by that time
    cy.wait(3000);
  });

  it("should not error when you close a toast manually before auto dismissal (templ)", () => {
    cy.loginAsLibrarian();

    cy.setUpPublication();
    cy.visitPublication();

    cy.contains(".btn", "Lock record").then((button) => {
      const csrfToken = button
        .get()[0]
        .ownerDocument.querySelector("meta[name='csrf-token']")
        .getAttribute("content");

      // Send the request via AJAX so the flash message cannot be shown immediately.
      // It will be shown on the next page visit (which is a templ layout page).
      cy.request({
        method: "POST",
        url: button.attr("hx-post"),
        headers: {
          "X-CSRF-Token": csrfToken,
        },
      });
    });

    cy.visit("/dashboard");

    cy.wait(3000);

    cy.ensureToast("Publication was successfully locked.").closeToast();

    // The error occurred after 5000ms, so we wait another 3000ms to make sure the test hasn't succeeded by that time
    cy.wait(3000);
  });
});
