/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  output: 'standalone',
  transpilePackages: [],
  env: {
    NEXT_PUBLIC_APP_NAME: 'unsia-assessment',
  },
  experimental: {
    webpackBuildWorker: false
  }
};

module.exports = nextConfig;
