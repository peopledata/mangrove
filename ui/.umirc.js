// https://umijs.org/config/
import { resolve } from 'path'
const fs = require('fs')
const path = require('path')
const lessToJs = require('less-vars-to-js')
const isDevelopment = process.env.NODE_ENV === 'development'

const { theme } = require('antd/lib')
const { convertLegacyToken } = require('@ant-design/compatible/lib')

const mapToken = theme.defaultAlgorithm(theme.defaultSeed)
const v4Token = convertLegacyToken(mapToken)

// how to speed compile: https://umijs.org/guide/boost-compile-speed
export default {
  // IMPORTANT! change next line to yours or delete. And hide in dev
  // publicPath: isDevelopment ? '/' : 'https://cdn.antd-admin.zuiidea.com/',
  publicPath: '/',
  alias: {
    api: resolve(__dirname, './src/services/'),
    components: resolve(__dirname, './src/components'),
    config: resolve(__dirname, './src/utils/config'),
    themes: resolve(__dirname, './src/themes'),
    utils: resolve(__dirname, './src/utils'),
  },
  antd: false,
  // a lower cost way to genereate sourcemap, default is cheap-module-source-map, could save 60% time in dev hotload
  devtool: 'eval',
  dva: { immer: true },
  dynamicImport: {
    loading: 'components/Loader/Loader',
  },
  extraBabelPlugins: [
    [
      'import',
      {
        libraryName: 'lodash',
        libraryDirectory: '',
        camel2DashComponentName: false,
      },
      'lodash',
    ],
    [
      'import',
      {
        libraryName: '@ant-design/icons',
        libraryDirectory: 'es/icons',
        camel2DashComponentName: false,
      },
      'ant-design-icons',
    ],
    ['macros'],
  ],
  hash: true,
  ignoreMomentLocale: true,
  // umi3 comple node_modules by default, could be disable
  nodeModulesTransform: {
    type: 'none',
    exclude: [],
  },
  mock: false,
  // Webpack Configuration
  proxy: {
    '/admin': {
      target: 'http://localhost:8081/',
      changeOrigin: true,
      // pathRewrite: { '^/api': '/api' },
    },
  },
  // Theme for antd
  // https://ant.design/docs/react/customize-theme
  theme: {
    ...v4Token,
    ...lessToJs(
      fs.readFileSync(path.join(__dirname, './src/themes/default.less'), 'utf8')
    ),
  },
  webpack5: {},
  mfsu: {},
  chainWebpack: function (config, { webpack }) {
    !isDevelopment &&
      config.merge({
        optimization: {
          minimize: false,
          splitChunks: {
            chunks: 'all',
            minSize: 30000,
            minChunks: 3,
            automaticNameDelimiter: '.',
            cacheGroups: {
              react: {
                name: 'react',
                priority: 20,
                test: /[\\/]node_modules[\\/](react|react-dom|react-dom-router)[\\/]/,
              },
              antd: {
                name: 'antd',
                priority: 20,
                test: /[\\/]node_modules[\\/](antd|@ant-design\/icons)[\\/]/,
              },
              'echarts-gl': {
                name: 'echarts-gl',
                priority: 30,
                test: /[\\/]node_modules[\\/]echarts-gl[\\/]/,
              },
              echarts: {
                name: 'echarts',
                priority: 20,
                test: /[\\/]node_modules[\\/](echarts|echarts-for-react|echarts-liquidfill)[\\/]/,
              },
              highcharts: {
                name: 'highcharts',
                priority: 20,
                test: /[\\/]node_modules[\\/]highcharts[\\/]/,
              },
              recharts: {
                name: 'recharts',
                priority: 20,
                test: /[\\/]node_modules[\\/]recharts[\\/]/,
              },
              draftjs: {
                name: 'draftjs',
                priority: 30,
                test: /[\\/]node_modules[\\/](draft-js|react-draft-wysiwyg|draftjs-to-html|draftjs-to-markdown)[\\/]/,
              },
              async: {
                chunks: 'async',
                minChunks: 2,
                name: 'async',
                maxInitialRequests: 1,
                minSize: 0,
                priority: 5,
                reuseExistingChunk: true,
              },
            },
          },
        },
      })
  },
}
