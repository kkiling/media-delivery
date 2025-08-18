const path = require('path');

module.exports = {
  webpack: {
    alias: {
      '@': path.resolve(__dirname, 'src'),
    },
  },
  devServer: {
    proxy: {
      '/v1': {
        target: 'http://localhost:8083',
        changeOrigin: true,
      }
    }
  }
};