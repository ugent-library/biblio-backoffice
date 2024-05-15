// https://github.com/ugent-library/biblio-backoffice/issues/1241

import * as dayjs from "dayjs";

describe("Issue #1241: As a researcher or librarian, I should not be able to select an embargo date in the past - embargo switch does not happen because of that.", () => {
  const tomorrow = dayjs().add(1, "day").format("YYYY-MM-DD");
  const today = dayjs().format("YYYY-MM-DD");

  describe("For publications", () => {
    it("should not be possible to select a date before tomorrow for embargo end date", () => {
      cy.loginAsResearcher();

      cy.setUpPublication("Book");

      cy.visitPublication();

      cy.contains(".nav-tabs .nav-item", "Full text & Files").click();

      cy.get("input[type=file][name=file]").selectFile(
        "cypress/fixtures/empty-pdf.pdf",
      );

      cy.ensureModal("Document details for file empty-pdf.pdf")
        .within(() => {
          cy.intercept("/publication/*/files/*/refresh-form*").as(
            "refreshForm",
          );

          cy.getLabel(/Embargoed access/).click();

          cy.wait("@refreshForm");

          cy.setFieldByLabel(
            "Access level during embargo",
            "UGent access - Local access",
          );
          cy.setFieldByLabel(
            "Access level after embargo",
            "Public access - Open access",
          );

          cy.get("#embargo_date")
            .should("have.attr", "min", tomorrow)
            .type(today);
        })
        .closeModal(true); // Save

      // Make sure modal didn't close
      cy.ensureModal("Document details for file empty-pdf.pdf")
        .within(() => {
          cy.get(".alert.alert-danger")
            .should("be.visible")
            .should(
              "contain.text",
              "The embargo end date should be in the future",
            );

          cy.get("#embargo_date")
            .scrollIntoView()
            .should("have.class", "is-invalid")
            .find(".invalid-feedback")
            .should("be.visible")
            .should(
              "have.text",
              "The embargo end date should be in the future",
            );

          cy.get("#embargo_date").type(tomorrow);
        })
        .closeModal(true); // Save

      cy.ensureNoModal();

      cy.get(".list-group-item")
        .should("have.length", 1)
        .find(".bc-toolbar")
        .should("contain.text", `Public access - Open access from ${tomorrow}`);
    });
  });

  describe("For datasets", () => {
    it("should not be possible to select a date before tomorrow for embargo end date", () => {
      cy.loginAsResearcher();

      cy.setUpDataset();

      cy.visitDataset();

      cy.contains(".card-header", "Dataset details")
        .contains(".btn", "Edit")
        .click();

      cy.ensureModal("Edit dataset details")
        .within(() => {
          cy.intercept("PUT", "/dataset/*/details/edit/refresh-form*").as(
            "refreshForm",
          );

          cy.setFieldByLabel("Access level", "Embargoed access");
          cy.wait("@refreshForm");

          cy.setFieldByLabel("Access level after embargo", "Open access");

          cy.get("#embargo_date")
            .should("have.attr", "min", tomorrow)
            .type(today);
        })
        .closeModal(true); // Save

      // Make sure modal didn't close
      cy.ensureModal("Edit dataset details")
        .within(() => {
          cy.get(".alert.alert-danger")
            .should("be.visible")
            .should("contain.text", "Embargo end date should be in the future");

          cy.get("#embargo_date")
            .should("have.class", "is-invalid")
            .scrollIntoView()
            .parent()
            .find(".invalid-feedback")
            .should("be.visible")
            .should("have.text", "Embargo end date should be in the future");

          cy.get("#embargo_date").type(tomorrow);
        })
        .closeModal(true); // Save

      cy.ensureNoModal();
    });
  });
});
