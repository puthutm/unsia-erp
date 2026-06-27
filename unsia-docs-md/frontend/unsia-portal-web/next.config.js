/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  output: 'standalone',
  transpilePackages: [],
  env: {
    NEXT_PUBLIC_APP_NAME: 'unsia-portal-web',
  },
  experimental: {
    webpackBuildWorker: false
  },
  webpack: (config) => {
    config.resolve.alias = {
      ...config.resolve.alias,
      '@/': require('path').resolve(__dirname, './'),
      '@/contexts': require('path').resolve(__dirname, './contexts'),
      '@/hooks': require('path').resolve(__dirname, './hooks'),
      '@/lib': require('path').resolve(__dirname, './lib'),
      '@/components': require('path').resolve(__dirname, './components'),
      '@/app': require('path').resolve(__dirname, './app'),
    };
    return config;
  },
};

module.exports = nextConfig;
