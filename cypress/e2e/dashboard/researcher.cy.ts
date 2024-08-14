import * as dayjs from "dayjs";
import { getRandomText } from "support/util";

describe("The researcher dashboard", () => {
  describe("Recent activity", () => {
    it("should display recent publication activity", () => {
      const PUBLICATION_TITLE = "My publication title " + getRandomText();

      cy.loginAsLibrarian("librarian1", "Researcher");

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
          cy.setFieldByLabel("Title", PUBLICATION_TITLE);
          cy.setFieldByLabel("Publication year", "1981");
        },
        true,
      );
      cy.contains(".btn", "Publish to Biblio").click();
      cy.ensureModal("Are you sure?").closeModal("Publish");

      verifyMostRecentActivity(
        `You published a publication: ${PUBLICATION_TITLE}.`,
      );

      cy.visitPublication();
      cy.contains(".btn", "Withdraw").click();
      cy.ensureModal("Are you sure?").closeModal("Withdraw");
      verifyMostRecentActivity(
        `You withdrew a publication: ${PUBLICATION_TITLE}.`,
      );

      cy.visitPublication();
      cy.contains(".btn", "Republish").click();
      cy.ensureModal("Are you sure?").closeModal("Republish");
      verifyMostRecentActivity(
        `You republished a publication: ${PUBLICATION_TITLE}.`,
      );

      cy.visitPublication();
      cy.contains(".btn", "Lock record").click();
      verifyMostRecentActivity(
        `You locked a publication: ${PUBLICATION_TITLE}.`,
      );

      cy.visitPublication();
      cy.contains(".btn", "Unlock record").click();
      verifyMostRecentActivity(
        `You unlocked a publication: ${PUBLICATION_TITLE}.`,
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
        `You left a comment on a publication: ${PUBLICATION_TITLE}.`,
      );
    });

    it("should display recent dataset activity", () => {
      const DATASET_TITLE = "My dataset title " + getRandomText();

      cy.loginAsLibrarian("librarian1", "Researcher");

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
          cy.setFieldByLabel("Title", DATASET_TITLE);
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

      verifyMostRecentActivity(`You published a dataset: ${DATASET_TITLE}.`);

      cy.visitDataset();
      cy.contains(".btn", "Withdraw").click();
      cy.ensureModal("Are you sure?").closeModal("Withdraw");
      verifyMostRecentActivity(`You withdrew a dataset: ${DATASET_TITLE}.`);

      cy.visitDataset();
      cy.contains(".btn", "Republish").click();
      cy.ensureModal("Are you sure?").closeModal("Republish");
      verifyMostRecentActivity(`You republished a dataset: ${DATASET_TITLE}.`);

      cy.visitDataset();
      cy.contains(".btn", "Lock record").click();
      verifyMostRecentActivity(`You locked a dataset: ${DATASET_TITLE}.`);

      cy.visitDataset();
      cy.contains(".btn", "Unlock record").click();
      verifyMostRecentActivity(`You unlocked a dataset: ${DATASET_TITLE}.`);

      cy.visitDataset();
      cy.updateFields(
        "Messages from and for Biblio team",
        () => {
          cy.setFieldByLabel("Message", "The biblio message");
        },
        true,
      );
      verifyMostRecentActivity(
        `You left a comment on a dataset: ${DATASET_TITLE}.`,
      );
    });

    it("should show other actor names but not biblio team member names", () => {
      const PUBLICATION_TITLE = "The Dissertation title " + getRandomText();

      // As researcher1
      cy.loginAsResearcher("researcher1");
      cy.setUpPublication("Dissertation", {
        prepareForPublishing: true,
        title: PUBLICATION_TITLE,
      });
      cy.addAuthor("Biblio", "Researcher2");
      cy.addAuthor("Biblio", "Librarian1");
      cy.addAuthor("Biblio", "Librarian2");
      verifyMostRecentActivity(
        `You edited a publication: ${PUBLICATION_TITLE}.`,
      );

      // As researcher2
      cy.loginAsResearcher("researcher2");
      verifyMostRecentActivity(
        `Biblio Researcher1 edited a publication: ${PUBLICATION_TITLE}.`,
      );
      cy.addSupervisor("Biblio", "Researcher1");
      verifyMostRecentActivity(
        `You edited a publication: ${PUBLICATION_TITLE}.`,
      );

      // As researcher1
      cy.loginAsResearcher("researcher1");
      verifyMostRecentActivity(
        `Biblio Researcher2 edited a publication: ${PUBLICATION_TITLE}.`,
      );

      // As librarian1
      cy.loginAsLibrarian("librarian1", "Researcher");
      verifyMostRecentActivity(
        `Biblio Researcher2 edited a publication: ${PUBLICATION_TITLE}.`,
      );
      cy.addSupervisor("Biblio", "Researcher2");
      verifyMostRecentActivity(
        `You edited a publication: ${PUBLICATION_TITLE}.`,
      );

      // As researcher1
      cy.loginAsResearcher("researcher1");
      verifyMostRecentActivity(
        `A Biblio team member edited a publication: ${PUBLICATION_TITLE}.`,
      );

      // As librarian2
      cy.loginAsLibrarian("librarian2", "Researcher");
      verifyMostRecentActivity(
        `Biblio Librarian1 edited a publication: ${PUBLICATION_TITLE}.`,
      );
    });
  });

  const NO_LOG = { log: false };

  function verifyMostRecentActivity(expectedText: string) {
    cy.wait(1000, NO_LOG);

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
