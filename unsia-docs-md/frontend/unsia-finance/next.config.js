/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  transpilePackages: [],
  experimental: {
    webpackBuildWorker: false
  }
}

module.exports = nextConfig
