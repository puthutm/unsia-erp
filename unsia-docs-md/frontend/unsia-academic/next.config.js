/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  transpilePackages: [],
  experimental: {
    webpackBuildWorker: false
  },
  output: 'standalone',
}

module.exports = nextConfig
