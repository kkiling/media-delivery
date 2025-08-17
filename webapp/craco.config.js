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
        target: 'http://10.10.10.202:8083',
        changeOrigin: true,
      }
    }
  }
};