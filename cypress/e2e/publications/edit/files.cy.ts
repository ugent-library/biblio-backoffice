import * as dayjs from "dayjs";
import { testFormAccessibility } from "support/util";

describe("Editing publication files", () => {
  beforeEach(() => {
    cy.loginAsResearcher("researcher1");

    cy.setUpPublication();
    cy.visitPublication();

    cy.contains(".nav-link", "Full text & Files").click();
  });

  it("should be possible to add, edit and delete files", () => {
    cy.get("#files-body").should("contain", "No files");

    cy.get("input[name=file]").selectFile("cypress/fixtures/empty-pdf.pdf");
    cy.ensureModal("Document details for file empty-pdf.pdf")
      .within(() => {
        cy.intercept("/publication/*/files/*/refresh*").as("refreshForm");

        cy.setFieldByLabel("Document type", "Peer review report");
        cy.wait("@refreshForm");
        cy.contains("label", "Embargoed access").click();
        cy.wait("@refreshForm");
        cy.setFieldByLabel(
          "License granted by the rights holder",
          "No license (in copyright)",
        );
      })
      .closeModal("Save");

    cy.ensureModal("Document details for file empty-pdf.pdf")
      .within(() => {
        cy.get(".alert-danger ul li")
          .map("textContent")
          .should("have.members", [
            "The embargo end date is not a valid date",
            "The selected access level is not a valid access level value",
            "The selected access level is not a valid access level value",
          ]);
        cy.get("select[name=access_level_during_embargo]")
          .should("have.class", "is-invalid")
          .next(".invalid-feedback")
          .should(
            "have.text",
            "The selected access level is not a valid access level value",
          );
        cy.get("select[name=access_level_after_embargo]")
          .should("have.class", "is-invalid")
          .nextAll(".invalid-feedback")
          .should(
            "have.text",
            "The selected access level is not a valid access level value",
          );
        cy.get("input[name=embargo_date]")
          .should("have.class", "is-invalid")
          .nextAll(".invalid-feedback")
          .should("have.text", "The embargo end date is not a valid date");

        cy.setFieldByLabel(
          "Access level during embargo",
          "Private access - Closed access",
        );
        cy.setFieldByLabel(
          "Access level after embargo",
          "Public access - Open access",
        );
        cy.setFieldByLabel(
          "Embargo end",
          dayjs().add(1, "day").format("YYYY-MM-DD"),
        );
      })
      .closeModal("Save");
    cy.ensureNoModal();

    cy.get("#files-body .list-group .list-group-item")
      .as("row")
      .should("have.length", 1);

    cy.get("@row").should("contain", "empty-pdf.pdf");
    cy.get("@row").should("contain", "Peer review report");
    cy.get("@row").should("contain", "Embargoed access");
    cy.get("@row").should("contain", "Private access - Closed access");
    cy.get("@row").should(
      "contain",
      `Public access - Open access from ${dayjs()
        .add(1, "day")
        .format("YYYY-MM-DD")}`,
    );

    cy.get("@row").find(".if-edit").click();

    cy.ensureModal("Document details for file empty-pdf.pdf")
      .within(() => {
        cy.setFieldByLabel("Document type", "Full text");
        cy.wait("@refreshForm");
        cy.setFieldByLabel("Publication version", "Author's original (AO)");
        cy.contains("label", "UGent access - Local access").click();
        cy.wait("@refreshForm");
        cy.setFieldByLabel(
          "License granted by the rights holder",
          "CC BY (4.0)",
        );
      })
      .closeModal("Save");

    cy.get("@row").should("have.length", 1);

    cy.get("@row")
      .should("not.contain", "Peer review report")
      .should("contain", "Full text")
      .should("not.contain", "Private access - Closed access")
      .should("not.contain", "Public access - Open access")
      .should("contain", "UGent access - Local access");

    cy.get("@row").find(".if-delete").click();

    cy.ensureModal("Confirm deletion")
      .within(() => {
        cy.get(".modal-body").should(
          "contain",
          "Are you sure you want to remove empty-pdf.pdf from the publication?",
        );
      })
      .closeModal("Delete");
    cy.ensureNoModal();

    cy.get("#files-body").should("contain", "No files");
  });

  it("should display the upload date and embargo date in the correct format", () => {
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

  it("should have clickable labels in the add/edit file upload form", () => {
    cy.get("input[type=file][name=file]").selectFile(
      "cypress/fixtures/empty-pdf.pdf",
    );

    cy.ensureModal("Document details for file empty-pdf.pdf").within(() => {
      // First stub the refresh so nothing changes when clicking the access level radio button labels
      cy.intercept("/publication/*/files/*/refresh-form*", "");

      testFormAccessibility(
        {
          "select[name=relation]": "Document type",
          "select[name=publication_version]": "Publication version",

          "input[type=radio][name=access_level][value='info:eu-repo/semantics/openAccess']":
            /Public access - Open access/,
          "input[type=radio][name=access_level][value='info:eu-repo/semantics/restrictedAccess']":
            /UGent access - Local access/,
          "input[type=radio][name=access_level][value='info:eu-repo/semantics/embargoedAccess']":
            /Embargoed access/,
          "input[type=radio][name=access_level][value='info:eu-repo/semantics/closedAccess']":
            /Private access - Closed access/,

          "select[name=license]": "License granted by the rights holder",
        },
        "Document type",
      );

      // Now activate the access level radio buttons again, so we can load the embargoed access specific fields
      cy.intercept("/publication/*/files/*/refresh-form*", (req) =>
        req.continue(),
      );

      // This enables the embargo fields
      cy.getLabel(/Embargoed access/).click();

      cy.setFieldByLabel(
        "Access level during embargo",
        "Private access - Closed access",
      );
      cy.setFieldByLabel(
        "Access level after embargo",
        "Public access - Open access",
      );
      cy.setFieldByLabel(
        "Embargo end",
        dayjs().add(5, "days").format("YYYY-MM-DD"),
      );

      cy.contains(".btn", "Save").click();
    });

    cy.ensureNoModal();

    cy.contains(".list-group-item", "empty-pdf.pdf").find(".if-edit").click();

    cy.ensureModal("Document details for file empty-pdf.pdf").within(() => {
      // Stub the refresh again so nothing changes when clicking the access level radio button labels
      cy.intercept("/publication/*/files/*/refresh-form*", "");

      // Test form again with embargo fields
      testFormAccessibility(
        {
          "select[name=relation]": "Document type",
          "select[name=publication_version]": "Publication version",

          "input[type=radio][name=access_level][value='info:eu-repo/semantics/openAccess']":
            /Public access - Open access/,
          "input[type=radio][name=access_level][value='info:eu-repo/semantics/restrictedAccess']":
            /UGent access - Local access/,
          "input[type=radio][name=access_level][value='info:eu-repo/semantics/embargoedAccess']":
            /Embargoed access/,
          "input[type=radio][name=access_level][value='info:eu-repo/semantics/closedAccess']":
            /Private access - Closed access/,

          "select[name=access_level_during_embargo]":
            "Access level during embargo",
          "select[name=access_level_after_embargo]":
            "Access level after embargo",
          "input[type=date][name=embargo_date]": "Embargo end",

          "select[name=license]": "License granted by the rights holder",
        },
        "Document type",
      );
    });
  });

  it("should not set autofocus when popup is refreshed", () => {
    cy.get("input[type=file][name=file]").selectFile(
      "cypress/fixtures/empty-pdf.pdf",
    );

    cy.ensureModal("Document details for file empty-pdf.pdf").within(() => {
      cy.focused().should("have.attr", "id", "relation");

      cy.intercept("/publication/*/files/*/refresh-form*", (req) => {
        req.on("response", (res) => {
          // Pre-check of assertion so command log doesn't get bloated with massive HTML blocks
          if (typeof res.body === "string" && res.body.includes("autofocus")) {
            expect(res.body).to.not.contain("autofocus");
          }
        });
      }).as("refreshForm");

      cy.contains("label", "Embargoed access").click();
      cy.wait("@refreshForm");
      cy.focused().should(
        "have.attr",
        "id",
        "access-level-info:eu-repo/semantics/embargoedAccess",
      );

      cy.setFieldByLabel("Document type", "Data fact sheet");
      cy.wait("@refreshForm");
      cy.focused().should("have.attr", "id", "relation");

      cy.contains("label", "Public access - Open access").click();
      cy.wait("@refreshForm");
      cy.focused().should(
        "have.attr",
        "id",
        "access-level-info:eu-repo/semantics/openAccess",
      );

      cy.get("@refreshForm.all").should("have.length", 3);
    });
  });
});
