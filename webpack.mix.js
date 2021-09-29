const mix = require('laravel-mix')

mix.js('./assets/js/app.js', 'js')
mix.sass('./assets/css/app.scss', 'css')
mix.setPublicPath('./static')

// copy images
mix.copy('assets/ugent/images/**/*', 'static/images');
mix.copy('assets/ugent/favicon.ico', 'static/favicon.ico');

// set the resourceroot for fonts so it points to the static assets path
mix.setResourceRoot('/static/fonts/')
mix.webpackConfig({
  module: {
    rules: [
      // copy font files to the ./static/fonts folder
      {
        test: '/(\\.(woff2?|ttf|eot|otf)$|font.*\\.svg$)/',
        use: [{
          loader: 'file-loader',
          options: {
            name: '[name].[ext]',
            outputPath: './fonts/'
          }
        }]
      }
    ]
  }
})

if (mix.inProduction()) {
  mix.version()
} else {
  mix.sourceMaps()
}
