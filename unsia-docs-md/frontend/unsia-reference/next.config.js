/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  transpilePackages: ['@phosphor-icons/react'],
  experimental: {
    webpackBuildWorker: false
  }
};

module.exports = nextConfig;
