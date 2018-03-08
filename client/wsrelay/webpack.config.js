//
// This is to generate the umd bundle only
//
const webpack = require('webpack')
const path = require('path')
const ForkTsCheckerWebpackPlugin = require('fork-ts-checker-webpack-plugin')
const UglifyJSPlugin = require('uglifyjs-webpack-plugin')
const production = process.env.NODE_ENV === 'production'

let entry = {
  'wsrelay': './index.ts',
};
if (production) {
  entry = Object.assign({}, entry, {'wsrelay.min': './index.ts'});
}

module.exports = {
  entry,
  output: {
    path: path.resolve(__dirname, 'dist/umd'),
    filename: '[name].js',
    libraryTarget: 'umd',
    library: 'dgame',
    umdNamedDefine: true,
  },
  module: {
    rules: [{
      test: /.tsx?$/,
      use: [{
        loader: 'ts-loader', options: { transpileOnly: true }
      }],
      exclude: /node_modules/
    }]
  },
  resolve: {
    extensions: ['.tsx', '.ts', '.js'],
    alias: {
      src: path.join(process.cwd(), 'src')
    }
  },
  plugins: [
    new webpack.NoEmitOnErrorsPlugin(),
    new UglifyJSPlugin({
      // sourceMap: true,
      include: /\.min\.js$/,
      uglifyOptions: {
        ie8: false,
        compress: {
          dead_code: true,
          unused: true
        },
        output: {
          comments: false,
          beautify: false
        }
      }
    }),
    new ForkTsCheckerWebpackPlugin({
      async: false,
      memoryLimit: 4096,
      checkSyntacticErrors: true
    })
  ]
}
