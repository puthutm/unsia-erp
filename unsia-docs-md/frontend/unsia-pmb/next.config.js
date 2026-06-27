/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  transpilePackages: [],
  output: 'standalone',
  experimental: {
    webpackBuildWorker: false
  }
};

module.exports = nextConfig;
