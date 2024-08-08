import { testFormAccessibility } from "support/util";

describe("Batch publication update", () => {
  beforeEach(() => {
    cy.loginAsLibrarian("librarian1");
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
