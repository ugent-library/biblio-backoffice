import { testFocusForLabel } from "support/util";

describe("Editing publication description", () => {
  beforeEach(() => {
    cy.loginAsResearcher();
  });

  describe("Publication details", () => {
    it("should be possible to change the publication type", () => {
      cy.setUpPublication("Dissertation");
      cy.visitPublication();

      cy.contains(".card", "Publication details")
        .contains(".btn", "Edit")
        .click();

      cy.intercept("/publication/*/type/confirm?type=journal_article").as(
        "changeType",
      );

      cy.ensureModal("Edit publication details").within(() => {
        cy.getLabel("Publication type")
          .next()
          .find("select > option:selected")
          .should("have.value", "dissertation")
          .should("have.text", "Dissertation");

        cy.setFieldByLabel("Publication type", "Journal article");
      });

      cy.wait("@changeType");

      cy.ensureModal("Changing the publication type might result in data loss")
        .within(() => {
          cy.get(".modal-body").should(
            "contain",
            "Are you sure you want to change the type to Journal article?",
          );
        })
        .closeModal("Proceed");
      cy.ensureNoModal();

      cy.contains(".card", "Publication details")
        .find(".card-body")
        .within(() => {
          cy.getLabel("Publication type")
            .next()
            .should("contain", "Journal article");
        });
    });

    it("should error when changing publication type and publication was updated concurrently", () => {
      cy.setUpPublication();
      cy.visitPublication();

      // First perform an update but also capture the snapshot ID
      cy.updateFields(
        "Abstract",
        () => {
          cy.setFieldByLabel("Abstract", "This is an abstract");
          cy.setFieldByLabel("Language", "Danish");

          cy.contains(".modal-footer .btn", "Add abstract")
            .attr("hx-headers")
            .as("initialSnapshot");
        },
        "Add abstract",
      );

      cy.contains(".card", "Publication details")
        .contains(".btn", "Edit")
        .click();

      cy.intercept("/publication/*/type/confirm?type=journal_article").as(
        "changeType",
      );

      cy.ensureModal("Edit publication details").within(() => {
        cy.setFieldByLabel("Publication type", "Journal article");
      });
      cy.wait("@changeType");

      cy.ensureModal("Changing the publication type might result in data loss")
        .within(() => {
          cy.get(".modal-body").should(
            "contain",
            "Are you sure you want to change the type to Journal article?",
          );

          // "Fix" the proceed button with the old (outdated) snapshot ID
          cy.contains(".modal-footer .btn", "Proceed").then((button) => {
            cy.get<string>("@initialSnapshot").then((initialSnapshot) => {
              button.attr("hx-headers", initialSnapshot);
            });
          });
        })
        .closeModal("Proceed");

      cy.ensureModal(null).within(() => {
        cy.contains(
          "Publication has been modified by another user. Please reload the page.",
        ).should("be.visible");
      });
    });

    it("should have clickable labels in the Journal Article form", () => {
      cy.setUpPublication("Journal Article");
      cy.visitPublication();

      cy.updateFields("Publication details", () => {
        testFocusForLabel("Publication type", 'select[name="type"]');
        testFocusForLabel(
          "Article type",
          'select[name="journal_article_type"]',
        );
        testFocusForLabel("DOI", 'input[type=text][name="doi"]');

        testFocusForLabel("Title", 'input[type=text][name="title"]');
        testFocusForLabel(
          "Alternative title",
          'input[type=text][name="alternative_title"]',
        );
        testFocusForLabel(
          "Journal title",
          'input[type=text][name="publication"]',
        );
        testFocusForLabel(
          "Short journal title",
          'input[type=text][name="publication_abbreviation"]',
        );

        testFocusForLabel("Languages", 'select[name="language"]');
        testFocusForLabel(
          "Publishing status",
          'select[name="publication_status"]',
        );
        testFocusForLabel(
          "Published while none of the authors and editors were employed at UGent",
          ':checkbox[name="extern"]',
        );
        testFocusForLabel("Publication year", 'input[type=text][name="year"]');
        testFocusForLabel(
          "Place of publication",
          'input[type=text][name="place_of_publication"]',
        );
        testFocusForLabel("Publisher", 'input[type=text][name="publisher"]');

        testFocusForLabel("Volume", 'input[type=text][name="volume"]');
        testFocusForLabel("Issue", 'input[type=text][name="issue"]');
        testFocusForLabel(
          "Special issue title",
          'input[type=text][name="issue_title"]',
        );
        testFocusForLabel("First page", 'input[type=text][name="page_first"]');
        testFocusForLabel("Last page", 'input[type=text][name="page_last"]');
        testFocusForLabel(
          "Number of pages",
          'input[type=text][name="page_count"]',
        );
        testFocusForLabel(
          "Article number",
          'input[type=text][name="article_number"]',
        );

        testFocusForLabel(
          "Web of Science ID",
          'input[type=text][name="wos_id"]',
        );
        testFocusForLabel("ISSN", 'input[type=text][name="issn"]');
        testFocusForLabel("E-ISSN", 'input[type=text][name="eissn"]');
        testFocusForLabel("ISBN", 'input[type=text][name="isbn"]');
        testFocusForLabel("E-ISBN", 'input[type=text][name="eisbn"]');
        testFocusForLabel("PubMed ID", 'input[type=text][name="pubmed_id"]');
        testFocusForLabel("Arxiv ID", 'input[type=text][name="arxiv_id"]');
        testFocusForLabel("ESCI ID", 'input[type=text][name="esci_id"]');
      });
    });

    it("should have clickable labels in the Book Chapter form", () => {
      cy.setUpPublication("Book Chapter");

      cy.visitPublication();

      cy.updateFields("Publication details", () => {
        testFocusForLabel("Publication type", 'select[name="type"]');
        testFocusForLabel("DOI", 'input[type=text][name="doi"]');

        testFocusForLabel("Title", 'input[type=text][name="title"]');
        testFocusForLabel(
          "Alternative title",
          'input[type=text][name="alternative_title"]',
        );
        testFocusForLabel("Book title", 'input[type=text][name="publication"]');

        testFocusForLabel("Languages", 'select[name="language"]');
        testFocusForLabel(
          "Publishing status",
          'select[name="publication_status"]',
        );
        testFocusForLabel(
          "Published while none of the authors and editors were employed at UGent",
          ':checkbox[name="extern"]',
        );
        testFocusForLabel("Publication year", 'input[type=text][name="year"]');
        testFocusForLabel(
          "Place of publication",
          'input[type=text][name="place_of_publication"]',
        );
        testFocusForLabel("Publisher", 'input[type=text][name="publisher"]');

        testFocusForLabel(
          "Series title",
          'input[type=text][name="series_title"]',
        );
        testFocusForLabel("Volume", 'input[type=text][name="volume"]');
        testFocusForLabel("Edition", 'input[type=text][name="edition"]');
        testFocusForLabel("First page", 'input[type=text][name="page_first"]');
        testFocusForLabel("Last page", 'input[type=text][name="page_last"]');
        testFocusForLabel(
          "Number of pages",
          'input[type=text][name="page_count"]',
        );

        testFocusForLabel(
          "Web of Science ID",
          'input[type=text][name="wos_id"]',
        );
        testFocusForLabel("ISSN", 'input[type=text][name="issn"]');
        testFocusForLabel("E-ISSN", 'input[type=text][name="eissn"]');
        testFocusForLabel("ISBN", 'input[type=text][name="isbn"]');
        testFocusForLabel("E-ISBN", 'input[type=text][name="eisbn"]');
      });
    });

    it("should have clickable labels in the Book form", () => {
      cy.setUpPublication("Book");

      cy.visitPublication();

      cy.updateFields("Publication details", () => {
        testFocusForLabel("Publication type", 'select[name="type"]');
        testFocusForLabel("DOI", 'input[type=text][name="doi"]');

        testFocusForLabel("Title", 'input[type=text][name="title"]');
        testFocusForLabel(
          "Alternative title",
          'input[type=text][name="alternative_title"]',
        );

        testFocusForLabel("Languages", 'select[name="language"]');
        testFocusForLabel(
          "Publishing status",
          'select[name="publication_status"]',
        );
        testFocusForLabel(
          "Published while none of the authors and editors were employed at UGent",
          ':checkbox[name="extern"]',
        );
        testFocusForLabel("Publication year", 'input[type=text][name="year"]');
        testFocusForLabel(
          "Place of publication",
          'input[type=text][name="place_of_publication"]',
        );
        testFocusForLabel("Publisher", 'input[type=text][name="publisher"]');

        testFocusForLabel(
          "Series title",
          'input[type=text][name="series_title"]',
        );
        testFocusForLabel("Volume", 'input[type=text][name="volume"]');
        testFocusForLabel("Edition", 'input[type=text][name="edition"]');
        testFocusForLabel(
          "Number of pages",
          'input[type=text][name="page_count"]',
        );

        testFocusForLabel(
          "Web of Science ID",
          'input[type=text][name="wos_id"]',
        );
        testFocusForLabel("ISSN", 'input[type=text][name="issn"]');
        testFocusForLabel("E-ISSN", 'input[type=text][name="eissn"]');
        testFocusForLabel("ISBN", 'input[type=text][name="isbn"]');
        testFocusForLabel("E-ISBN", 'input[type=text][name="eisbn"]');
      });
    });

    it("should have clickable labels in the Conference contribution form", () => {
      cy.setUpPublication("Conference contribution");

      cy.visitPublication();

      cy.updateFields("Publication details", () => {
        testFocusForLabel("Publication type", 'select[name="type"]');
        testFocusForLabel("Conference type", 'select[name="conference_type"]');
        testFocusForLabel("DOI", 'input[type=text][name="doi"]');

        testFocusForLabel("Title", 'input[type=text][name="title"]');
        testFocusForLabel(
          "Alternative title",
          'input[type=text][name="alternative_title"]',
        );
        testFocusForLabel(
          "Proceedings title",
          'input[type=text][name="publication"]',
        );
        testFocusForLabel(
          "Publication short title",
          'input[type=text][name="publication_abbreviation"]',
        );

        testFocusForLabel("Languages", 'select[name="language"]');
        testFocusForLabel(
          "Publishing status",
          'select[name="publication_status"]',
        );
        testFocusForLabel(
          "Published while none of the authors and editors were employed at UGent",
          ':checkbox[name="extern"]',
        );
        testFocusForLabel("Publication year", 'input[type=text][name="year"]');
        testFocusForLabel(
          "Place of publication",
          'input[type=text][name="place_of_publication"]',
        );
        testFocusForLabel("Publisher", 'input[type=text][name="publisher"]');

        testFocusForLabel(
          "Series title",
          'input[type=text][name="series_title"]',
        );
        testFocusForLabel("Volume", 'input[type=text][name="volume"]');
        testFocusForLabel("Issue", 'input[type=text][name="issue"]');
        testFocusForLabel(
          "Special issue title",
          'input[type=text][name="issue_title"]',
        );
        testFocusForLabel("First page", 'input[type=text][name="page_first"]');
        testFocusForLabel("Last page", 'input[type=text][name="page_last"]');
        testFocusForLabel(
          "Number of pages",
          'input[type=text][name="page_count"]',
        );
        testFocusForLabel(
          "Article number",
          'input[type=text][name="article_number"]',
        );

        testFocusForLabel(
          "Web of Science ID",
          'input[type=text][name="wos_id"]',
        );
        testFocusForLabel("ISSN", 'input[type=text][name="issn"]');
        testFocusForLabel("E-ISSN", 'input[type=text][name="eissn"]');
        testFocusForLabel("ISBN", 'input[type=text][name="isbn"]');
        testFocusForLabel("E-ISBN", 'input[type=text][name="eisbn"]');
      });
    });

    it("should have clickable labels in the Dissertation form", () => {
      cy.setUpPublication("Dissertation");

      cy.visitPublication();

      cy.updateFields("Publication details", () => {
        testFocusForLabel("Publication type", 'select[name="type"]');
        testFocusForLabel("DOI", 'input[type=text][name="doi"]');

        testFocusForLabel("Title", 'input[type=text][name="title"]');
        testFocusForLabel(
          "Alternative title",
          'input[type=text][name="alternative_title"]',
        );

        testFocusForLabel("Languages", 'select[name="language"]');
        testFocusForLabel(
          "Publishing status",
          'select[name="publication_status"]',
        );
        testFocusForLabel(
          "Published while none of the authors and editors were employed at UGent",
          ':checkbox[name="extern"]',
        );
        testFocusForLabel("Publication year", 'input[type=text][name="year"]');
        testFocusForLabel(
          "Place of publication",
          'input[type=text][name="place_of_publication"]',
        );
        testFocusForLabel("Publisher", 'input[type=text][name="publisher"]');

        testFocusForLabel(
          "Series title",
          'input[type=text][name="series_title"]',
        );
        testFocusForLabel("Volume", 'input[type=text][name="volume"]');
        testFocusForLabel(
          "Number of pages",
          'input[type=text][name="page_count"]',
        );

        testFocusForLabel(
          "Date of defense",
          'input[type=text][name="defense_date"]',
        );
        testFocusForLabel(
          "Place of defense",
          'input[type=text][name="defense_place"]',
        );

        testFocusForLabel(
          "Web of Science ID",
          'input[type=text][name="wos_id"]',
        );
        testFocusForLabel("ISSN", 'input[type=text][name="issn"]');
        testFocusForLabel("E-ISSN", 'input[type=text][name="eissn"]');
        testFocusForLabel("ISBN", 'input[type=text][name="isbn"]');
        testFocusForLabel("E-ISBN", 'input[type=text][name="eisbn"]');
      });
    });

    it("should have clickable labels in the Miscellaneous form", () => {
      cy.setUpPublication("Miscellaneous");

      cy.visitPublication();

      cy.updateFields("Publication details", () => {
        testFocusForLabel("Publication type", 'select[name="type"]');
        testFocusForLabel(
          "Miscellaneous type",
          'select[name="miscellaneous_type"]',
        );
        testFocusForLabel("DOI", 'input[type=text][name="doi"]');

        testFocusForLabel("Title", 'input[type=text][name="title"]');
        testFocusForLabel(
          "Alternative title",
          'input[type=text][name="alternative_title"]',
        );
        testFocusForLabel(
          "Publication title",
          'input[type=text][name="publication"]',
        );
        testFocusForLabel(
          "Publication short title",
          'input[type=text][name="publication_abbreviation"]',
        );

        testFocusForLabel("Languages", 'select[name="language"]');
        testFocusForLabel(
          "Publishing status",
          'select[name="publication_status"]',
        );
        testFocusForLabel(
          "Published while none of the authors and editors were employed at UGent",
          ':checkbox[name="extern"]',
        );
        testFocusForLabel("Publication year", 'input[type=text][name="year"]');
        testFocusForLabel(
          "Place of publication",
          'input[type=text][name="place_of_publication"]',
        );
        testFocusForLabel("Publisher", 'input[type=text][name="publisher"]');

        testFocusForLabel(
          "Series title",
          'input[type=text][name="series_title"]',
        );
        testFocusForLabel("Volume", 'input[type=text][name="volume"]');
        testFocusForLabel("Issue", 'input[type=text][name="issue"]');
        testFocusForLabel(
          "Special issue title",
          'input[type=text][name="issue_title"]',
        );
        testFocusForLabel("Edition", 'input[type=text][name="edition"]');
        testFocusForLabel("First page", 'input[type=text][name="page_first"]');
        testFocusForLabel("Last page", 'input[type=text][name="page_last"]');
        testFocusForLabel(
          "Number of pages",
          'input[type=text][name="page_count"]',
        );
        testFocusForLabel(
          "Article number",
          'input[type=text][name="article_number"]',
        );
        testFocusForLabel(
          "Report number",
          'input[type=text][name="report_number"]',
        );

        testFocusForLabel(
          "Web of Science ID",
          'input[type=text][name="wos_id"]',
        );
        testFocusForLabel("ISSN", 'input[type=text][name="issn"]');
        testFocusForLabel("E-ISSN", 'input[type=text][name="eissn"]');
        testFocusForLabel("ISBN", 'input[type=text][name="isbn"]');
        testFocusForLabel("PubMed ID", 'input[type=text][name="pubmed_id"]');
        testFocusForLabel("Arxiv ID", 'input[type=text][name="arxiv_id"]');
        testFocusForLabel("ESCI ID", 'input[type=text][name="esci_id"]');
      });
    });

    it("should have clickable labels in the Issue editor form", () => {
      cy.setUpPublication("Issue");

      cy.visitPublication();

      cy.updateFields("Publication details", () => {
        testFocusForLabel("Publication type", 'select[name="type"]');
        testFocusForLabel("DOI", 'input[type=text][name="doi"]');

        testFocusForLabel("Title", 'input[type=text][name="title"]');
        testFocusForLabel(
          "Alternative title",
          'input[type=text][name="alternative_title"]',
        );
        testFocusForLabel(
          "Journal title",
          'input[type=text][name="publication"]',
        );
        testFocusForLabel(
          "Short journal title",
          'input[type=text][name="publication_abbreviation"]',
        );

        testFocusForLabel("Languages", 'select[name="language"]');
        testFocusForLabel(
          "Publishing status",
          'select[name="publication_status"]',
        );
        testFocusForLabel(
          "Published while none of the authors and editors were employed at UGent",
          ':checkbox[name="extern"]',
        );
        testFocusForLabel("Publication year", 'input[type=text][name="year"]');
        testFocusForLabel(
          "Place of publication",
          'input[type=text][name="place_of_publication"]',
        );
        testFocusForLabel("Publisher", 'input[type=text][name="publisher"]');

        testFocusForLabel("Volume", 'input[type=text][name="volume"]');
        testFocusForLabel("Issue", 'input[type=text][name="issue"]');
        testFocusForLabel(
          "Special issue title",
          'input[type=text][name="issue_title"]',
        );
        testFocusForLabel("Edition", 'input[type=text][name="edition"]');
        testFocusForLabel(
          "Number of pages",
          'input[type=text][name="page_count"]',
        );

        testFocusForLabel(
          "Web of Science ID",
          'input[type=text][name="wos_id"]',
        );
        testFocusForLabel("ISSN", 'input[type=text][name="issn"]');
        testFocusForLabel("E-ISSN", 'input[type=text][name="eissn"]');
        testFocusForLabel("ISBN", 'input[type=text][name="isbn"]');
      });
    });
  });

  describe("Projects", () => {
    beforeEach(() => {
      cy.setUpPublication();
      cy.visitPublication();
    });

    it("should be possible to add and delete projects", () => {
      cy.get("#projects-body").should("contain", "No projects");

      cy.contains(".card", "Project").contains(".btn", "Add project").click();

      cy.ensureModal("Select projects").within(() => {
        cy.intercept("/publication/*/projects/suggestions?*").as(
          "suggestProject",
        );

        cy.getLabel("Search project").next("input").type("001D07903");
        cy.wait("@suggestProject");

        cy.contains(".list-group-item", "001D07903")
          .contains(".btn", "Add project")
          .click();
      });
      cy.ensureNoModal();

      cy.get("#projects-body")
        .contains(".list-group-item", "001D07903")
        .find(".if-more")
        .click();
      cy.contains(".dropdown-item", "Remove from publication").click();

      cy.ensureModal("Confirm deletion")
        .within(() => {
          cy.get(".modal-body").should(
            "contain",
            "Are you sure you want to remove this project from the publication?",
          );
        })
        .closeModal("Delete");
      cy.ensureNoModal();

      cy.get("#projects-body").should("contain", "No projects");
    });
  });

  describe("Abstracts", () => {
    beforeEach(() => {
      cy.setUpPublication();
      cy.visitPublication();
    });

    it("should be possible to add, edit and delete abstracts", () => {
      cy.get("#abstracts-body").should("contain", "No abstracts");

      cy.contains(".btn", "Add abstract").click();
      cy.ensureModal("Add abstract")
        .within(() => {
          cy.setFieldByLabel("Abstract", "");
          cy.setFieldByLabel("Language", "Danish");
        })
        .closeModal("Add abstract");

      cy.ensureModal("Add abstract")
        .within(() => {
          cy.contains(".alert-danger", "Abstract text can't be empty").should(
            "be.visible",
          );
          cy.get("textarea[name=text]")
            .should("have.class", "is-invalid")
            .next(".invalid-feedback")
            .should("have.text", "Abstract text can't be empty");

          cy.setFieldByLabel("Abstract", "The initial abstract");
        })
        .closeModal("Add abstract");
      cy.ensureNoModal();

      cy.get("#abstracts-body")
        .find("table tbody tr")
        .as("row")
        .should("have.length", 1);

      cy.get("@row").should("contain", "The initial abstract");
      cy.get("@row").should("contain", "Danish");

      cy.get("@row").within(() => {
        cy.get(".if-more").click();

        cy.contains("button", "Edit").click();
      });

      cy.ensureModal("Edit abstract")
        .within(() => {
          cy.setFieldByLabel("Abstract", "");
          cy.setFieldByLabel("Language", "Northern Sami");
        })
        .closeModal("Update abstract");

      cy.ensureModal("Edit abstract")
        .within(() => {
          cy.contains(".alert-danger", "Abstract text can't be empty").should(
            "be.visible",
          );
          cy.get("textarea[name=text]")
            .should("have.class", "is-invalid")
            .next(".invalid-feedback")
            .should("have.text", "Abstract text can't be empty");

          cy.setFieldByLabel("Abstract", "The updated abstract");
        })
        .closeModal("Update abstract");
      cy.ensureNoModal();

      cy.get("@row").should("have.length", 1);

      cy.get("@row").should("not.contain", "The initial abstract");
      cy.get("@row").should("contain", "The updated abstract");
      cy.get("@row").should("not.contain", "Danish");
      cy.get("@row").should("contain", "Northern Sami");

      cy.get("@row").within(() => {
        cy.get(".if-more").click();

        cy.contains("button", "Delete").click();
      });

      cy.ensureModal("Confirm deletion")
        .within(() => {
          cy.get(".modal-body").should(
            "contain",
            "Are you sure you want to remove this abstract?",
          );
        })
        .closeModal("Delete");
      cy.ensureNoModal();

      cy.get("#abstracts-body").should("contain", "No abstracts");
    });

    it("should have clickable labels in the Abstract dialog", () => {
      cy.updateFields("Abstract", () => {
        testFocusForLabel("Abstract", 'textarea[name="text"]');
        testFocusForLabel("Language", 'select[name="lang"]');
      });
    });
  });

  describe("Links", () => {
    beforeEach(() => {
      cy.setUpPublication();
      cy.visitPublication();
    });

    it("should be possible to add, edit and delete links", () => {
      cy.get("#links-body").should("contain", "No links");

      cy.contains(".btn", "Add link").click();
      cy.ensureModal("Add link")
        .within(() => {
          cy.setFieldByLabel("URL", "https://www.ugent.be");
          cy.setFieldByLabel("Relation", "Related information");
          cy.setFieldByLabel("Description", "The initial website");
        })
        .closeModal("Add link");
      cy.ensureNoModal();

      cy.get("#links-body")
        .find("table tbody tr")
        .as("row")
        .should("have.length", 1);

      cy.get("@row").should("contain", "https://www.ugent.be");
      cy.get("@row").should("contain", "Related information");
      cy.get("@row").should("contain", "The initial website");

      cy.get("@row").within(() => {
        cy.get(".if-more").click();

        cy.contains("button", "Edit").click();
      });

      cy.ensureModal("Edit link")
        .within(() => {
          cy.setFieldByLabel("URL", "https://lib.ugent.be");
          cy.setFieldByLabel("Relation", "Accompanying website");
          cy.setFieldByLabel("Description", "The updated website");
        })
        .closeModal("Update link");
      cy.ensureNoModal();

      cy.get("@row").should("have.length", 1);

      cy.get("@row").should("not.contain", "https://www.ugent.be");
      cy.get("@row").should("contain", "https://lib.ugent.be");
      cy.get("@row").should("not.contain", "Related information");
      cy.get("@row").should("contain", "Accompanying website");
      cy.get("@row").should("not.contain", "The initial website");
      cy.get("@row").should("contain", "The updated website");

      cy.get("@row").within(() => {
        cy.get(".if-more").click();

        cy.contains("button", "Delete").click();
      });

      cy.ensureModal("Confirm deletion")
        .within(() => {
          cy.get(".modal-body").should(
            "contain",
            "Are you sure you want to remove this link?",
          );
        })
        .closeModal("Delete");
      cy.ensureNoModal();

      cy.get("#links-body").should("contain", "No links");
    });

    it("should have clickable labels in the Link dialog", () => {
      cy.updateFields("Link", () => {
        testFocusForLabel("URL", 'input[type=text][name="url"]');
        testFocusForLabel("Relation", 'select[name="relation"]');
        testFocusForLabel(
          "Description",
          'input[type=text][name="description"]',
        );
      });
    });
  });

  describe("Lay summaries", () => {
    beforeEach(() => {
      cy.setUpPublication("Dissertation");
      cy.visitPublication();
    });

    it("should be possible to add, edit and delete lay summaries", () => {
      cy.get("#lay-summaries-body").should("contain", "No lay summaries");

      cy.contains(".btn", "Add lay summary").click();
      cy.ensureModal("Add lay summary")
        .within(() => {
          cy.setFieldByLabel("Lay summary", "");
          cy.setFieldByLabel("Language", "Italian");
        })
        .closeModal("Add lay summary");

      cy.ensureModal("Add lay summary")
        .within(() => {
          cy.contains(
            ".alert-danger",
            "Lay summary text can't be empty",
          ).should("be.visible");
          cy.get("textarea[name=text]")
            .should("have.class", "is-invalid")
            .next(".invalid-feedback")
            .should("have.text", "Lay summary text can't be empty");

          cy.setFieldByLabel("Lay summary", "The initial lay summary");
        })
        .closeModal("Add lay summary");
      cy.ensureNoModal();

      cy.get("#lay-summaries-body")
        .find("table tbody tr")
        .as("row")
        .should("have.length", 1);

      cy.get("@row").should("contain", "The initial lay summary");
      cy.get("@row").should("contain", "Italian");

      cy.get("@row").within(() => {
        cy.get(".if-more").click();

        cy.contains("button", "Edit").click();
      });

      cy.ensureModal("Edit lay summary")
        .within(() => {
          cy.setFieldByLabel("Lay summary", "");
          cy.setFieldByLabel("Language", "Multiple languages");
        })
        .closeModal("Update lay summary");

      cy.ensureModal("Edit lay summary")
        .within(() => {
          cy.contains(
            ".alert-danger",
            "Lay summary text can't be empty",
          ).should("be.visible");
          cy.get("textarea[name=text]")
            .should("have.class", "is-invalid")
            .next(".invalid-feedback")
            .should("have.text", "Lay summary text can't be empty");

          cy.setFieldByLabel("Lay summary", "The updated lay summary");
        })
        .closeModal("Update lay summary");
      cy.ensureNoModal();

      cy.get("@row").should("have.length", 1);

      cy.get("@row").should("not.contain", "The initial lay summary");
      cy.get("@row").should("contain", "The updated lay summary");
      cy.get("@row").should("not.contain", "Italian");
      cy.get("@row").should("contain", "Multiple languages");

      cy.get("@row").within(() => {
        cy.get(".if-more").click();

        cy.contains("button", "Delete").click();
      });

      cy.ensureModal("Confirm deletion")
        .within(() => {
          cy.get(".modal-body").should(
            "contain",
            "Are you sure you want to remove this lay summary?",
          );
        })
        .closeModal("Delete");
      cy.ensureNoModal();

      cy.get("#lay-summaries-body").should("contain", "No lay summaries");
    });

    it("should have clickable labels in the Lay summary dialog", () => {
      cy.updateFields("Lay summary", () => {
        testFocusForLabel("Lay summary", 'textarea[name="text"]');
        testFocusForLabel("Language", 'select[name="lang"]');
      });
    });
  });

  describe("Conference details", () => {
    beforeEach(() => {
      cy.setUpPublication("Conference contribution");
      cy.visitPublication();
    });

    it("should be possible to add and edit conference details", () => {
      cy.get("#conference-details").contains(".btn", "Edit").click();
      cy.ensureModal("Edit conference details")
        .within(() => {
          cy.setFieldByLabel("Conference", "The conference name");
          cy.setFieldByLabel("Conference location", "The conference location");
          cy.setFieldByLabel(
            "Conference organiser",
            "The conference organiser",
          );
          cy.setFieldByLabel("Conference start date", "2021-01-01");
          cy.setFieldByLabel("Conference end date", "2022-02-02");
        })
        .closeModal(true);
      cy.ensureNoModal();

      cy.get("#conference-details")
        .should("contain", "The conference name")
        .should("contain", "The conference location")
        .should("contain", "The conference organiser")
        .should("contain", "2021-01-01")
        .should("contain", "2022-02-02");

      cy.get("#conference-details").contains(".btn", "Edit").click();
      cy.ensureModal("Edit conference details")
        .within(() => {
          cy.setFieldByLabel("Conference", "The updated conference name");
          cy.setFieldByLabel(
            "Conference location",
            "The updated conference location",
          );
          cy.setFieldByLabel(
            "Conference organiser",
            "The updated conference organiser",
          );
          cy.setFieldByLabel("Conference start date", "2023-03-03");
          cy.setFieldByLabel("Conference end date", "2024-04-04");
        })
        .closeModal(true);
      cy.ensureNoModal();

      cy.get("#conference-details")
        .should("contain", "The updated conference name")
        .should("contain", "The updated conference location")
        .should("contain", "The updated conference organiser")
        .should("contain", "2023-03-03")
        .should("contain", "2024-04-04");
    });

    it("should have clickable labels in the Conference details dialog", () => {
      cy.updateFields("Conference details", () => {
        testFocusForLabel("Conference", 'input[type=text][name="name"]');
        testFocusForLabel(
          "Conference location",
          'input[type=text][name="location"]',
        );
        testFocusForLabel(
          "Conference organiser",
          'input[type=text][name="organizer"]',
        );
        testFocusForLabel(
          "Conference start date",
          'input[type=text][name="start_date"]',
        );
        testFocusForLabel(
          "Conference end date",
          'input[type=text][name="end_date"]',
        );
      });
    });
  });

  describe("Additional info", () => {
    beforeEach(() => {
      cy.setUpPublication();
      cy.visitPublication();
    });

    it("should be possible to add and edit additional info", () => {
      cy.get("#additional-information").contains(".btn", "Edit").click();
      cy.ensureModal("Edit additional information")
        .within(() => {
          cy.setFieldByLabel("Research field", "Performing Arts")
            .next("button.form-value-add")
            .click()
            .closest(".form-value")
            .next(".form-value")
            .find("select")
            .select("Social Sciences");

          cy.setFieldByLabel(
            "Keywords",
            "these{enter}are{enter}the{enter}keywords",
          );

          cy.setFieldByLabel(
            "Additional information",
            "The additional information",
          );
        })
        .closeModal(true);
      cy.ensureNoModal();

      cy.get("#additional-information")
        .should("contain", "Performing Arts")
        .should("contain", "Social Sciences")
        .should("contain", "The additional information")
        .find(".badge")
        .should("have.length", 4)
        .map("textContent")
        .should("have.ordered.members", ["these", "are", "the", "keywords"]);

      cy.get("#additional-information").contains(".btn", "Edit").click();
      cy.ensureModal("Edit additional information")
        .within(() => {
          cy.getLabel("Research field")
            .next()
            .find(".form-values")
            .as("formValues")
            .contains("select", "Performing Arts")
            .next("button:contains(Delete)")
            .click();

          cy.get("@formValues")
            .find(".form-value")
            .last()
            .find("select")
            .select("Chemistry");

          cy.get("tags").contains("tag", "these").find("x").click();
          cy.get("tags").contains("tag", "are").find("x").click();
          cy.get("tags span[contenteditable]").type("updated");

          cy.setFieldByLabel(
            "Additional information",
            "The updated information",
          );
        })
        .closeModal(true);
      cy.ensureNoModal();

      cy.get("#additional-information")
        .should("contain", "Social Sciences")
        .should("contain", "Chemistry")
        .should("not.contain", "Performing Arts")
        .should("contain", "The updated information")
        .should("not.contain", "The additional information")
        .find(".badge")
        .should("have.length", 3)
        .map("textContent")
        .should("have.ordered.members", ["the", "keywords", "updated"]);
    });

    it("should have clickable labels in the Additional info dialog", () => {
      cy.updateFields("Additional information", () => {
        testFocusForLabel("Research field", 'select[name="research_field"]');
        testFocusForLabel(
          "Keywords",
          ".tags:has(textarea#keyword) tags span.tagify__input[contenteditable]",
        );
        testFocusForLabel(
          "Additional information",
          'textarea[name="additional_info"]',
        );
      });
    });
  });
});
