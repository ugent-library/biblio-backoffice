import { defineConfig } from 'cypress'

export default defineConfig({
  e2e: {
    baseUrl: 'https://backoffice.bibliotest.ugent.be',
    experimentalStudio: true,

    setupNodeEvents(on, config) {
      // implement node event listeners here
    },
  },
})
