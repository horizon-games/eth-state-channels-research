const path = require('path')
const fs = require('fs')
const webpack = require('webpack')
const ForkTsCheckerWebpackPlugin = require('fork-ts-checker-webpack-plugin')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const shared = require('./shared')

const apps = {
  // 'index': {
  //   src:      './src/index.ts',
  //   template: './src/index.html'
  // },
  'ttt': {
    src:      './examples/ttt/index.ts',
    template: './examples/ttt/index.html'
  }
}

const main = [
  'webpack-dev-server/client?http://0.0.0.0:3000',
  'webpack/hot/only-dev-server',
]
const vendor = shared.vendorEntry({
  mainModules: main,
  modulesToExclude: ['']
})

const buildConfig = (name, app) => {
  const config = {
    context: process.cwd(), // to automatically find tsconfig.json
    entry: {},
    output: {
      path: path.resolve(__dirname, 'dist'),
      filename: '[name].js',
      publicPath: '/'
    },
    plugins: [
      new webpack.optimize.CommonsChunkPlugin({
        name: 'vendor'
      }),
      new webpack.NamedModulesPlugin(),
      new webpack.HotModuleReplacementPlugin(),
      new ForkTsCheckerWebpackPlugin({
        tslint: true,
        checkSyntacticErrors: true,
        watch: ['./client'] // optional but improves performance (fewer stat calls)
      }),
      // new webpack.NoEmitOnErrorsPlugin(),
      // new webpack.DefinePlugin(shared.appEnvVars('config/app.dev.env')),
      new webpack.DefinePlugin({
        'process.env.NODE_ENV': JSON.stringify('development'),
      }),
      new HtmlWebpackPlugin({
        inject: 'head',
        template: app.template,
        filename: path.basename(app.template)
      }),
    ],
    module: {
      rules: [{
        test: /.tsx?$/,
        use: [{
          loader: 'ts-loader', options: { transpileOnly: true }
        }],
        exclude: path.resolve(process.cwd(), 'node_modules'),
        include: [
          path.resolve(process.cwd(), 'client'),
          path.resolve(process.cwd(), 'examples')
        ]
      }]
    },
    resolve: {
      extensions: ['.tsx', '.ts', '.js'],
      alias: {
        dgame: path.join(process.cwd(), 'client', 'dgame'),
        wsrelay: path.join(process.cwd(), 'client', 'wsrelay')
      }
    },
    devtool: 'inline-source-map',
    devServer: {
      host: '0.0.0.0',
      port: 3000,
      open: false,
      hot: true,
      historyApiFallback: true,
      stats: 'errors-only'
    }
  }
  config.entry[name] = [...main, app.src]
  config.entry['vendor'] = vendor
  return config
}

const configs = []
for (const [k, v] of Object.entries(apps)) {
  configs.push(buildConfig(k, v))
}

module.exports = configs
