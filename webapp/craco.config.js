const path = require('path');
const TerserPlugin = require('terser-webpack-plugin');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const CssMinimizerPlugin = require('css-minimizer-webpack-plugin');

module.exports = {
  webpack: {
    alias: {
      '@': path.resolve(__dirname, 'src'),
    },
    configure: (webpackConfig) => {
      // Оптимизация для продакшена
      if (process.env.NODE_ENV === 'production') {
        webpackConfig.optimization = {
          ...webpackConfig.optimization,
          minimizer: [
            new TerserPlugin({
              terserOptions: {
                compress: {
                  drop_console: true, // удаляет console.log
                },
              },
            }),
            new CssMinimizerPlugin(),
          ],
          splitChunks: {
            chunks: 'all',
            minSize: 20000,
            minChunks: 1,
            maxAsyncRequests: 30,
            maxInitialRequests: 30,
            cacheGroups: {
              defaultVendors: {
                test: /[\\/]node_modules[\\/]/,
                priority: -10,
                reuseExistingChunk: true,
              },
              default: {
                minChunks: 2,
                priority: -20,
                reuseExistingChunk: true,
              },
            },
          },
        };

        // Отключаем source maps для продакшена
        webpackConfig.devtool = false;
      }

      return webpackConfig;
    },
  },
  devServer: {
    proxy: {
      '/v1': {
        target: 'http://localhost:8083',
        changeOrigin: true,
      },
    },
  },
};
