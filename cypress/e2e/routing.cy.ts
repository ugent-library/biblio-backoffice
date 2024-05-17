describe("Authorization", () => {
  it("should not be possible to edit or delete publications from another user", () => {
    cy.loginAsLibrarian();
    cy.setUpPublication();

    cy.loginAsResearcher();

    testForbiddenPublicationRoute("/add/description");
    testForbiddenPublicationRoute("/add/confirm");
    testForbiddenPublicationRoute("/add/publish", "POST");
    testForbiddenPublicationRoute("/add/finish");

    testForbiddenPublicationRoute(""); // The regular view publication route
    testForbiddenPublicationRoute("/description");
    testForbiddenPublicationRoute("/files");
    testForbiddenPublicationRoute("/contributors");
    testForbiddenPublicationRoute("/datasets");
    testForbiddenPublicationRoute("/activity");
    testForbiddenPublicationRoute("/files/file-123");

    testForbiddenPublicationRoute("/confirm-delete");
    testForbiddenPublicationRoute("", "DELETE");
    testForbiddenPublicationRoute("/lock", "POST");
    testForbiddenPublicationRoute("/unlock", "POST");
    testForbiddenPublicationRoute("/publish/confirm");
    testForbiddenPublicationRoute("/publish", "POST");
    testForbiddenPublicationRoute("/withdraw/confirm");
    testForbiddenPublicationRoute("/withdraw", "POST");
    testForbiddenPublicationRoute("/republish/confirm");
    testForbiddenPublicationRoute("/republish", "POST");

    testForbiddenPublicationRoute("/message/edit", "GET", "PUT");
    testForbiddenPublicationRoute("/message", "PUT");

    testForbiddenPublicationRoute("/details/edit", "GET", "PUT");
    testForbiddenPublicationRoute("/type/confirm");
    testForbiddenPublicationRoute("/type", "PUT");

    testForbiddenPublicationRoute("/conference/edit");
    testForbiddenPublicationRoute("/conference", "PUT");
    testForbiddenPublicationRoute("/additional-info/edit");
    testForbiddenPublicationRoute("/additional-info", "PUT");

    testForbiddenPublicationRoute("/projects/add");
    testForbiddenPublicationRoute("/projects/suggestions");
    testForbiddenPublicationRoute("/projects", "POST");
    testForbiddenPublicationRoute(
      "/snapshot-123/projects/confirm-delete/project-123",
    );
    testForbiddenPublicationRoute("/projects/project-123", "DELETE");

    testForbiddenPublicationRoute("/links/add");
    testForbiddenPublicationRoute("/links", "POST");
    testForbiddenPublicationRoute("/links/link-123/edit");
    testForbiddenPublicationRoute("/links/link-123", "PUT", "DELETE");
    testForbiddenPublicationRoute(
      "/snapshot-123/links/link-123/confirm-delete",
    );

    testForbiddenPublicationRoute("/departments/add");
    testForbiddenPublicationRoute("/departments/suggestions");
    testForbiddenPublicationRoute("/departments", "POST");
    testForbiddenPublicationRoute(
      "/snapshot-123/departments/department-123/confirm-delete",
    );
    testForbiddenPublicationRoute("/departments/department-123", "DELETE");

    testForbiddenPublicationRoute("/abstracts/add");
    testForbiddenPublicationRoute("/abstracts", "POST");
    testForbiddenPublicationRoute("/abstracts/abstract-123/edit");
    testForbiddenPublicationRoute("/abstracts/abstract-123", "PUT", "DELETE");
    testForbiddenPublicationRoute(
      "/snapshot-123/abstracts/abstract-123/confirm-delete",
    );

    testForbiddenPublicationRoute("/lay_summaries/add");
    testForbiddenPublicationRoute("/lay_summaries", "POST");
    testForbiddenPublicationRoute("/lay_summaries/lay-summary-123/edit");
    testForbiddenPublicationRoute(
      "/lay_summaries/lay-summary-123",
      "PUT",
      "DELETE",
    );
    testForbiddenPublicationRoute(
      "/snapshot-123/lay_summaries/lay-summary-123/confirm-delete",
    );

    testForbiddenPublicationRoute("/datasets/add");
    testForbiddenPublicationRoute("/datasets/suggestions");
    testForbiddenPublicationRoute("/datasets", "POST");
    testForbiddenPublicationRoute(
      "/snapshot-123/datasets/dataset-123/confirm-delete",
    );
    testForbiddenPublicationRoute("/datasets/dataset-123", "DELETE");

    testForbiddenPublicationRoute("/contributors/role-123/order", "POST");
    testForbiddenPublicationRoute("/contributors/role-123/add");
    testForbiddenPublicationRoute("/contributors/role-123/suggestions");
    testForbiddenPublicationRoute("/contributors/role-123/confirm-create");
    testForbiddenPublicationRoute("/contributors/role-123", "POST");
    testForbiddenPublicationRoute("/contributors/role-123/position-123/edit");
    testForbiddenPublicationRoute(
      "/contributors/role-123/position-123/suggestions",
    );
    testForbiddenPublicationRoute(
      "/contributors/role-123/position-123/confirm-update",
    );
    testForbiddenPublicationRoute(
      "/contributors/role-123/position-123/confirm-delete",
    );
    testForbiddenPublicationRoute(
      "/contributors/role-123/position-123",
      "PUT",
      "DELETE",
    );

    testForbiddenPublicationRoute("/files", "POST");
    testForbiddenPublicationRoute("/files/file-123/edit");
    testForbiddenPublicationRoute("/refresh-files");
    testForbiddenPublicationRoute("/files/file-123/refresh-form");
    testForbiddenPublicationRoute("/files/file-123", "PUT", "DELETE");
    testForbiddenPublicationRoute("/snapshot123/files/file-123/confirm-delete");
  });

  it("should not be possible to edit or delete datasets from another user", () => {
    cy.loginAsLibrarian();
    cy.setUpDataset();

    cy.loginAsResearcher();

    testForbiddenDatasetRoute("/add/description");
    testForbiddenDatasetRoute("/add/confirm");
    testForbiddenDatasetRoute("/save", "POST");
    testForbiddenDatasetRoute("/add/publish", "POST");
    testForbiddenDatasetRoute("/add/finish");

    testForbiddenDatasetRoute(""); // The regular view dataset route
    testForbiddenDatasetRoute("/description");
    testForbiddenDatasetRoute("/contributors");
    testForbiddenDatasetRoute("/publications");
    testForbiddenDatasetRoute("/activity");

    testForbiddenDatasetRoute("/confirm-delete");
    testForbiddenDatasetRoute("", "DELETE");
    testForbiddenDatasetRoute("/lock", "POST");
    testForbiddenDatasetRoute("/unlock", "POST");
    testForbiddenDatasetRoute("/publish/confirm");
    testForbiddenDatasetRoute("/publish", "POST");
    testForbiddenDatasetRoute("/withdraw/confirm");
    testForbiddenDatasetRoute("/withdraw", "POST");
    testForbiddenDatasetRoute("/republish/confirm");
    testForbiddenDatasetRoute("/republish", "POST");

    testForbiddenDatasetRoute("/message/edit", "GET", "PUT");
    testForbiddenDatasetRoute("/message", "PUT");

    testForbiddenDatasetRoute("/details/edit", "GET", "PUT");
    testForbiddenDatasetRoute("/details/edit/refresh", "PUT");

    testForbiddenDatasetRoute("/projects/add");
    testForbiddenDatasetRoute("/projects/suggestions");
    testForbiddenDatasetRoute("/projects", "POST");
    testForbiddenDatasetRoute(
      "/snapshot-123/projects/confirm-delete/project-123",
    );
    testForbiddenDatasetRoute("/projects/project-123", "DELETE");

    testForbiddenDatasetRoute("/links/add");
    testForbiddenDatasetRoute("/links", "POST");
    testForbiddenDatasetRoute("/links/link-123/edit");
    testForbiddenDatasetRoute("/links/link-123", "PUT", "DELETE");
    testForbiddenDatasetRoute("/snapshot-123/links/link-123/confirm-delete");

    testForbiddenDatasetRoute("/departments/add");
    testForbiddenDatasetRoute("/departments/suggestions");
    testForbiddenDatasetRoute("/departments", "POST");
    testForbiddenDatasetRoute(
      "/snapshot-123/departments/department-123/confirm-delete",
    );
    testForbiddenDatasetRoute("/departments/department-123", "DELETE");

    testForbiddenDatasetRoute("/abstracts/add");
    testForbiddenDatasetRoute("/abstracts", "POST");
    testForbiddenDatasetRoute("/abstracts/abstract-123/edit");
    testForbiddenDatasetRoute("/abstracts/abstract-123", "PUT", "DELETE");
    testForbiddenDatasetRoute(
      "/snapshot-123/abstracts/abstract-123/confirm-delete",
    );

    testForbiddenDatasetRoute("/publications/add");
    testForbiddenDatasetRoute("/publications/suggestions");
    testForbiddenDatasetRoute("/publications", "POST");
    testForbiddenDatasetRoute(
      "/snapshot-123/publications/publication-123/confirm-delete",
    );
    testForbiddenDatasetRoute("/publications/publication-123", "DELETE");

    testForbiddenDatasetRoute("/contributors/role-123/order", "POST");
    testForbiddenDatasetRoute("/contributors/role-123/add");
    testForbiddenDatasetRoute("/contributors/role-123/suggestions");
    testForbiddenDatasetRoute("/contributors/role-123/confirm-create");
    testForbiddenDatasetRoute("/contributors/role-123", "POST");
    testForbiddenDatasetRoute("/contributors/role-123/position-123/edit");
    testForbiddenDatasetRoute(
      "/contributors/role-123/position-123/suggestions",
    );
    testForbiddenDatasetRoute(
      "/contributors/role-123/position-123/confirm-update",
    );
    testForbiddenDatasetRoute(
      "/contributors/role-123/position-123/confirm-delete",
    );
    testForbiddenDatasetRoute(
      "/contributors/role-123/position-123",
      "PUT",
      "DELETE",
    );
  });

  it("should not be possible to edit publication reviewer tags and notes as a regular user", () => {
    cy.loginAsResearcher();

    cy.setUpPublication();

    testUnauthorizedPublicationRoute("/reviewer-tags/edit");
    testUnauthorizedPublicationRoute("/reviewer-tags", "PUT");
    testUnauthorizedPublicationRoute("/reviewer-note/edit");
    testUnauthorizedPublicationRoute("/reviewer-note", "PUT");
  });

  it("should not be possible to edit dataset reviewer tags and notes as a regular user", () => {
    cy.loginAsResearcher();

    cy.setUpDataset();

    testUnauthorizedDatasetRoute("/reviewer-tags/edit");
    testUnauthorizedDatasetRoute("/reviewer-tags", "PUT");
    testUnauthorizedDatasetRoute("/reviewer-note/edit");
    testUnauthorizedDatasetRoute("/reviewer-note", "PUT");
  });

  type HttpMethods = ("GET" | "PUT" | "POST" | "DELETE")[];

  function testForbiddenDatasetRoute(route: string, ...methods: HttpMethods) {
    cy.then(function () {
      testRouteHttpStatus(403, `/dataset/${this.biblioId}${route}`, ...methods);
    });
  }

  function testUnauthorizedDatasetRoute(
    route: string,
    ...methods: HttpMethods
  ) {
    cy.then(function () {
      testRouteHttpStatus(401, `/dataset/${this.biblioId}${route}`, ...methods);
    });
  }

  function testForbiddenPublicationRoute(
    route: string,
    ...methods: HttpMethods
  ) {
    cy.then(function () {
      testRouteHttpStatus(
        403,
        `/publication/${this.biblioId}${route}`,
        ...methods,
      );
    });
  }

  function testUnauthorizedPublicationRoute(
    route: string,
    ...methods: HttpMethods
  ) {
    cy.then(function () {
      testRouteHttpStatus(
        401,
        `/publication/${this.biblioId}${route}`,
        ...methods,
      );
    });
  }

  function testRouteHttpStatus(
    httpStatus: number,
    url: string,
    ...methods: HttpMethods
  ) {
    if (methods.length === 0) {
      methods.push("GET");
    }

    return cy.then(function () {
      for (const method of methods) {
        return cy
          .request({
            url,
            method,
            headers: {
              "X-CSRF-Token": this.CSRFToken,
            },
            followRedirect: false,
            failOnStatusCode: false,
          })
          .should((response) => {
            expect(response.status).to.equal(httpStatus);
          });
      }
    });
  }
});
