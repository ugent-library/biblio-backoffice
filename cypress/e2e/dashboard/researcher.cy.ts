import * as dayjs from "dayjs";

describe("The researcher dashboard", () => {
  describe("Recent activity", () => {
    it("should display recent publication activity", () => {
      cy.loginAsLibrarian();

      cy.setUpPublication(undefined, {
        otherFields: {
          // Make sure the publication is not immediately updated after creation
          title: null,
          classification: null,
        },
      });
      verifyMostRecentActivity(
        "You started a draft publication: Untitled record.",
      );

      cy.addAuthor("John", "Doe");
      verifyMostRecentActivity("You edited a publication: Untitled record.");

      cy.visitPublication();
      cy.updateFields(
        "Publication details",
        () => {
          cy.setFieldByLabel("Title", "My publication title");
          cy.setFieldByLabel("Publication year", "1981");
        },
        true,
      );
      cy.contains(".btn", "Publish to Biblio").click();
      cy.ensureModal("Are you sure?").closeModal("Publish");

      verifyMostRecentActivity(
        "You published a publication: My publication title.",
      );

      cy.visitPublication();
      cy.contains(".btn", "Withdraw").click();
      cy.ensureModal("Are you sure?").closeModal("Withdraw");
      verifyMostRecentActivity(
        "You withdrew a publication: My publication title.",
      );

      cy.visitPublication();
      cy.contains(".btn", "Republish").click();
      cy.ensureModal("Are you sure?").closeModal("Republish");
      verifyMostRecentActivity(
        "You republished a publication: My publication title.",
      );

      cy.visitPublication();
      cy.contains(".btn", "Lock record").click();
      verifyMostRecentActivity(
        "You locked a publication: My publication title.",
      );

      cy.visitPublication();
      cy.contains(".btn", "Unlock record").click();
      verifyMostRecentActivity(
        "You unlocked a publication: My publication title.",
      );

      cy.visitPublication();
      cy.updateFields(
        "Messages from and for Biblio team",
        () => {
          cy.setFieldByLabel("Message", "The biblio message");
        },
        true,
      );
      verifyMostRecentActivity(
        "You left a comment on a publication: My publication title.",
      );
    });

    it("should display recent dataset activity", () => {
      cy.loginAsLibrarian();

      cy.setUpDataset({
        // Make sure the dataset is not immediately updated after creation
        otherFields: {
          title: null,
          identifier_type: null,
          identifier: null,
        },
      });
      verifyMostRecentActivity("You started a draft dataset: Untitled record.");

      cy.addCreator("John", "Doe");
      verifyMostRecentActivity("You edited a dataset: Untitled record.");

      cy.visitDataset();
      cy.updateFields(
        "Dataset details",
        () => {
          cy.setFieldByLabel("Title", "My dataset title");
          cy.setFieldByLabel("Persistent identifier type", "DOI");
          cy.setFieldByLabel("Identifier", "10.5072/test/t");

          cy.setFieldByLabel("Publication year", "1981");
          cy.setFieldByLabel("Publisher", "UGent");

          cy.setFieldByLabel("Data format", "text/csv")
            .next(".autocomplete-hits")
            .contains(".badge", "text/csv")
            .click();

          cy.setFieldByLabel("License", "CC BY (4.0)");
          cy.setFieldByLabel("Access level", "Open access");
        },
        true,
      );
      cy.contains(".btn", "Publish to Biblio").click();
      cy.ensureModal("Are you sure?").closeModal("Publish");

      verifyMostRecentActivity("You published a dataset: My dataset title.");

      cy.visitDataset();
      cy.contains(".btn", "Withdraw").click();
      cy.ensureModal("Are you sure?").closeModal("Withdraw");
      verifyMostRecentActivity("You withdrew a dataset: My dataset title.");

      cy.visitDataset();
      cy.contains(".btn", "Republish").click();
      cy.ensureModal("Are you sure?").closeModal("Republish");
      verifyMostRecentActivity("You republished a dataset: My dataset title.");

      cy.visitDataset();
      cy.contains(".btn", "Lock record").click();
      verifyMostRecentActivity("You locked a dataset: My dataset title.");

      cy.visitDataset();
      cy.contains(".btn", "Unlock record").click();
      verifyMostRecentActivity("You unlocked a dataset: My dataset title.");

      cy.visitDataset();
      cy.updateFields(
        "Messages from and for Biblio team",
        () => {
          cy.setFieldByLabel("Message", "The biblio message");
        },
        true,
      );
      verifyMostRecentActivity(
        "You left a comment on a dataset: My dataset title.",
      );
    });

    it("should show other actor names but not biblio team member names", () => {
      // As researcher
      cy.loginAsResearcher();
      cy.setUpPublication("Dissertation", { prepareForPublishing: true });
      cy.addAuthor("Biblio", "Librarian");
      verifyMostRecentActivity(
        "You edited a publication: The Dissertation title.",
      );

      // As librarian
      cy.loginAsLibrarian();
      cy.visit("/");
      cy.contains("button.dropdown-toggle", "Biblio Librarian").click();
      cy.contains(".dropdown-item", "View as").click();
      cy.ensureModal("View as other user").within(() => {
        cy.setFieldByLabel("First name", "John");
        cy.setFieldByLabel("Last name", "Doe");
        cy.contains(".btn", "Change user").click();
      });

      // As John Doe
      cy.contains("Viewing the perspective of John Doe.").should("be.visible");
      verifyMostRecentActivity(
        "Biblio Researcher edited a publication: The Dissertation title.",
      );

      cy.visitPublication();
      cy.updateFields(
        "Messages from and for Biblio team",
        () => {
          cy.setFieldByLabel("Message", "The biblio message");
        },
        true,
      );
      verifyMostRecentActivity(
        "You left a comment on a publication: The Dissertation title.",
      );

      // As researcher
      cy.loginAsResearcher();
      verifyMostRecentActivity(
        "John Doe left a comment on a publication: The Dissertation title.",
      );

      // As librarian
      cy.loginAsLibrarian();
      cy.visitPublication();
      cy.contains(".btn", "Lock record").click();
      verifyMostRecentActivity(
        "You locked a publication: The Dissertation title.",
      );

      // As researcher
      cy.loginAsResearcher();
      verifyMostRecentActivity(
        "A Biblio team member locked a publication: The Dissertation title.",
      );
    });
  });

  const NO_LOG = { log: false };

  function verifyMostRecentActivity(expectedText: string) {
    cy.wait(500, NO_LOG);

    cy.visit("/dashboard", NO_LOG);

    cy.get(".c-activity-item", NO_LOG)
      .first(NO_LOG)
      .as("activity")
      .find(".c-activity-item__date", NO_LOG)
      .invoke(NO_LOG, "text")
      .should("be.oneOf", [
        dayjs().subtract(1, "minute").format("YYYY-MM-DD HH:mm"),
        dayjs().format("YYYY-MM-DD HH:mm"),
        dayjs().add(1, "minute").format("YYYY-MM-DD HH:mm"),
      ]);

    cy.get("@activity", NO_LOG)
      .find(".c-activity-item__text", NO_LOG)
      .as("activityText")
      .invoke(NO_LOG, "text")
      .then((t) => t.replace(/ +/g, " "))
      .should("contain", expectedText);

    cy.get("@biblioId", NO_LOG).then((biblioId) => {
      cy.get("@activityText", NO_LOG)
        .find(".c-activity-item__link", NO_LOG)
        .invoke(NO_LOG, "attr", "href")
        .should("end.with", "/" + biblioId);
    });
  }
});
