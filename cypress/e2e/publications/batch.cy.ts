import { testFormAccessibility } from "support/util";

describe("Batch publication update", () => {
  beforeEach(() => {
    cy.loginAsLibrarian();
    cy.switchMode("Librarian");
    cy.visit("/publication/batch");
  });

  it("should have clickable labels in the form", () => {
    testFormAccessibility(
      {
        "textarea[name=mutations]": "Operations",
      },
      "Operations",
    );
  });
});
