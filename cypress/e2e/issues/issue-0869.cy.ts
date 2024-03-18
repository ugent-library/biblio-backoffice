// https://github.com/ugent-library/biblio-backoffice/issues/869

import * as dayjs from "dayjs";

describe("Issue #869: Embargo date input differs from displayed output", () => {
  it("should display the upload date and embargo date in the correct format", () => {
    cy.loginAsResearcher();

    cy.setUpPublication();

    cy.contains("Full text & Files").click();

    cy.get("input[type=file][name=file]").selectFile(
      "cypress/fixtures/empty-pdf.pdf",
    );

    const embargoDate = dayjs().add(5, "days");

    cy.ensureModal("Document details for file empty-pdf.pdf")
      .within(() => {
        cy.intercept("/publication/*/files/*/refresh-form*").as("refreshForm");

        cy.getLabel(/Embargoed access/).click();
        cy.wait("@refreshForm");

        cy.setFieldByLabel(
          "Access level during embargo",
          "Private access - Closed access",
        );
        cy.setFieldByLabel(
          "Access level after embargo",
          "Public access - Open access",
        );

        cy.get("#embargo_date").type(embargoDate.format("YYYY-MM-DD"));
      })
      .closeModal(true);

    cy.get("#files-body .list-group .list-group-item")
      .as("item")
      .should("have.length", 1);

    cy.get("@item")
      .find(".bc-toolbar")
      .should("contain", "Embargoed access")
      .should("contain", "Private access - Closed access")
      .should(
        "contain",
        "Public access - Open access from " + embargoDate.format("YYYY-MM-DD"),
      );

    cy.get("@item")
      .contains(".list-group-item-main .c-body-small", "Uploaded")
      .invoke("text")
      .should("match", /^Uploaded \d{4}-\d{2}-\d{2} at \d{2}:\d{2}$/);
  });
});
