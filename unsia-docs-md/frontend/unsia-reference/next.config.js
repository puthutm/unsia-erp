/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  transpilePackages: ['@phosphor-icons/react'],
  experimental: {
    webpackBuildWorker: false
  },
  output: 'standalone',
};

module.exports = nextConfig;
