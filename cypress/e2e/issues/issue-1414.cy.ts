// https://github.com/ugent-library/biblio-backoffice/issues/1414

describe("Issue #1414: JS error when closing toast", () => {
  it("should not error when you close a toast manually before auto dismissal", () => {
    cy.loginAsLibrarian();

    cy.setUpPublication();
    cy.visitPublication();

    cy.contains(".btn", "Lock record").click();

    cy.ensureToast("Publication was successfully locked.");

    cy.wait(3000);

    cy.ensureToast("Publication was successfully locked.").closeToast();

    // The error occurred after 5000ms, so we wait another 3000ms to make sure the test hasn't succeeded by that time
    cy.wait(3000);
  });
});
