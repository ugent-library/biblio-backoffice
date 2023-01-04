import { defineConfig } from 'cypress'
import dotenvPlugin from 'cypress-dotenv'

export default defineConfig({
  projectId: 'mjg74d',

  e2e: {
    baseUrl: 'https://backoffice.bibliotest.ugent.be',
    experimentalStudio: true,

    setupNodeEvents(_on, config) {
      config = dotenvPlugin(config)

      return config
    },
  },
})
