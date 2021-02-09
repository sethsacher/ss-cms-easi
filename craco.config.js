const sassResourcesLoader = require('craco-sass-resources-loader');
const TerserPlugin = require('terser-webpack-plugin');

module.exports = {
  webpack: {
    configure: (webpackConfig) => ({
      ...webpackConfig,
      optimization: {
        ...webpackConfig.optimization,
        // Workaround for CircleCI bug caused by the number of CPUs shown
        // https://github.com/facebook/create-react-app/issues/8320
        minimizer: webpackConfig.optimization.minimizer.map(item => {
          if (item instanceof TerserPlugin) {
            item.options.parallel = 0;
          }

          return item;
        })
      },
    }),
  },
  plugins: [
    {
      plugin: sassResourcesLoader,
      options: {
        resources: [
          './src/stylesheets/_uswdsUtilities.scss',
          './src/stylesheets/_colors.scss',
          './src/stylesheets/_variables.scss'
        ]
      }
    }
  ]
};