const _transform = require('lodash/transform')
const webpack = require('webpack')
const path = require('path')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const ESLintPlugin = require('eslint-webpack-plugin')
const NodePolyfillPlugin = require('node-polyfill-webpack-plugin')
const ForkTsCheckerWebpackPlugin = require('fork-ts-checker-webpack-plugin')
const FileManagerPlugin = require('filemanager-webpack-plugin')

const host = process.env.HOST
const protocol = process.env.PROTOCOL
const buildPath = process.env.SOURCE_PATH
const port = parseInt(process.env.PORT, 10)
const projectRoot = path.resolve(__dirname, '../')

const publicPath = `${protocol}://${host}:${port}/${buildPath}/`

/** переопределяем таргет сборки для автопрефиксера */
const autoprefixerBrowserslist = [
  'last 1 chrome version',
  'last 1 firefox version',
  'last 1 safari version',
]

const plugins = [
  new webpack.ProvidePlugin({
    process: 'process/browser',
  }),
  new webpack.DefinePlugin({
    ..._transform(
      process.env,
      (result, value, key) => {
        const parsedValue = parseFloat(value)
        result[`process.env.${key}`] = isNaN(parsedValue)
          ? `'${value}'`
          : parsedValue
      },
      {}
    ),
  }),
  new MiniCssExtractPlugin({
    filename: '[name]-[chunkhash].css',
  }),
  new HtmlWebpackPlugin({
    filename: 'index.html',
    template: path.resolve(projectRoot, 'public', 'index.html'),
    inject: 'body',
    publicPath,
  }),
  new NodePolyfillPlugin(),
  new ESLintPlugin({
    extensions: ['js', 'jsx', 'ts', 'tsx'],
    overrideConfigFile: path.join(__dirname, '..', '.eslintrc.js'),
  }),
  new FileManagerPlugin({
    events: {
      onStart: {
        delete: [buildPath],
      },
    },
    runOnceInWatchMode: true,
  }),
  new webpack.ProgressPlugin(),
  new ForkTsCheckerWebpackPlugin({
    async: true,
    typescript: {
      configOverwrite: {
        compilerOptions: {
          noEmit: true,
          skipLibCheck: true,
        },
      },
    },
    issue: {
      exclude: [
        { file: 'node_modules/**/*.tsx' },
        { file: 'node_modules/**/*.ts' },
      ],
    },
  }),
]

const devServer = {
  server: {
    type: 'http',
  },
  host,
  static: {
    directory: path.join(__dirname, '..', buildPath),
  },
  headers: { 'Access-Control-Allow-Origin': '*' },
  historyApiFallback: true, // Apply HTML5 History API if routes are used
  open: true,
  compress: true,
  allowedHosts: 'all',
  hot: true, // Reload the page after changes saved (HotModuleReplacementPlugin)
  client: {
    // Shows a full-screen overlay in the browser when there are compiler errors or warnings
    overlay: {
      errors: false,
      warnings: false,
      runtimeErrors: (error) => {
        if (
          error.message ===
          'ResizeObserver loop completed with undelivered notifications.'
        ) {
          return false
        }
        return true
      },
    },
  },
  port: 3000,
  /**
   * Writes files to output path (default: false)
   * Build dir is not cleared using <output: {clean:true}>
   * To resolve should use FileManager
   */
  devMiddleware: {
    writeToDisk: true,
  },
}

module.exports = {
  mode: 'development',
  target: 'web',
  devtool: 'inline-source-map',
  devServer,
  plugins,
  context: projectRoot,
  ignoreWarnings: [{ module: /node_modules/ }],
  watchOptions: {
    ignored: /node_modules/,
  },
  entry: {
    bundle: path.join(__dirname, '..', 'src', 'index.ts'),
  },
  output: {
    path: path.join(__dirname, '..', buildPath),
    filename: '[name].js',
    chunkFilename: '[name].js',
    publicPath,
  },
  // Checking the maximum weight of the bundle is disabled
  performance: {
    hints: false,
  },
  // Modules resolved
  resolve: {
    modules: ['src', 'node_modules'],
    extensions: ['.tsx', '.ts', '.js', '.scss'],
    fallback: {
      module: false,
      fs: false,
    },
  },
  module: {
    exprContextCritical: false,
    strictExportPresence: true, // Strict mod to avoid of importing non-existent objects
    rules: [
      {
        test: /\.[jt]sx?$/,
        exclude: /node_modules/,
        use: {
          loader: 'babel-loader',
          options: {
            cacheDirectory: true, // Using a cache to avoid of recompilation
          },
        },
      },
      {
        test: /\.module\.(scss|sass)$/,
        use: [
          'style-loader',
          {
            loader: 'css-loader',
            options: {
              modules: {
                mode: 'local',
                namedExport: false,
                exportLocalsConvention: 'camel-case',
                localIdentName: '[name]_[local]__[hash]',
              },
              importLoaders: 2,
              sourceMap: true,
            },
          },
          {
            loader: 'postcss-loader',
            options: {
              postcssOptions: {
                plugins: [
                  {
                    autoprefixer: {
                      overrideBrowserslist: autoprefixerBrowserslist,
                    },
                  },
                ],
              },
            },
          },
          {
            loader: 'sass-loader',
            options: {
              sassOptions: {
                outputStyle: 'expanded',
              },
              sourceMap: true,
            },
          },
        ],
      },
      {
        test: /\.(scss|sass)$/,
        exclude: /\.module\.(scss|sass)$/,
        use: [
          'style-loader',
          {
            loader: 'css-loader',
            options: {
              modules: false,
              importLoaders: 2,
              sourceMap: true,
            },
          },
          {
            loader: 'postcss-loader',
            options: {
              postcssOptions: {
                plugins: [
                  {
                    autoprefixer: {
                      overrideBrowserslist: autoprefixerBrowserslist,
                    },
                  },
                ],
              },
            },
          },
          {
            loader: 'sass-loader',
            options: {
              sassOptions: {
                outputStyle: 'expanded',
              },
              sourceMap: true,
            },
          },
        ],
      },
    ],
  },
}
